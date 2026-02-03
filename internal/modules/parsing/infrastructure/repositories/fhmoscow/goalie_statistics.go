package fhmoscow

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type GoalieStatistics struct {
	ID             int       `db:"id"`
	PlayerID       string    `db:"player_id"`
	TeamID         string    `db:"team_id"`
	TournamentID   string    `db:"tournament_id"`
	Games          int       `db:"games"`
	MinutesPlayed  *int      `db:"minutes"`
	GoalsAgainst   int       `db:"goals_against"`
	SavePercentage *float64  `db:"save_percentage"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type GoalieStatisticsRepository struct {
	db *sqlx.DB
}

func NewGoalieStatisticsRepository(db *sqlx.DB) *GoalieStatisticsRepository {
	return &GoalieStatisticsRepository{db: db}
}

func (r *GoalieStatisticsRepository) Upsert(ctx context.Context, s *GoalieStatistics) error {
	query := `
		INSERT INTO goalie_statistics (
			player_id, team_id, tournament_id,
			games, minutes, goals_against, save_percentage,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (player_id, team_id, tournament_id) DO UPDATE SET
			games = EXCLUDED.games,
			minutes = EXCLUDED.minutes,
			goals_against = EXCLUDED.goals_against,
			save_percentage = EXCLUDED.save_percentage,
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query,
		s.PlayerID, s.TeamID, s.TournamentID,
		s.Games, s.MinutesPlayed, s.GoalsAgainst, s.SavePercentage,
	)
	return err
}
