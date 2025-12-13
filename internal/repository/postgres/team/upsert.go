package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// Upsert создает или обновляет команду (ON CONFLICT DO UPDATE)
func (r *repository) Upsert(ctx context.Context, t *team.Team) (*team.Team, error) {
	query := `
		INSERT INTO teams (id, url, name, city, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) 
		DO UPDATE SET 
			url = EXCLUDED.url,
			name = EXCLUDED.name,
			city = EXCLUDED.city
		RETURNING id, url, name, city, created_at
	`

	var result team.Team
	err := r.db.QueryRowContext(ctx, query,
		t.ID,
		t.URL,
		t.Name,
		t.City,
		t.CreatedAt,
	).Scan(
		&result.ID,
		&result.URL,
		&result.Name,
		&result.City,
		&result.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert team: %w", err)
	}

	return &result, nil
}
