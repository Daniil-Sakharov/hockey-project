package player_statistics

import (
	"context"
	"fmt"
)

// CountAll возвращает общее количество записей статистики
func (r *repository) CountAll(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM player_statistics`

	var count int
	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, fmt.Errorf("failed to count statistics: %w", err)
	}

	return count, nil
}
