package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// GetByID получает игрока по ID
func (r *PlayerPostgres) GetByID(ctx context.Context, id string) (*entities.Player, error) {
	query := `SELECT * FROM players WHERE id = $1`

	var p entities.Player
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get player by ID: %w", err)
	}

	return &p, nil
}

// GetByProfileURL получает игрока по profile URL
func (r *PlayerPostgres) GetByProfileURL(ctx context.Context, profileURL string) (*entities.Player, error) {
	query := `SELECT * FROM players WHERE profile_url = $1`

	var p entities.Player
	err := r.db.GetContext(ctx, &p, query, profileURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get player by profile URL: %w", err)
	}

	return &p, nil
}

// GetByExternalID получает игрока по external ID и source
func (r *PlayerPostgres) GetByExternalID(ctx context.Context, externalID, source string) (*entities.Player, error) {
	query := `SELECT * FROM players WHERE external_id = $1 AND source = $2`

	var p entities.Player
	err := r.db.GetContext(ctx, &p, query, externalID, source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get player by external ID: %w", err)
	}

	return &p, nil
}

// ExistsByExternalID проверяет существование игрока по external ID
func (r *PlayerPostgres) ExistsByExternalID(ctx context.Context, externalID, source string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM players WHERE external_id = $1 AND source = $2)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, externalID, source)
	if err != nil {
		return false, fmt.Errorf("failed to check player exists: %w", err)
	}

	return exists, nil
}

// Update обновляет игрока
func (r *PlayerPostgres) Update(ctx context.Context, p *entities.Player) error {
	query := `
		UPDATE players SET
			name = $2, profile_url = $3, birth_date = $4, position = $5,
			height = $6, weight = $7, handedness = $8, birth_place = $9,
			citizenship = $10, role = $11, region = $12, updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.Name, p.ProfileURL, p.BirthDate, p.Position,
		p.Height, p.Weight, p.Handedness, p.BirthPlace,
		p.Citizenship, p.Role, p.Region,
	)
	if err != nil {
		return fmt.Errorf("failed to update player: %w", err)
	}

	return nil
}
