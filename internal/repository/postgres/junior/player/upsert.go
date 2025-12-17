package player

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// Upsert создает или обновляет игрока по external_id и source
func (r *repository) Upsert(ctx context.Context, p *player.Player) error {
	query := `
		INSERT INTO players (
			id, profile_url, name, birth_date, position,
			height, weight, handedness,
			registry_id, school, rank, data_season,
			external_id, citizenship, role, birth_place,
			source, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19
		)
		ON CONFLICT (external_id, source) WHERE external_id IS NOT NULL
		DO UPDATE SET
			name = EXCLUDED.name,
			birth_date = EXCLUDED.birth_date,
			position = EXCLUDED.position,
			height = EXCLUDED.height,
			weight = EXCLUDED.weight,
			handedness = EXCLUDED.handedness,
			school = EXCLUDED.school,
			citizenship = EXCLUDED.citizenship,
			role = EXCLUDED.role,
			birth_place = EXCLUDED.birth_place,
			updated_at = EXCLUDED.updated_at
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
		p.ExternalID,
		p.Citizenship,
		p.Role,
		p.BirthPlace,
		p.Source,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert player: %w", err)
	}

	return nil
}
