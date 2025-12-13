package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// Create создает новую команду в БД
func (r *repository) Create(ctx context.Context, t *team.Team) error {
	query := `
		INSERT INTO teams (id, url, name, city, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (url) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.URL,
		t.Name,
		t.City,
		t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}
