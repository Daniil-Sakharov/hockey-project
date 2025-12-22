package fhspb

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type GoalieStatistics struct {
	ID              int       `db:"id"`
	PlayerID        string    `db:"player_id"`
	TeamID          string    `db:"team_id"`
	TournamentID    string    `db:"tournament_id"`
	Games           int       `db:"games"`
	Minutes         int       `db:"minutes"`
	GoalsAgainst    int       `db:"goals_against"`
	ShotsAgainst    int       `db:"shots_against"`
	SavePercentage  *float64  `db:"save_percentage"`
	GoalsAgainstAvg *float64  `db:"goals_against_avg"`
	Wins            int       `db:"wins"`
	Shutouts        int       `db:"shutouts"`
	Assists         int       `db:"assists"`
	PenaltyMinutes  int       `db:"penalty_minutes"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
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
			games, minutes, goals_against, shots_against, save_percentage, goals_against_avg,
			wins, shutouts, assists, penalty_minutes,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
		ON CONFLICT (player_id, team_id, tournament_id) DO UPDATE SET
			games = EXCLUDED.games,
			minutes = EXCLUDED.minutes,
			goals_against = EXCLUDED.goals_against,
			shots_against = EXCLUDED.shots_against,
			save_percentage = EXCLUDED.save_percentage,
			goals_against_avg = EXCLUDED.goals_against_avg,
			wins = EXCLUDED.wins,
			shutouts = EXCLUDED.shutouts,
			assists = EXCLUDED.assists,
			penalty_minutes = EXCLUDED.penalty_minutes,
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query,
		s.PlayerID, s.TeamID, s.TournamentID,
		s.Games, s.Minutes, s.GoalsAgainst, s.ShotsAgainst, s.SavePercentage, s.GoalsAgainstAvg,
		s.Wins, s.Shutouts, s.Assists, s.PenaltyMinutes,
	)
	return err
}
