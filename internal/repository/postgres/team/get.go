package team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// GetByID возвращает команду по ID
func (r *repository) GetByID(ctx context.Context, id string) (*team.Team, error) {
	query := `
		SELECT id, url, name, city, created_at
		FROM teams
		WHERE id = $1
	`

	var t team.Team
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("team with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &t, nil
}
