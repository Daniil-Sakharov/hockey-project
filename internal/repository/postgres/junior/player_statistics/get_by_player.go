package player_statistics

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// GetByPlayerID возвращает всю статистику игрока
func (r *repository) GetByPlayerID(ctx context.Context, playerID string) ([]*player_statistics.PlayerStatistic, error) {
	query := `
		SELECT * FROM player_statistics 
		WHERE player_id = $1
		ORDER BY tournament_id, birth_year, group_name
	`

	var stats []*player_statistics.PlayerStatistic
	if err := r.db.SelectContext(ctx, &stats, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return stats, nil
}
