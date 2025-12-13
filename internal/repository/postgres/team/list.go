package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

func (r *repository) List(ctx context.Context) ([]*team.Team, error) {
	query := `
		SELECT id, url, name, city, created_at
		FROM teams
		ORDER BY name
	`

	var teams []*team.Team
	err := r.db.SelectContext(ctx, &teams, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	return teams, nil
}
