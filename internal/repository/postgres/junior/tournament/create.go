package tournament

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// Create создает новый турнир в БД
func (r *repository) Create(ctx context.Context, t *tournament.Tournament) error {
	query := `
		INSERT INTO tournaments (id, url, name, domain, season, start_date, end_date, is_ended, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (url) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID,
		t.URL,
		t.Name,
		t.Domain,
		t.Season,
		t.StartDate,
		t.EndDate,
		t.IsEnded,
		t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create tournament: %w", err)
	}

	return nil
}
