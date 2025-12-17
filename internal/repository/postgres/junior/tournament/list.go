package tournament

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

func (r *repository) List(ctx context.Context) ([]*tournament.Tournament, error) {
	query := `
		SELECT id, url, name, domain, season, start_date, end_date, is_ended, created_at
		FROM tournaments
		WHERE url IS NOT NULL AND domain IS NOT NULL
		ORDER BY name
	`

	var tournaments []*tournament.Tournament
	err := r.db.SelectContext(ctx, &tournaments, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tournaments: %w", err)
	}

	return tournaments, nil
}
