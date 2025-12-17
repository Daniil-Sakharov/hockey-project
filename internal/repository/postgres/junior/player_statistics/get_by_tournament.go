package player_statistics

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// GetByTournament возвращает всю статистику турнира
func (r *repository) GetByTournament(ctx context.Context, tournamentID string) ([]*player_statistics.PlayerStatistic, error) {
	query := `
		SELECT * FROM player_statistics 
		WHERE tournament_id = $1
		ORDER BY player_id, birth_year, group_name
	`

	var stats []*player_statistics.PlayerStatistic
	if err := r.db.SelectContext(ctx, &stats, query, tournamentID); err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return stats, nil
}
