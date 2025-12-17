package tournament

import (
	"context"
	"database/sql"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// GetByURL возвращает турнир по URL (для дедупликации)
func (r *repository) GetByURL(ctx context.Context, url string) (*tournament.Tournament, error) {
	query := `
		SELECT id, url, name, domain, season, start_date, end_date, is_ended, created_at
		FROM tournaments
		WHERE url = $1
	`

	var t tournament.Tournament
	err := r.db.GetContext(ctx, &t, query, url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Турнир не найден - это нормально для дедупликации
		}
		return nil, err
	}

	return &t, nil
}
