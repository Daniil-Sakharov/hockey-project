package fhmoscow

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
	// ID формат: fhm:teamExtID (команда уникальна по external_id + source)
	id := fmt.Sprintf("fhm:%s", t.ExternalID)

	query := `
		INSERT INTO teams (id, external_id, tournament_id, name, url, city, region, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (external_id, source) WHERE external_id IS NOT NULL DO UPDATE SET
			name = EXCLUDED.name,
			url = COALESCE(EXCLUDED.url, teams.url),
			city = COALESCE(EXCLUDED.city, teams.city),
			tournament_id = EXCLUDED.tournament_id
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, t.ExternalID, t.TournamentID, t.Name, t.URL, t.City, RegionMoscow, SourceFHMoscow).Scan(&returnedID)
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

// FindByExternalIDAndTournamentName ищет команду по external_id и названию турнира
// Возвращает команду с привязкой к турниру с группой (если есть)
func (r *TeamRepository) FindByExternalIDAndTournamentName(ctx context.Context, externalID, tournamentName string) (*Team, error) {
	var t Team
	// Ищем команду где турнир имеет такое же название, предпочитая турниры с группами
	query := `
		SELECT tm.id, tm.external_id, tm.tournament_id, tm.name, tm.url, tm.city, tm.region, tm.created_at
		FROM teams tm
		JOIN tournaments t ON tm.tournament_id = t.id
		WHERE tm.external_id = $1
		  AND tm.source = $2
		  AND t.name = $3
		ORDER BY t.group_name NULLS LAST
		LIMIT 1`
	err := r.db.GetContext(ctx, &t, query, externalID, SourceFHMoscow, tournamentName)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// FindByNameAndTournamentName ищет команду по названию команды и названию турнира
// Возвращает команду с привязкой к турниру с группой (если есть)
func (r *TeamRepository) FindByNameAndTournamentName(ctx context.Context, teamName, tournamentName string) (*Team, error) {
	var t Team
	query := `
		SELECT tm.id, tm.external_id, tm.tournament_id, tm.name, tm.url, tm.city, tm.region, tm.created_at
		FROM teams tm
		JOIN tournaments t ON tm.tournament_id = t.id
		WHERE tm.name = $1
		  AND tm.source = $2
		  AND t.name = $3
		  AND t.group_name IS NOT NULL
		ORDER BY t.group_name
		LIMIT 1`
	err := r.db.GetContext(ctx, &t, query, teamName, SourceFHMoscow, tournamentName)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
