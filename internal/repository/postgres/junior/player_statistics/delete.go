package player_statistics

import (
	"context"
	"fmt"
)

// DeleteByTournament удаляет всю статистику турнира
func (r *repository) DeleteByTournament(ctx context.Context, tournamentID string) error {
	query := `DELETE FROM player_statistics WHERE tournament_id = $1`
	_, err := r.db.ExecContext(ctx, query, tournamentID)
	if err != nil {
		return fmt.Errorf("failed to delete statistics: %w", err)
	}
	return nil
}

// DeleteAll удаляет всю статистику (TRUNCATE)
func (r *repository) DeleteAll(ctx context.Context) error {
	query := `TRUNCATE TABLE player_statistics`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate table: %w", err)
	}
	return nil
}
