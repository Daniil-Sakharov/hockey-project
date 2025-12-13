package team

import (
	"context"
	"database/sql"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// GetByURL возвращает команду по URL (для дедупликации)
func (r *repository) GetByURL(ctx context.Context, url string) (*team.Team, error) {
	query := `
		SELECT id, url, name, city, created_at
		FROM teams
		WHERE url = $1
	`

	var t team.Team
	err := r.db.GetContext(ctx, &t, query, url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Команда не найдена - это нормально для дедупликации
		}
		return nil, err
	}

	return &t, nil
}
