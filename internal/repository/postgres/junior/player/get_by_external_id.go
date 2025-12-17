package player

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// GetByExternalID возвращает игрока по внешнему ID и источнику
func (r *repository) GetByExternalID(ctx context.Context, externalID, source string) (*player.Player, error) {
	query := `
		SELECT id, profile_url, name, birth_date, position,
			   height, weight, handedness,
			   registry_id, school, rank, data_season,
			   external_id, citizenship, role, birth_place,
			   source, created_at, updated_at
		FROM players
		WHERE external_id = $1 AND source = $2
	`

	var p player.Player
	err := r.db.GetContext(ctx, &p, query, externalID, source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get player by external_id: %w", err)
	}

	return &p, nil
}

// ExistsByExternalID проверяет существование игрока по внешнему ID и источнику
func (r *repository) ExistsByExternalID(ctx context.Context, externalID, source string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM players WHERE external_id = $1 AND source = $2)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, externalID, source)
	if err != nil {
		return false, fmt.Errorf("failed to check player existence: %w", err)
	}

	return exists, nil
}
