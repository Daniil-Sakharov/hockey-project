package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

func (r *repository) Upsert(ctx context.Context, t *team.Team) error {
	query := `
		INSERT INTO teams (id, url, name, city, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) 
		DO UPDATE SET 
			url = EXCLUDED.url,
			name = EXCLUDED.name,
			city = EXCLUDED.city
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.URL,
		t.Name,
		t.City,
		t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert team: %w", err)
	}

	return nil
}
