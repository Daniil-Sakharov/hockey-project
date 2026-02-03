package mihf

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Player struct {
	ID          string     `db:"id"`
	ExternalID  string     `db:"external_id"`
	FullName    string     `db:"name"`
	ProfileURL  *string    `db:"profile_url"`
	BirthDate   *time.Time `db:"birth_date"`
	BirthPlace  *string    `db:"birth_place"`
	Position    *string    `db:"position"`
	Height      *int       `db:"height"`
	Weight      *int       `db:"weight"`
	Handedness  *string    `db:"handedness"`
	Citizenship *string    `db:"citizenship"`
	Source      string     `db:"source"`
	Region      string     `db:"region"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type PlayerRepository struct {
	db *sqlx.DB
}

func NewPlayerRepository(db *sqlx.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) Upsert(ctx context.Context, p *Player) (string, error) {
	id := fmt.Sprintf("msk:%s", p.ExternalID)

	query := `
		INSERT INTO players (id, external_id, name, profile_url, birth_date, birth_place, position, height, weight, handedness, citizenship, source, region, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			profile_url = COALESCE(EXCLUDED.profile_url, players.profile_url),
			birth_date = COALESCE(EXCLUDED.birth_date, players.birth_date),
			birth_place = COALESCE(EXCLUDED.birth_place, players.birth_place),
			position = COALESCE(EXCLUDED.position, players.position),
			height = COALESCE(EXCLUDED.height, players.height),
			weight = COALESCE(EXCLUDED.weight, players.weight),
			handedness = COALESCE(EXCLUDED.handedness, players.handedness),
			citizenship = COALESCE(EXCLUDED.citizenship, players.citizenship),
			updated_at = NOW()
		RETURNING id`

	var returnedID string
	err := r.db.QueryRowContext(ctx, query, id, p.ExternalID, p.FullName, p.ProfileURL, p.BirthDate, p.BirthPlace, p.Position, p.Height, p.Weight, p.Handedness, p.Citizenship, SourceMIHF, RegionMoscow).Scan(&returnedID)
	return returnedID, err
}

func (r *PlayerRepository) GetByExternalID(ctx context.Context, externalID string) (*Player, error) {
	var p Player
	err := r.db.GetContext(ctx, &p, `SELECT id, external_id, name, profile_url, birth_date, birth_place, position, height, weight, handedness, citizenship, source, region, created_at, updated_at FROM players WHERE external_id = $1 AND source = $2`, externalID, SourceMIHF)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PlayerRepository) Exists(ctx context.Context, externalID string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM players WHERE external_id = $1 AND source = $2)`, externalID, SourceMIHF)
	return exists, err
}
