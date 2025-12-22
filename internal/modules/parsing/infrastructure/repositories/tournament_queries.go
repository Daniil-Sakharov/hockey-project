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
