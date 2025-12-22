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

// ProfileRepository реализация получения профиля
type ProfileRepository struct {
	db *sqlx.DB
}

// NewProfileRepository создает новый репозиторий
func NewProfileRepository(db *sqlx.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetByID возвращает профиль игрока по ID
// Объединяет данные из всех источников (Junior/FHSPB) по name+birth_date
func (r *ProfileRepository) GetByID(ctx context.Context, playerID string) (*services.PlayerProfile, error) {
	// Сначала получаем базовую информацию об игроке и находим все связанные ID
	query := `
		WITH target_player AS (
			SELECT name, birth_date FROM players WHERE id = $1
		),
		all_player_ids AS (
			SELECT p.id FROM players p, target_player tp
			WHERE p.name = tp.name AND p.birth_date = tp.birth_date
		),
		latest_data AS (
			SELECT DISTINCT ON (1)
				p.id, p.name, 
				COALESCE(EXTRACT(YEAR FROM p.birth_date)::int, 0) as birth_year,
				p.position, p.height, p.weight,
				t.name as team,
				COALESCE(p.region, '') as region,
				tr.start_date
			FROM players p
			JOIN all_player_ids api ON p.id = api.id
			LEFT JOIN player_teams pt ON pt.player_id = p.id
			LEFT JOIN teams t ON pt.team_id = t.id
			LEFT JOIN tournaments tr ON pt.tournament_id = tr.id
			ORDER BY 1, tr.start_date DESC NULLS LAST
		)
		SELECT id, name, birth_year, 
		       COALESCE(position, '') as position, 
		       height, weight, 
		       COALESCE(team, '') as team, 
		       region
		FROM latest_data
		ORDER BY start_date DESC NULLS LAST
		LIMIT 1
	`

	var basicInfo services.PlayerBasicInfo
	err := r.db.QueryRowxContext(ctx, query, playerID).Scan(
		&basicInfo.ID, &basicInfo.Name, &basicInfo.BirthYear,
		&basicInfo.Position, &basicInfo.Height, &basicInfo.Weight,
		&basicInfo.Team, &basicInfo.Region,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("player not found")
	}
	if err != nil {
		return nil, err
	}

	return &services.PlayerProfile{BasicInfo: basicInfo}, nil
}

// GetStats возвращает статистику игрока за всё время
// Объединяет статистику из всех источников (Junior/FHSPB)
func (r *ProfileRepository) GetStats(ctx context.Context, playerID string) (*services.PlayerStats, error) {
	query := `
		WITH target_player AS (
			SELECT name, birth_date FROM players WHERE id = $1
		),
		all_player_ids AS (
			SELECT p.id FROM players p, target_player tp
			WHERE p.name = tp.name AND p.birth_date = tp.birth_date
		)
		SELECT 
			COUNT(DISTINCT tournament_id) as tournaments,
			COALESCE(SUM(games), 0) as games,
			COALESCE(SUM(goals), 0) as goals,
			COALESCE(SUM(assists), 0) as assists,
			COALESCE(SUM(points), 0) as points,
			COALESCE(SUM(plus_minus), 0) as plus_minus,
			COALESCE(SUM(penalty_minutes), 0) as penalties,
			COALESCE(SUM(hat_tricks), 0) as hat_tricks,
			COALESCE(SUM(game_winning_goals), 0) as gwg
		FROM player_statistics ps
		JOIN all_player_ids api ON ps.player_id = api.id
	`

	var stats services.PlayerStats
	err := r.db.QueryRowxContext(ctx, query, playerID).Scan(
		&stats.Tournaments, &stats.Games, &stats.Goals, &stats.Assists,
		&stats.Points, &stats.PlusMinus, &stats.Penalties,
		&stats.HatTricks, &stats.GameWinningGoals,
	)
	if err != nil {
		return nil, err
	}

	if stats.Games == 0 {
		return nil, nil
	}

	stats.GoalsPerGame = calcAvg(stats.Goals, stats.Games)
	stats.AssistsPerGame = calcAvg(stats.Assists, stats.Games)
	stats.PointsPerGame = calcAvg(stats.Points, stats.Games)
	stats.PenaltiesPerGame = calcAvg(stats.Penalties, stats.Games)

	return &stats, nil
}

// GetSeasonStats возвращает статистику за сезон
// Объединяет статистику из всех источников (Junior/FHSPB)
func (r *ProfileRepository) GetSeasonStats(ctx context.Context, playerID, season string) (*services.SeasonStats, error) {
	query := `
		WITH target_player AS (
			SELECT name, birth_date FROM players WHERE id = $1
		),
		all_player_ids AS (
			SELECT p.id FROM players p, target_player tp
			WHERE p.name = tp.name AND p.birth_date = tp.birth_date
		)
		SELECT 
			COUNT(DISTINCT ps.tournament_id) as tournaments,
			COALESCE(SUM(ps.games), 0) as games,
			COALESCE(SUM(ps.goals), 0) as goals,
			COALESCE(SUM(ps.assists), 0) as assists,
			COALESCE(SUM(ps.points), 0) as points,
			COALESCE(SUM(ps.plus_minus), 0) as plus_minus,
			COALESCE(SUM(ps.penalty_minutes), 0) as penalties,
			COALESCE(SUM(ps.hat_tricks), 0) as hat_tricks,
			COALESCE(SUM(ps.game_winning_goals), 0) as gwg
		FROM player_statistics ps
		JOIN all_player_ids api ON ps.player_id = api.id
		JOIN tournaments t ON ps.tournament_id = t.id
		WHERE t.season = $2
	`

	var stats services.SeasonStats
	stats.Season = season
	err := r.db.QueryRowxContext(ctx, query, playerID, season).Scan(
		&stats.Tournaments, &stats.Games, &stats.Goals, &stats.Assists,
		&stats.Points, &stats.PlusMinus, &stats.Penalties,
		&stats.HatTricks, &stats.GameWinningGoals,
	)
	if err != nil {
		return nil, err
	}

	if stats.Games == 0 {
		return nil, nil
	}

	stats.GoalsPerGame = calcAvg(stats.Goals, stats.Games)
	stats.AssistsPerGame = calcAvg(stats.Assists, stats.Games)
	stats.PointsPerGame = calcAvg(stats.Points, stats.Games)
	stats.PenaltiesPerGame = calcAvg(stats.Penalties, stats.Games)

	return &stats, nil
}

// GetCurrentSeason возвращает текущий сезон
func (r *ProfileRepository) GetCurrentSeason() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	if month >= time.September {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}

// GetRecentTournaments возвращает последние турниры игрока
// Объединяет турниры из всех источников (Junior/FHSPB)
// Логика: если у игрока для турнира есть группы кроме "Общая статистика",
// то "Общую статистику" не показываем. Если только "Общая статистика" - показываем.
// Сортировка: сезон (новые→старые) → дата турнира (новые→старые)
func (r *ProfileRepository) GetRecentTournaments(ctx context.Context, playerID string, limit int) ([]*services.ProfileTournamentStats, error) {
	query := `
		WITH target_player AS (
			SELECT name, birth_date FROM players WHERE id = $1
		),
		all_player_ids AS (
			SELECT p.id FROM players p, target_player tp
			WHERE p.name = tp.name AND p.birth_date = tp.birth_date
		),
		tournament_groups AS (
			SELECT 
				ps.tournament_id,
				array_agg(DISTINCT ps.group_name) as groups
			FROM player_statistics ps
			JOIN all_player_ids api ON ps.player_id = api.id
			GROUP BY ps.tournament_id
		)
		SELECT 
			t.name as tournament_name,
			COALESCE(ps.group_name, '') as group_name,
			t.season,
			ps.games, ps.goals, ps.assists, ps.points,
			ps.plus_minus, ps.penalty_minutes as penalties,
			COALESCE(ps.hat_tricks, 0) as hat_tricks,
			COALESCE(ps.game_winning_goals, 0) as winning_goals,
			CASE WHEN t.name ILIKE '%чемпионат%' THEN true ELSE false END as is_championship
		FROM player_statistics ps
		JOIN all_player_ids api ON ps.player_id = api.id
		JOIN tournaments t ON ps.tournament_id = t.id
		JOIN tournament_groups tg ON tg.tournament_id = ps.tournament_id
		WHERE (
		    -- NULL group_name - показываем всегда (FHSPB stats)
		    ps.group_name IS NULL
		    OR
		    -- Показываем "Общая статистика" только если это единственная группа
		    (ps.group_name = 'Общая статистика' AND array_length(tg.groups, 1) = 1)
		    OR
		    -- Показываем все группы кроме "Общая статистика"
		    ps.group_name != 'Общая статистика'
		  )
		ORDER BY t.season DESC, t.start_date DESC NULLS LAST, t.name, ps.group_name
		LIMIT $2
	`

	rows, err := r.db.QueryxContext(ctx, query, playerID, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var tournaments []*services.ProfileTournamentStats
	for rows.Next() {
		var t services.ProfileTournamentStats
		if err := rows.Scan(
			&t.TournamentName, &t.GroupName, &t.Season,
			&t.Games, &t.Goals, &t.Assists, &t.Points,
			&t.PlusMinus, &t.Penalties, &t.HatTricks, &t.WinningGoals, &t.IsChampionship,
		); err != nil {
			continue
		}
		tournaments = append(tournaments, &t)
	}
	return tournaments, nil
}

// GetTournamentsBySeason возвращает турниры сгруппированные по сезонам
func (r *ProfileRepository) GetTournamentsBySeason(ctx context.Context, playerID string) ([]*services.SeasonTournaments, error) {
	tournaments, err := r.GetRecentTournaments(ctx, playerID, 20)
	if err != nil {
		return nil, err
	}

	seasonMap := make(map[string][]*services.ProfileTournamentStats)
	var seasons []string

	for _, t := range tournaments {
		if _, exists := seasonMap[t.Season]; !exists {
			seasons = append(seasons, t.Season)
		}
		seasonMap[t.Season] = append(seasonMap[t.Season], t)
	}

	var result []*services.SeasonTournaments
	for _, season := range seasons {
		result = append(result, &services.SeasonTournaments{
			Season:      season,
			Tournaments: seasonMap[season],
		})
	}
	return result, nil
}

func calcAvg(val, games int) float64 {
	if games == 0 {
		return 0
	}
	return float64(val) / float64(games)
}
