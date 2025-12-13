package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// List возвращает список всех команд
func (r *repository) List(ctx context.Context, limit, offset int) ([]*team.Team, error) {
	query := `
		SELECT id, url, name, city, created_at
		FROM teams
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	var teams []*team.Team
	err := r.db.SelectContext(ctx, &teams, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	return teams, nil
}
