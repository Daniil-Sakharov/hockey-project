package player

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// Update обновляет данные игрока
func (r *repository) Update(ctx context.Context, p *player.Player) error {
	query := `
		UPDATE players SET
			name = $1,
			birth_date = $2,
			position = $3,
			height = $4,
			weight = $5,
			handedness = $6,
			registry_id = $7,
			school = $8,
			rank = $9,
			data_season = $10,
			source = $11,
			updated_at = $12
		WHERE id = $13
	`

	p.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
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
		p.UpdatedAt,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("player with id %s not found", p.ID)
	}

	return nil
}
