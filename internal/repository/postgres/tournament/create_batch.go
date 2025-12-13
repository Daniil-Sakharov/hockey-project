package tournament

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// CreateBatch создает несколько турниров за одну транзакцию
func (r *repository) CreateBatch(ctx context.Context, tournaments []*tournament.Tournament) error {
	if len(tournaments) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tournaments (id, url, name, domain, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (url) DO NOTHING
	`

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, t := range tournaments {
		_, err := stmt.ExecContext(ctx,
			t.ID,
			t.URL,
			t.Name,
			t.Domain,
			t.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert tournament %s: %w", t.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
