package tournament

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// GetByID возвращает турнир по ID
func (r *repository) GetByID(ctx context.Context, id string) (*tournament.Tournament, error) {
	query := `
		SELECT id, url, name, domain, season, start_date, end_date, is_ended, created_at
		FROM tournaments
		WHERE id = $1
	`

	var t tournament.Tournament
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tournament with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get tournament: %w", err)
	}

	return &t, nil
}
