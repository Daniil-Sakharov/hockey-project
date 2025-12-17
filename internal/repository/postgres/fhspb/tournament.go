package fhspb

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	RegionSPB   = "Санкт-Петербург"
	SourceFHSPB = "fhspb.ru"
)

type Tournament struct {
	ID         string     `db:"id"`
	ExternalID string     `db:"external_id"`
	Name       string     `db:"name"`
	BirthYear  *int       `db:"birth_year"`
	GroupName  *string    `db:"group_name"`
	Season     *string    `db:"season"`
	StartDate  *time.Time `db:"start_date"`
	EndDate    *time.Time `db:"end_date"`
	IsEnded    bool       `db:"is_ended"`
	Region     string     `db:"region"`
	CreatedAt  time.Time  `db:"created_at"`
}

type TournamentRepository struct {
	db *sqlx.DB
}

func NewTournamentRepository(db *sqlx.DB) *TournamentRepository {
	return &TournamentRepository{db: db}
}

func (r *TournamentRepository) Upsert(ctx context.Context, t *Tournament) (string, error) {
	// ID формат: spb:<external_id>
	id := fmt.Sprintf("spb:%s", t.ExternalID)

	query := `
		INSERT INTO tournaments (id, external_id, name, birth_year, group_name, season, start_date, end_date, is_ended, region, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			birth_year = COALESCE(EXCLUDED.birth_year, tournaments.birth_year),
			group_name = COALESCE(EXCLUDED.group_name, tournaments.group_name),
			season = COALESCE(EXCLUDED.season, tournaments.season),
			start_date = COALESCE(EXCLUDED.start_date, tournaments.start_date),
			end_date = COALESCE(EXCLUDED.end_date, tournaments.end_date),
			is_ended = EXCLUDED.is_ended
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, t.ExternalID, t.Name, t.BirthYear, t.GroupName, t.Season, t.StartDate, t.EndDate, t.IsEnded, RegionSPB).Scan(&returnedID)
	return returnedID, err
}

func (r *TournamentRepository) GetByExternalID(ctx context.Context, externalID string) (*Tournament, error) {
	var t Tournament
	err := r.db.GetContext(ctx, &t, `SELECT id, external_id, name, birth_year, group_name, season, start_date, end_date, is_ended, region, created_at FROM tournaments WHERE external_id = $1 AND region = $2`, externalID, RegionSPB)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TournamentRepository) GetAll(ctx context.Context) ([]Tournament, error) {
	var tournaments []Tournament
	err := r.db.SelectContext(ctx, &tournaments, `SELECT id, external_id, name, birth_year, group_name, season, start_date, end_date, is_ended, region, created_at FROM tournaments WHERE region = $1 ORDER BY id`, RegionSPB)
	return tournaments, err
}
