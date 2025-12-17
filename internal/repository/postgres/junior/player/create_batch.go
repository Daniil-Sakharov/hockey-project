package player

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// CreateBatch создает несколько игроков за одну транзакцию
func (r *repository) CreateBatch(ctx context.Context, players []*player.Player) error {
	if len(players) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO players (
			id, profile_url, name, birth_date, position, 
			height, weight, handedness, 
			registry_id, school, rank, data_season,
			source, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, 
			$9, $10, $11, $12,
			$13, $14, $15
		)
		ON CONFLICT (profile_url) DO NOTHING
	`

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, p := range players {
		_, err := stmt.ExecContext(ctx,
			p.ID,
			p.ProfileURL,
			p.Name,
			p.BirthDate,
			p.Position,
			p.Height,
			p.Weight,
			p.Handedness,
			p.RegistryID,
			p.School,
			p.Rank,
			p.DataSeason,
			p.Source,
			p.CreatedAt,
			p.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert player %s: %w", p.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
