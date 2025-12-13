package player_statistics

import (
	"context"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// GetRecentTournaments возвращает последние N турниров игрока
// Сортирует по сезону (новые первыми), затем по дате начала турнира
func (r *repository) GetRecentTournaments(ctx context.Context, playerID string, limit int) ([]*player_statistics.TournamentStat, error) {
	query := `
		SELECT 
			ps.tournament_id,
			t.name as tournament_name,
			ps.group_name,
			t.season,
			CASE 
				WHEN t.start_date IS NOT NULL 
				THEN to_char(t.start_date, 'YYYY-MM-DD')
				ELSE NULL 
			END as start_date,
			
			SUM(ps.games) as games,
			SUM(ps.goals) as goals,
			SUM(ps.assists) as assists,
			SUM(ps.points) as points,
			SUM(ps.plus_minus) as plus_minus,
			SUM(ps.penalty_minutes) as penalty_minutes,
			SUM(ps.hat_tricks) as hat_tricks,
			SUM(ps.game_winning_goals) as game_winning_goals
		FROM player_statistics ps
		INNER JOIN tournaments t ON ps.tournament_id = t.id
		WHERE ps.player_id = $1
		GROUP BY ps.tournament_id, t.name, ps.group_name, t.season, t.start_date
		ORDER BY t.season DESC, t.start_date DESC NULLS LAST, t.name
		LIMIT $2
	`

	type queryResult struct {
		TournamentID     string  `db:"tournament_id"`
		TournamentName   string  `db:"tournament_name"`
		GroupName        string  `db:"group_name"`
		Season           string  `db:"season"`
		StartDate        *string `db:"start_date"`
		Games            int     `db:"games"`
		Goals            int     `db:"goals"`
		Assists          int     `db:"assists"`
		Points           int     `db:"points"`
		PlusMinus        int     `db:"plus_minus"`
		PenaltyMinutes   int     `db:"penalty_minutes"`
		HatTricks        int     `db:"hat_tricks"`
		GameWinningGoals int     `db:"game_winning_goals"`
	}

	var results []queryResult
	err := r.db.SelectContext(ctx, &results, query, playerID, limit)
	if err != nil {
		return nil, err
	}

	// Конвертируем в доменную модель
	tournaments := make([]*player_statistics.TournamentStat, 0, len(results))
	for _, res := range results {
		tournaments = append(tournaments, &player_statistics.TournamentStat{
			TournamentID:     res.TournamentID,
			TournamentName:   res.TournamentName,
			GroupName:        res.GroupName,
			Season:           res.Season,
			IsChampionship:   isChampionship(res.TournamentName),
			StartDate:        res.StartDate,
			Games:            res.Games,
			Goals:            res.Goals,
			Assists:          res.Assists,
			Points:           res.Points,
			PlusMinus:        res.PlusMinus,
			PenaltyMinutes:   res.PenaltyMinutes,
			HatTricks:        res.HatTricks,
			GameWinningGoals: res.GameWinningGoals,
		})
	}

	return tournaments, nil
}

// isChampionship проверяет является ли турнир первенством
func isChampionship(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.Contains(lowerName, "первенство")
}
