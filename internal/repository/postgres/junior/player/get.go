package player

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// GetByID возвращает игрока по ID
func (r *repository) GetByID(ctx context.Context, id string) (*player.Player, error) {
	query := `
		SELECT 
			id, profile_url, name, birth_date, position,
			height, weight, handedness,
			data_season,
			source, created_at, updated_at
		FROM players
		WHERE id = $1
	`

	var p player.Player
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return &p, nil
}
