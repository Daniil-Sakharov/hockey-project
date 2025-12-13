package tournament

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// List возвращает список всех турниров
func (r *repository) List(ctx context.Context, limit, offset int) ([]*tournament.Tournament, error) {
	query := `
		SELECT id, url, name, domain, season, start_date, end_date, is_ended, created_at
		FROM tournaments
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	var tournaments []*tournament.Tournament
	err := r.db.SelectContext(ctx, &tournaments, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tournaments: %w", err)
	}

	return tournaments, nil
}
