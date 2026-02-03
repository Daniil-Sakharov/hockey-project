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
	URL          *string   `db:"url"`
	City         *string   `db:"city"`
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
	// ID формируется как spb:external_id (команда уникальна по external_id + source)
	id := fmt.Sprintf("spb:%s", t.ExternalID)

	query := `
		INSERT INTO teams (id, external_id, tournament_id, name, url, city, region, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (external_id, source) WHERE external_id IS NOT NULL DO UPDATE SET
			name = EXCLUDED.name,
			url = COALESCE(EXCLUDED.url, teams.url),
			city = COALESCE(EXCLUDED.city, teams.city)
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, t.ExternalID, t.TournamentID, t.Name, t.URL, t.City, RegionSPB, SourceFHSPB).Scan(&returnedID)
	return returnedID, err
}

func (r *TeamRepository) GetByExternalID(ctx context.Context, externalID, tournamentID string) (*Team, error) {
	var t Team
	err := r.db.GetContext(ctx, &t, `SELECT id, external_id, tournament_id, name, url, city, region, created_at FROM teams WHERE external_id = $1 AND tournament_id = $2`, externalID, tournamentID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeamRepository) GetByName(ctx context.Context, name, tournamentID string) (*Team, error) {
	var t Team
	err := r.db.GetContext(ctx, &t, `SELECT id, external_id, tournament_id, name, url, city, region, created_at FROM teams WHERE name = $1 AND tournament_id = $2 AND source = $3`, name, tournamentID, SourceFHSPB)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
