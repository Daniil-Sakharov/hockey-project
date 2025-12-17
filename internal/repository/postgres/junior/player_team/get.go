package player_team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
)

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

func (r *repository) GetByTeam(ctx context.Context, teamID string) ([]*player_team.PlayerTeam, error) {
	query := `
		SELECT player_id, team_id, tournament_id, season, 
		       started_at, ended_at, is_active, 
		       jersey_number, role, source, 
		       created_at, updated_at
		FROM player_teams
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	var players []*player_team.PlayerTeam
	if err := r.db.SelectContext(ctx, &players, query, teamID); err != nil {
		return nil, fmt.Errorf("failed to get team players: %w", err)
	}

	return players, nil
}

func (r *repository) GetByTournament(ctx context.Context, tournamentID string) ([]*player_team.PlayerTeam, error) {
	query := `
		SELECT player_id, team_id, tournament_id, season, 
		       started_at, ended_at, is_active, 
		       jersey_number, role, source, 
		       created_at, updated_at
		FROM player_teams
		WHERE tournament_id = $1
		ORDER BY created_at DESC
	`

	var players []*player_team.PlayerTeam
	if err := r.db.SelectContext(ctx, &players, query, tournamentID); err != nil {
		return nil, fmt.Errorf("failed to get tournament players: %w", err)
	}

	return players, nil
}
