package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// StatsOverview represents aggregated statistics.
type StatsOverview struct {
	Players     int64
	Teams       int64
	Tournaments int64
}

// StatsService provides statistics aggregation.
type StatsService struct {
	db *sqlx.DB
}

// NewStatsService creates a new stats service.
func NewStatsService(db *sqlx.DB) *StatsService {
	return &StatsService{db: db}
}

// GetOverview returns aggregated stats overview.
func (s *StatsService) GetOverview(ctx context.Context) (*StatsOverview, error) {
	var result StatsOverview

	// Count players
	err := s.db.GetContext(ctx, &result.Players, "SELECT COUNT(*) FROM players")
	if err != nil {
		return nil, fmt.Errorf("failed to count players: %w", err)
	}

	// Count teams
	err = s.db.GetContext(ctx, &result.Teams, "SELECT COUNT(*) FROM teams")
	if err != nil {
		return nil, fmt.Errorf("failed to count teams: %w", err)
	}

	// Count tournaments
	err = s.db.GetContext(ctx, &result.Tournaments, "SELECT COUNT(*) FROM tournaments")
	if err != nil {
		return nil, fmt.Errorf("failed to count tournaments: %w", err)
	}

	return &result, nil
}
