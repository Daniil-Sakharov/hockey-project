package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// GetByID получает турнир по ID
func (r *TournamentPostgres) GetByID(ctx context.Context, id string) (*entities.Tournament, error) {
	query := `SELECT * FROM tournaments WHERE id = $1`

	var t entities.Tournament
	err := r.db.GetContext(ctx, &t, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tournament by ID: %w", err)
	}
	return &t, nil
}

// GetByURL получает турнир по URL
func (r *TournamentPostgres) GetByURL(ctx context.Context, url string) (*entities.Tournament, error) {
	query := `SELECT * FROM tournaments WHERE url = $1`

	var t entities.Tournament
	err := r.db.GetContext(ctx, &t, query, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tournament by URL: %w", err)
	}
	return &t, nil
}

// GetBySource получает турниры по источнику
func (r *TournamentPostgres) GetBySource(ctx context.Context, source string) ([]*entities.Tournament, error) {
	query := `SELECT * FROM tournaments WHERE source = $1 ORDER BY name`

	var tournaments []*entities.Tournament
	err := r.db.SelectContext(ctx, &tournaments, query, source)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournaments by source: %w", err)
	}
	return tournaments, nil
}

// GetWithTeams получает турниры где уже есть команды (через player_teams)
func (r *TournamentPostgres) GetWithTeams(ctx context.Context, source string) ([]*entities.Tournament, error) {
	query := `
		SELECT DISTINCT t.*
		FROM tournaments t
		INNER JOIN player_teams pt ON pt.tournament_id = t.id AND pt.source = t.source
		WHERE t.source = $1
		ORDER BY t.name
	`

	var tournaments []*entities.Tournament
	err := r.db.SelectContext(ctx, &tournaments, query, source)
	if err != nil {
		return nil, fmt.Errorf("failed to get tournaments with teams: %w", err)
	}
	return tournaments, nil
}

// List получает список junior турниров (исключая FHSPB)
func (r *TournamentPostgres) List(ctx context.Context) ([]*entities.Tournament, error) {
	query := `SELECT * FROM tournaments WHERE id NOT LIKE 'spb:%' AND url IS NOT NULL AND domain IS NOT NULL ORDER BY name`

	var tournaments []*entities.Tournament
	err := r.db.SelectContext(ctx, &tournaments, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tournaments: %w", err)
	}
	return tournaments, nil
}

