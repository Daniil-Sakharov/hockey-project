package player

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// GetByProfileURL возвращает игрока по URL профиля (для дедупликации)
func (r *repository) GetByProfileURL(ctx context.Context, url string) (*player.Player, error) {
	query := `
		SELECT 
			id, profile_url, name, birth_date, position,
			height, weight, handedness,
			data_season,
			source, created_at, updated_at
		FROM players
		WHERE profile_url = $1
	`

	var p player.Player
	err := r.db.GetContext(ctx, &p, query, url)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Игрок не найден - это нормально для дедупликации
		}
		return nil, fmt.Errorf("failed to get player by url: %w", err)
	}

	return &p, nil
}
