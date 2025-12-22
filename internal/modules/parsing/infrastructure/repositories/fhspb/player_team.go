package fhspb

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlayerTeam struct {
	PlayerID     string     `db:"player_id"`
	TeamID       string     `db:"team_id"`
	TournamentID string     `db:"tournament_id"`
	Season       *string    `db:"season"`
	StartedAt    *time.Time `db:"started_at"`
	EndedAt      *time.Time `db:"ended_at"`
	IsActive     *bool      `db:"is_active"`
	Number       *int       `db:"jersey_number"`
	Role         *string    `db:"role"`
	Position     *string    `db:"position"`
	Source       string     `db:"source"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

type PlayerTeamRepository struct {
	db *sqlx.DB
}

func NewPlayerTeamRepository(db *sqlx.DB) *PlayerTeamRepository {
	return &PlayerTeamRepository{db: db}
}

func (r *PlayerTeamRepository) Upsert(ctx context.Context, pt *PlayerTeam) error {
	query := `
		INSERT INTO player_teams (player_id, team_id, tournament_id, season, started_at, ended_at, is_active, jersey_number, role, position, source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		ON CONFLICT (player_id, team_id, tournament_id) DO UPDATE SET
			season = COALESCE(EXCLUDED.season, player_teams.season),
			started_at = COALESCE(EXCLUDED.started_at, player_teams.started_at),
			ended_at = COALESCE(EXCLUDED.ended_at, player_teams.ended_at),
			is_active = COALESCE(EXCLUDED.is_active, player_teams.is_active),
			jersey_number = EXCLUDED.jersey_number,
			role = EXCLUDED.role,
			position = EXCLUDED.position,
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query, pt.PlayerID, pt.TeamID, pt.TournamentID, pt.Season, pt.StartedAt, pt.EndedAt, pt.IsActive, pt.Number, pt.Role, pt.Position, SourceFHSPB)
	return err
}
