package player

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// Create создает нового игрока в БД
func (r *repository) Create(ctx context.Context, p *player.Player) error {
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

	_, err := r.db.ExecContext(ctx, query,
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
		return fmt.Errorf("failed to create player: %w", err)
	}

	return nil
}
