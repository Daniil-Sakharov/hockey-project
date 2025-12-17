package fhspb

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Player struct {
	ID         string     `db:"id"`
	ExternalID string     `db:"external_id"`
	FullName   string     `db:"name"`
	BirthDate  *time.Time `db:"birth_date"`
	BirthPlace *string    `db:"birth_place"`
	Source     string     `db:"source"`
	Region     string     `db:"region"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type PlayerRepository struct {
	db *sqlx.DB
}

func NewPlayerRepository(db *sqlx.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Upsert(ctx context.Context, p *Player) (string, error) {
	// ID формат: spb:<external_id>
	id := fmt.Sprintf("spb:%s", p.ExternalID)

	query := `
		INSERT INTO players (id, external_id, name, birth_date, birth_place, source, region, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			birth_date = COALESCE(EXCLUDED.birth_date, players.birth_date),
			birth_place = COALESCE(EXCLUDED.birth_place, players.birth_place),
			updated_at = NOW()
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, p.ExternalID, p.FullName, p.BirthDate, p.BirthPlace, SourceFHSPB, RegionSPB).Scan(&returnedID)
	return returnedID, err
}

func (r *PlayerRepository) GetByExternalID(ctx context.Context, externalID string) (*Player, error) {
	var p Player
	err := r.db.GetContext(ctx, &p, `SELECT id, external_id, name, birth_date, birth_place, source, region, created_at, updated_at FROM players WHERE external_id = $1 AND source = $2`, externalID, SourceFHSPB)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
