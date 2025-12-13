package player_statistics

import (
	"context"
	"database/sql"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// GetSeasonStats возвращает агрегированную статистику игрока за конкретный сезон
func (r *repository) GetSeasonStats(ctx context.Context, playerID string, season string) (*player_statistics.AggregatedStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT ps.tournament_id) as tournaments_count,
			COALESCE(SUM(ps.games), 0) as games,
			COALESCE(SUM(ps.goals), 0) as goals,
			COALESCE(SUM(ps.assists), 0) as assists,
			COALESCE(SUM(ps.points), 0) as points,
			COALESCE(SUM(ps.plus), 0) as plus,
			COALESCE(SUM(ps.minus), 0) as minus,
			COALESCE(SUM(ps.plus_minus), 0) as plus_minus,
			COALESCE(SUM(ps.penalty_minutes), 0) as penalty_minutes,
			COALESCE(SUM(ps.hat_tricks), 0) as hat_tricks,
			COALESCE(SUM(ps.game_winning_goals), 0) as game_winning_goals
		FROM player_statistics ps
		INNER JOIN tournaments t ON ps.tournament_id = t.id
		WHERE ps.player_id = $1 AND t.season = $2
	`

	var stats player_statistics.AggregatedStats
	err := r.db.GetContext(ctx, &stats, query, playerID, season)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Если нет турниров, значит нет статистики
	if stats.TournamentsCount == 0 {
		return nil, nil
	}

	return &stats, nil
}
