package fhmoscow

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Tournament struct {
	ID         string     `db:"id"`
	ExternalID string     `db:"external_id"` // формат: "tournament_group" (14_137)
	URL        *string    `db:"url"`
	Name       string     `db:"name"`
	Domain     *string    `db:"domain"`
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
	id := fmt.Sprintf("fhm:%s", t.ExternalID)
	domain := "https://www.fhmoscow.com"

	query := `
		INSERT INTO tournaments (id, external_id, url, name, domain, birth_year, group_name, season, start_date, end_date, is_ended, region, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
		ON CONFLICT (id) DO UPDATE SET
			url = COALESCE(EXCLUDED.url, tournaments.url),
			name = EXCLUDED.name,
			domain = COALESCE(EXCLUDED.domain, tournaments.domain),
			birth_year = COALESCE(EXCLUDED.birth_year, tournaments.birth_year),
			group_name = COALESCE(EXCLUDED.group_name, tournaments.group_name),
			season = COALESCE(EXCLUDED.season, tournaments.season),
			start_date = COALESCE(EXCLUDED.start_date, tournaments.start_date),
			end_date = COALESCE(EXCLUDED.end_date, tournaments.end_date),
			is_ended = EXCLUDED.is_ended
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, t.ExternalID, t.URL, t.Name, domain, t.BirthYear, t.GroupName, t.Season, t.StartDate, t.EndDate, t.IsEnded, RegionMoscow, SourceFHMoscow).Scan(&returnedID)
	return returnedID, err
}

func (r *TournamentRepository) GetByExternalID(ctx context.Context, externalID string) (*Tournament, error) {
	var t Tournament
	id := fmt.Sprintf("fhm:%s", externalID)
	err := r.db.GetContext(ctx, &t, `SELECT id, external_id, url, name, domain, birth_year, group_name, season, start_date, end_date, is_ended, region, created_at FROM tournaments WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
