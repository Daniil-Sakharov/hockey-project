package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type TournamentPostgres struct {
	db *sqlx.DB
}

func NewTournamentPostgres(db *sqlx.DB) *TournamentPostgres {
	return &TournamentPostgres{db: db}
}

// Create создает новый турнир в БД
func (r *TournamentPostgres) Create(ctx context.Context, t *entities.Tournament) error {
	query := `
		INSERT INTO tournaments (
			id, url, name, domain, season, start_date, end_date, 
			is_ended, external_id, birth_year, group_name, region, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW())
		ON CONFLICT (id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.URL, t.Name, t.Domain, t.Season, t.StartDate, t.EndDate,
		t.IsEnded, t.ExternalID, t.BirthYear, t.GroupName, t.Region,
	)
	if err != nil {
		return fmt.Errorf("failed to create tournament: %w", err)
	}

	return nil
}

// CreateBatch создает множество турниров за один запрос
func (r *TournamentPostgres) CreateBatch(ctx context.Context, tournaments []*entities.Tournament) error {
	if len(tournaments) == 0 {
		return nil
	}

	query := `
		INSERT INTO tournaments (
			id, url, name, domain, season, start_date, end_date, 
			is_ended, external_id, birth_year, group_name, region, created_at
		) VALUES (
			:id, :url, :name, :domain, :season, :start_date, :end_date, 
			:is_ended, :external_id, :birth_year, :group_name, :region, NOW()
		)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := r.db.NamedExecContext(ctx, query, tournaments)
	if err != nil {
		return fmt.Errorf("failed to create tournaments batch: %w", err)
	}

	return nil
}

// Update обновляет турнир
func (r *TournamentPostgres) Update(ctx context.Context, t *entities.Tournament) error {
	query := `
		UPDATE tournaments SET
			name = $2, domain = $3, season = $4, start_date = $5, end_date = $6,
			is_ended = $7, external_id = $8, birth_year = $9, group_name = $10, region = $11
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.Name, t.Domain, t.Season, t.StartDate, t.EndDate,
		t.IsEnded, t.ExternalID, t.BirthYear, t.GroupName, t.Region,
	)
	if err != nil {
		return fmt.Errorf("failed to update tournament: %w", err)
	}

	return nil
}
