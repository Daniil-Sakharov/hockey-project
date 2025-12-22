package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type playerTeamRepository struct {
	db *sqlx.DB
}

// NewPlayerTeamRepository создает новый репозиторий для player_teams
func NewPlayerTeamRepository(db *sqlx.DB) *playerTeamRepository {
	return &playerTeamRepository{db: db}
}

func (r *playerTeamRepository) Create(ctx context.Context, link *entities.PlayerTeam) error {
	return r.Upsert(ctx, link)
}

func (r *playerTeamRepository) CreateBatch(ctx context.Context, links []*entities.PlayerTeam) error {
	for _, link := range links {
		if err := r.Upsert(ctx, link); err != nil {
			return err
		}
	}
	return nil
}

func (r *playerTeamRepository) Upsert(ctx context.Context, pt *entities.PlayerTeam) error {
	query := `
		INSERT INTO player_teams (
			player_id, team_id, tournament_id, season, 
			started_at, ended_at, is_active, 
			jersey_number, role, source, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (player_id, team_id, tournament_id) 
		DO UPDATE SET 
			season = EXCLUDED.season,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		pt.PlayerID, pt.TeamID, pt.TournamentID, pt.Season,
		pt.StartedAt, pt.EndedAt, pt.IsActive,
		pt.JerseyNumber, pt.Role, pt.Source,
		pt.CreatedAt, pt.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert player_team: %w", err)
	}
	return nil
}

func (r *playerTeamRepository) GetByPlayerID(ctx context.Context, playerID string) ([]*entities.PlayerTeam, error) {
	var links []*entities.PlayerTeam
	query := `SELECT * FROM player_teams WHERE player_id = $1`
	if err := r.db.SelectContext(ctx, &links, query, playerID); err != nil {
		return nil, fmt.Errorf("failed to get player_teams by player_id: %w", err)
	}
	return links, nil
}

func (r *playerTeamRepository) GetByTeamID(ctx context.Context, teamID string) ([]*entities.PlayerTeam, error) {
	var links []*entities.PlayerTeam
	query := `SELECT * FROM player_teams WHERE team_id = $1`
	if err := r.db.SelectContext(ctx, &links, query, teamID); err != nil {
		return nil, fmt.Errorf("failed to get player_teams by team_id: %w", err)
	}
	return links, nil
}

func (r *playerTeamRepository) Exists(ctx context.Context, playerID, teamID, tournamentID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM player_teams WHERE player_id = $1 AND team_id = $2 AND tournament_id = $3`
	if err := r.db.GetContext(ctx, &count, query, playerID, teamID, tournamentID); err != nil {
		return false, fmt.Errorf("failed to check player_team exists: %w", err)
	}
	return count > 0, nil
}
