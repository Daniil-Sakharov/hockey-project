package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
)

// CreateBatch создает несколько команд за одну транзакцию
func (r *repository) CreateBatch(ctx context.Context, teams []*team.Team) error {
	if len(teams) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO teams (id, url, name, city, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (url) DO NOTHING
	`

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, t := range teams {
		_, err := stmt.ExecContext(ctx,
			t.ID,
			t.URL,
			t.Name,
			t.City,
			t.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert team %s: %w", t.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
