package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type PlayerPostgres struct {
	db *sqlx.DB
}

func NewPlayerPostgres(db *sqlx.DB) *PlayerPostgres {
	return &PlayerPostgres{db: db}
}

// Create создает нового игрока в БД
func (r *PlayerPostgres) Create(ctx context.Context, p *entities.Player) error {
	query := `
		INSERT INTO players (
			id, profile_url, name, birth_date, position, 
			height, weight, handedness, 
			data_season, external_id, birth_place, citizenship, role, region,
			source, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.ProfileURL, p.Name, p.BirthDate, p.Position,
		p.Height, p.Weight, p.Handedness,
		p.DataSeason, p.ExternalID, p.BirthPlace, p.Citizenship, p.Role, p.Region,
		p.Source, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	return nil
}

// CreateBatch создает множество игроков за один запрос
func (r *PlayerPostgres) CreateBatch(ctx context.Context, players []*entities.Player) error {
	if len(players) == 0 {
		return nil
	}

	query := `
		INSERT INTO players (
			id, profile_url, name, birth_date, position, 
			height, weight, handedness, 
			data_season, external_id, birth_place, citizenship, role, region,
			source, created_at, updated_at
		) VALUES (
			:id, :profile_url, :name, :birth_date, :position, 
			:height, :weight, :handedness, 
			:data_season, :external_id, :birth_place, :citizenship, :role, :region,
			:source, :created_at, :updated_at
		)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := r.db.NamedExecContext(ctx, query, players)
	if err != nil {
		return fmt.Errorf("failed to create players batch: %w", err)
	}

	return nil
}

// Upsert создает или обновляет игрока
func (r *PlayerPostgres) Upsert(ctx context.Context, p *entities.Player) error {
	query := `
		INSERT INTO players (
			id, external_id, name, profile_url, birth_date, birth_place, 
			position, height, weight, handedness, citizenship, role, region,
			source, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
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
			role = COALESCE(EXCLUDED.role, players.role),
			region = COALESCE(EXCLUDED.region, players.region),
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.ExternalID, p.Name, p.ProfileURL, p.BirthDate, p.BirthPlace,
		p.Position, p.Height, p.Weight, p.Handedness, p.Citizenship, p.Role, p.Region, p.Source,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert player: %w", err)
	}

	return nil
}
