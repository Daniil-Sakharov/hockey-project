package tournament

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// Update обновляет турнир
func (r *repository) Update(ctx context.Context, t *tournament.Tournament) error {
	query := `
		UPDATE tournaments 
		SET name = $1,
		    domain = $2,
		    season = $3,
		    start_date = $4,
		    end_date = $5,
		    is_ended = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		t.Name,
		t.Domain,
		t.Season,
		t.StartDate,
		t.EndDate,
		t.IsEnded,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update tournament: %w", err)
	}

	return nil
}
