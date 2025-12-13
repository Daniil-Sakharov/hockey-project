package player_team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
)

// GetByPlayer возвращает все команды игрока
func (r *repository) GetByPlayer(ctx context.Context, playerID string) ([]*player_team.PlayerTeam, error) {
	query := `
		SELECT player_id, team_id, tournament_id, season, 
		       started_at, ended_at, is_active, 
		       jersey_number, role, source, 
		       created_at, updated_at
		FROM player_teams
		WHERE player_id = $1
		ORDER BY created_at DESC
	`

	var teams []*player_team.PlayerTeam
	if err := r.db.SelectContext(ctx, &teams, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get player teams: %w", err)
	}

	return teams, nil
}

// GetActiveByPlayer возвращает активные команды игрока
func (r *repository) GetActiveByPlayer(ctx context.Context, playerID string) ([]*player_team.PlayerTeam, error) {
	query := `
		SELECT player_id, team_id, tournament_id, season, 
		       started_at, ended_at, is_active, 
		       jersey_number, role, source, 
		       created_at, updated_at
		FROM player_teams
		WHERE player_id = $1 AND is_active = true AND ended_at IS NULL
		ORDER BY created_at DESC
	`

	var teams []*player_team.PlayerTeam
	if err := r.db.SelectContext(ctx, &teams, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get active player teams: %w", err)
	}

	return teams, nil
}

// GetByTeam возвращает всех игроков команды в турнире
func (r *repository) GetByTeam(ctx context.Context, teamID, tournamentID string) ([]*player_team.PlayerTeam, error) {
	query := `
		SELECT player_id, team_id, tournament_id, season, 
		       started_at, ended_at, is_active, 
		       jersey_number, role, source, 
		       created_at, updated_at
		FROM player_teams
		WHERE team_id = $1 AND tournament_id = $2
		ORDER BY created_at DESC
	`

	var players []*player_team.PlayerTeam
	if err := r.db.SelectContext(ctx, &players, query, teamID, tournamentID); err != nil {
		return nil, fmt.Errorf("failed to get team players: %w", err)
	}

	return players, nil
}
