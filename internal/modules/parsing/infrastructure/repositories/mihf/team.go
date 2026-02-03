package mihf

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
	tournamentExtID := t.TournamentID
	if len(tournamentExtID) > 4 && tournamentExtID[:4] == "msk:" {
		tournamentExtID = tournamentExtID[4:]
	}
	id := fmt.Sprintf("msk:%s:%s", tournamentExtID, t.ExternalID)

	query := `
		INSERT INTO teams (id, external_id, tournament_id, name, url, city, region, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			url = COALESCE(EXCLUDED.url, teams.url),
			city = COALESCE(EXCLUDED.city, teams.city)
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, t.ExternalID, t.TournamentID, t.Name, t.URL, t.City, RegionMoscow, SourceMIHF).Scan(&returnedID)
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

// UpdateCity обновляет город команды
func (r *TeamRepository) UpdateCity(ctx context.Context, id, city string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE teams SET city = $1 WHERE id = $2`, city, id)
	return err
}

// UpdateLogoURL обновляет URL логотипа команды
func (r *TeamRepository) UpdateLogoURL(ctx context.Context, id, logoURL string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE teams SET logo_url = $1 WHERE id = $2`, logoURL, id)
	return err
}
