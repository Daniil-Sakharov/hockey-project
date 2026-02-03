package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TopScorer represents a top scorer player.
type TopScorer struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	Team    string `db:"team_name"`
	Goals   int    `db:"goals"`
	Assists int    `db:"assists"`
	Games   int    `db:"games"`
}

// RankingService provides player rankings.
type RankingService struct {
	db *sqlx.DB
}

// NewRankingService creates a new ranking service.
func NewRankingService(db *sqlx.DB) *RankingService {
	return &RankingService{db: db}
}

// GetTopScorers returns top scorers by goals.
func (s *RankingService) GetTopScorers(ctx context.Context, limit int) ([]TopScorer, error) {
	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT
			p.id,
			p.name,
			COALESCE(t.name, '') as team_name,
			COALESCE(SUM(ps.goals), 0)::int as goals,
			COALESCE(SUM(ps.assists), 0)::int as assists,
			COALESCE(SUM(ps.games), 0)::int as games
		FROM players p
		LEFT JOIN player_statistics ps ON p.id = ps.player_id
		LEFT JOIN player_teams pt ON p.id = pt.player_id
		LEFT JOIN teams t ON pt.team_id = t.id
		GROUP BY p.id, p.name, t.name
		HAVING SUM(ps.goals) > 0
		ORDER BY goals DESC, assists DESC
		LIMIT $1
	`

	var scorers []TopScorer
	err := s.db.SelectContext(ctx, &scorers, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top scorers: %w", err)
	}

	return scorers, nil
}
