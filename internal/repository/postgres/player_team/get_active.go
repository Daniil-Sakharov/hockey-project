package player_team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
)

func (r *repository) GetActiveTeam(ctx context.Context, playerID string) (*player_team.PlayerTeam, error) {
	query := `
		SELECT player_id, team_id, tournament_id, season, started_at, ended_at, 
		       is_active, jersey_number, role, source, created_at, updated_at
		FROM player_teams
		WHERE player_id = $1 AND is_active = true
		ORDER BY started_at DESC NULLS LAST
		LIMIT 1
	`

	var pt player_team.PlayerTeam
	err := r.db.GetContext(ctx, &pt, query, playerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active team: %w", err)
	}

	return &pt, nil
}
