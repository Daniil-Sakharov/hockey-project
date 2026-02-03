package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type TeamPostgres struct {
	db *sqlx.DB
}

func NewTeamPostgres(db *sqlx.DB) *TeamPostgres {
	return &TeamPostgres{db: db}
}

func (r *TeamPostgres) Create(ctx context.Context, t *entities.Team) error {
	query := `
		INSERT INTO teams (id, url, name, city, external_id, tournament_id, region, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.URL, t.Name, t.City, t.ExternalID, t.TournamentID, t.Region, t.Source,
	)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}

func (r *TeamPostgres) CreateBatch(ctx context.Context, teams []*entities.Team) error {
	if len(teams) == 0 {
		return nil
	}

	query := `
		INSERT INTO teams (id, url, name, city, external_id, tournament_id, region, source, created_at)
		VALUES (:id, :url, :name, :city, :external_id, :tournament_id, :region, :source, NOW())
		ON CONFLICT (id) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, query, teams)
	if err != nil {
		return fmt.Errorf("failed to create teams batch: %w", err)
	}
	return nil
}

func (r *TeamPostgres) Upsert(ctx context.Context, t *entities.Team) error {
	query := `
		INSERT INTO teams (id, external_id, url, name, city, logo_url, tournament_id, region, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			city = COALESCE(EXCLUDED.city, teams.city),
			url = COALESCE(EXCLUDED.url, teams.url),
			logo_url = COALESCE(EXCLUDED.logo_url, teams.logo_url),
			tournament_id = COALESCE(EXCLUDED.tournament_id, teams.tournament_id),
			region = COALESCE(EXCLUDED.region, teams.region)`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.ExternalID, t.URL, t.Name, t.City, t.LogoURL, t.TournamentID, t.Region, t.Source,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert team: %w", err)
	}
	return nil
}

func (r *TeamPostgres) GetByID(ctx context.Context, id string) (*entities.Team, error) {
	query := `SELECT * FROM teams WHERE id = $1`

	var t entities.Team
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team by ID: %w", err)
	}
	return &t, nil
}

func (r *TeamPostgres) GetByURL(ctx context.Context, url string) (*entities.Team, error) {
	query := `SELECT * FROM teams WHERE url = $1`

	var t entities.Team
	err := r.db.GetContext(ctx, &t, query, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team by URL: %w", err)
	}
	return &t, nil
}

func (r *TeamPostgres) GetByName(ctx context.Context, name, source string) (*entities.Team, error) {
	query := `SELECT * FROM teams WHERE name = $1 AND source = $2 LIMIT 1`

	var t entities.Team
	err := r.db.GetContext(ctx, &t, query, name, source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team by name: %w", err)
	}
	return &t, nil
}

func (r *TeamPostgres) GetByNameAndCity(ctx context.Context, name, city, source string) (*entities.Team, error) {
	query := `SELECT * FROM teams WHERE name = $1 AND city = $2 AND source = $3`

	var t entities.Team
	err := r.db.GetContext(ctx, &t, query, name, city, source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team by name and city: %w", err)
	}
	return &t, nil
}

func (r *TeamPostgres) List(ctx context.Context) ([]*entities.Team, error) {
	query := `SELECT * FROM teams`

	var teams []*entities.Team
	err := r.db.SelectContext(ctx, &teams, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	return teams, nil
}

func (r *TeamPostgres) UpdateCity(ctx context.Context, id, city string) error {
	if city == "" {
		return nil
	}
	query := `UPDATE teams SET city = $2 WHERE id = $1 AND (city IS NULL OR city = '')`
	_, err := r.db.ExecContext(ctx, query, id, city)
	if err != nil {
		return fmt.Errorf("failed to update team city: %w", err)
	}
	return nil
}

func (r *TeamPostgres) UpdateLogoURL(ctx context.Context, id, logoURL string) error {
	if logoURL == "" {
		return nil
	}
	query := `UPDATE teams SET logo_url = $2 WHERE id = $1 AND (logo_url IS NULL OR logo_url = '')`
	_, err := r.db.ExecContext(ctx, query, id, logoURL)
	if err != nil {
		return fmt.Errorf("failed to update team logo URL: %w", err)
	}
	return nil
}
