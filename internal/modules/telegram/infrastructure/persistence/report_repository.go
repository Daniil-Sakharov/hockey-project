package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/jmoiron/sqlx"
)

// ReportRepository реализация получения данных для отчёта
type ReportRepository struct {
	db *sqlx.DB
}

// NewReportRepository создает новый репозиторий
func NewReportRepository(db *sqlx.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetFullReport возвращает полные данные для отчёта
func (r *ReportRepository) GetFullReport(ctx context.Context, playerID string) (*services.FullPlayerReport, error) {
	report := &services.FullPlayerReport{}

	// Получаем базовую информацию
	player, err := r.getPlayerInfo(ctx, playerID)
	if err != nil {
		return nil, err
	}
	report.Player = *player

	// Получаем статистику
	stats, err := r.getTotalStats(ctx, playerID)
	if err == nil && stats != nil {
		report.TotalStats = *stats
		report.HasStats = stats.TotalGames > 0

		report.GoalsByPeriod = services.PeriodGoals{
			Period1:  stats.GoalsPeriod1,
			Period2:  stats.GoalsPeriod2,
			Period3:  stats.GoalsPeriod3,
			Overtime: stats.GoalsOvertime,
		}

		report.GoalsByType = services.GoalsBreakdown{
			EvenStrength: stats.GoalsEvenStrength,
			PowerPlay:    stats.GoalsPowerPlay,
			ShortHanded:  stats.GoalsShortHanded,
		}
		report.HasDetailedStats = stats.GoalsEvenStrength > 0 || stats.GoalsPowerPlay > 0
	}

	// Получаем статистику по сезонам
	seasons, _ := r.getSeasonStats(ctx, playerID)
	report.SeasonStats = seasons
	report.HasMultipleSeasons = len(seasons) > 1

	// Получаем турниры
	tournaments, _ := r.getTournaments(ctx, playerID)
	report.Tournaments = tournaments

	return report, nil
}

func (r *ReportRepository) getPlayerInfo(ctx context.Context, playerID string) (*services.ReportPlayerInfo, error) {
	query := `
		SELECT p.id, p.name, 
		       COALESCE(EXTRACT(YEAR FROM p.birth_date)::int, 0) as birth_year,
		       COALESCE(p.position, '') as position, 
		       p.height, p.weight
		FROM players p WHERE p.id = $1
	`

	var info services.ReportPlayerInfo
	var height, weight sql.NullInt64
	err := r.db.QueryRowxContext(ctx, query, playerID).Scan(
		&info.ID, &info.Name, &info.BirthYear, &info.Position, &height, &weight,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("player not found")
	}
	if err != nil {
		return nil, err
	}

	if height.Valid {
		h := int(height.Int64)
		info.Height = &h
	}
	if weight.Valid {
		w := int(weight.Int64)
		info.Weight = &w
	}

	// Получаем команду и регион
	teamQuery := `
		SELECT COALESCE(t.name, ''), COALESCE(t.region, '')
		FROM player_teams pt
		JOIN teams t ON pt.team_id = t.id
		WHERE pt.player_id = $1
		ORDER BY pt.created_at DESC
		LIMIT 1
	`
	_ = r.db.QueryRowxContext(ctx, teamQuery, playerID).Scan(&info.Team, &info.Region)

	return &info, nil
}

func (r *ReportRepository) getTotalStats(ctx context.Context, playerID string) (*services.ReportTotalStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT tournament_id) as tournaments,
			COALESCE(SUM(games), 0) as games,
			COALESCE(SUM(goals), 0) as goals,
			COALESCE(SUM(assists), 0) as assists,
			COALESCE(SUM(points), 0) as points,
			COALESCE(SUM(plus_minus), 0) as plus_minus,
			COALESCE(SUM(penalty_minutes), 0) as penalties,
			COALESCE(SUM(hat_tricks), 0) as hat_tricks,
			COALESCE(SUM(game_winning_goals), 0) as gwg,
			COALESCE(SUM(goals_even_strength), 0) as esg,
			COALESCE(SUM(goals_power_play), 0) as ppg,
			COALESCE(SUM(goals_short_handed), 0) as shg,
			COALESCE(SUM(goals_period_1), 0) as g1,
			COALESCE(SUM(goals_period_2), 0) as g2,
			COALESCE(SUM(goals_period_3), 0) as g3,
			COALESCE(SUM(goals_overtime), 0) as got
		FROM player_statistics
		WHERE player_id = $1
	`

	var stats services.ReportTotalStats
	err := r.db.QueryRowxContext(ctx, query, playerID).Scan(
		&stats.TotalTournaments, &stats.TotalGames, &stats.TotalGoals, &stats.TotalAssists,
		&stats.TotalPoints, &stats.TotalPlusMinus, &stats.TotalPenalties,
		&stats.TotalHatTricks, &stats.TotalWinningGoals,
		&stats.GoalsEvenStrength, &stats.GoalsPowerPlay, &stats.GoalsShortHanded,
		&stats.GoalsPeriod1, &stats.GoalsPeriod2, &stats.GoalsPeriod3, &stats.GoalsOvertime,
	)
	if err != nil {
		return nil, err
	}

	if stats.TotalGames > 0 {
		stats.GoalsPerGame = float64(stats.TotalGoals) / float64(stats.TotalGames)
		stats.AssistsPerGame = float64(stats.TotalAssists) / float64(stats.TotalGames)
		stats.PointsPerGame = float64(stats.TotalPoints) / float64(stats.TotalGames)
		stats.PenaltiesPerGame = float64(stats.TotalPenalties) / float64(stats.TotalGames)
	}

	return &stats, nil
}

func (r *ReportRepository) getSeasonStats(ctx context.Context, playerID string) ([]services.SeasonSummary, error) {
	query := `
		SELECT t.season, 
		       COALESCE(SUM(ps.games), 0) as games,
		       COALESCE(SUM(ps.goals), 0) as goals,
		       COALESCE(SUM(ps.assists), 0) as assists,
		       COALESCE(SUM(ps.points), 0) as points
		FROM player_statistics ps
		JOIN tournaments t ON ps.tournament_id = t.id
		WHERE ps.player_id = $1
		GROUP BY t.season
		ORDER BY t.season
	`

	rows, err := r.db.QueryxContext(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var seasons []services.SeasonSummary
	for rows.Next() {
		var s services.SeasonSummary
		if err := rows.Scan(&s.Season, &s.Games, &s.Goals, &s.Assists, &s.Points); err != nil {
			continue
		}
		seasons = append(seasons, s)
	}
	return seasons, nil
}

func (r *ReportRepository) getTournaments(ctx context.Context, playerID string) ([]services.TournamentStats, error) {
	query := `
		WITH target_player AS (
			SELECT name, birth_date FROM players WHERE id = $1
		),
		all_player_ids AS (
			SELECT p.id FROM players p, target_player tp
			WHERE p.name = tp.name AND p.birth_date = tp.birth_date
		),
		tournament_groups AS (
			SELECT ps.tournament_id, array_agg(DISTINCT ps.group_name) as groups
			FROM player_statistics ps
			JOIN all_player_ids api ON ps.player_id = api.id
			GROUP BY ps.tournament_id
		)
		SELECT t.season, t.name, t.id, COALESCE(ps.group_name, ''), COALESCE(tm.name, ''),
		       ps.games, ps.goals, ps.assists, ps.points, ps.plus_minus, ps.penalty_minutes,
		       ps.hat_tricks, ps.game_winning_goals
		FROM player_statistics ps
		JOIN all_player_ids api ON ps.player_id = api.id
		JOIN tournaments t ON ps.tournament_id = t.id
		JOIN tournament_groups tg ON tg.tournament_id = ps.tournament_id
		LEFT JOIN teams tm ON ps.team_id = tm.id
		WHERE (
			ps.group_name IS NULL
			OR (ps.group_name = 'Общая статистика' AND array_length(tg.groups, 1) = 1)
			OR ps.group_name != 'Общая статистика'
		)
		ORDER BY t.season DESC, t.start_date DESC NULLS LAST, t.name, ps.group_name
	`

	rows, err := r.db.QueryxContext(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var tournaments []services.TournamentStats
	for rows.Next() {
		var ts services.TournamentStats
		if err := rows.Scan(
			&ts.Season, &ts.TournamentName, &ts.TournamentID, &ts.GroupName, &ts.TeamName,
			&ts.Games, &ts.Goals, &ts.Assists, &ts.Points, &ts.PlusMinus, &ts.PenaltyMinutes,
			&ts.HatTricks, &ts.GameWinningGoals,
		); err != nil {
			continue
		}
		tournaments = append(tournaments, ts)
	}
	return tournaments, nil
}

// GetCurrentSeason возвращает текущий сезон
func (r *ReportRepository) GetCurrentSeason() string {
	now := time.Now()
	year := now.Year()
	if now.Month() >= time.September {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}
