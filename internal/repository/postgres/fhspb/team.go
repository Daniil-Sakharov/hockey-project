package fhspb

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Team struct {
	ID           string    `db:"id"`
	ExternalID   string    `db:"external_id"`
	TournamentID string    `db:"tournament_id"`
	Name         string    `db:"name"`
	Region       string    `db:"region"`
	CreatedAt    time.Time `db:"created_at"`
}

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Upsert(ctx context.Context, t *Team) (string, error) {
	// ID формат: spb:<tournament_external_id>:<team_external_id>
	// Извлекаем tournament_external_id из tournament_id (spb:6366 -> 6366)
	tournamentExtID := t.TournamentID
	if len(tournamentExtID) > 4 && tournamentExtID[:4] == "spb:" {
		tournamentExtID = tournamentExtID[4:]
	}
	id := fmt.Sprintf("spb:%s:%s", tournamentExtID, t.ExternalID)

	query := `
		INSERT INTO teams (id, external_id, tournament_id, name, region, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, t.ExternalID, t.TournamentID, t.Name, RegionSPB).Scan(&returnedID)
	return returnedID, err
}

func (r *TeamRepository) GetByExternalID(ctx context.Context, externalID, tournamentID string) (*Team, error) {
	var t Team
	err := r.db.GetContext(ctx, &t, `SELECT id, external_id, tournament_id, name, region, created_at FROM teams WHERE external_id = $1 AND tournament_id = $2`, externalID, tournamentID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
