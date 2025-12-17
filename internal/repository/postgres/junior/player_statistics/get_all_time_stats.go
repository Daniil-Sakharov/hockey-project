package player_statistics

import (
	"context"
	"database/sql"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// GetAllTimeStats возвращает агрегированную статистику игрока за всё время
func (r *repository) GetAllTimeStats(ctx context.Context, playerID string) (*player_statistics.AggregatedStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT tournament_id) as tournaments_count,
			COALESCE(SUM(games), 0) as games,
			COALESCE(SUM(goals), 0) as goals,
			COALESCE(SUM(assists), 0) as assists,
			COALESCE(SUM(points), 0) as points,
			COALESCE(SUM(plus), 0) as plus,
			COALESCE(SUM(minus), 0) as minus,
			COALESCE(SUM(plus_minus), 0) as plus_minus,
			COALESCE(SUM(penalty_minutes), 0) as penalty_minutes,
			COALESCE(SUM(hat_tricks), 0) as hat_tricks,
			COALESCE(SUM(game_winning_goals), 0) as game_winning_goals
		FROM player_statistics
		WHERE player_id = $1
	`

	var stats player_statistics.AggregatedStats
	err := r.db.GetContext(ctx, &stats, query, playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Нет статистики
		}
		return nil, err
	}

	// Если нет турниров, значит нет статистики
	if stats.TournamentsCount == 0 {
		return nil, nil
	}

	return &stats, nil
}
