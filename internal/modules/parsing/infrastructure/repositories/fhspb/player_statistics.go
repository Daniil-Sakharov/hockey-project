package fhspb

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlayerStatistics struct {
	ID             int       `db:"id"`
	PlayerID       string    `db:"player_id"`
	TeamID         string    `db:"team_id"`
	TournamentID   string    `db:"tournament_id"`
	Games          int       `db:"games"`
	Points         int       `db:"points"`
	PointsAvg      *float64  `db:"points_avg"`
	Goals          int       `db:"goals"`
	Assists        int       `db:"assists"`
	PlusMinus      int       `db:"plus_minus"`
	PenaltyMinutes int       `db:"penalty_minutes"`
	PenaltyAvg     *float64  `db:"penalty_avg"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type PlayerStatisticsRepository struct {
	db *sqlx.DB
}

func NewPlayerStatisticsRepository(db *sqlx.DB) *PlayerStatisticsRepository {
	return &PlayerStatisticsRepository{db: db}
}

func (r *PlayerStatisticsRepository) Upsert(ctx context.Context, s *PlayerStatistics) error {
	query := `
		INSERT INTO player_statistics (
			player_id, team_id, tournament_id, group_name,
			games, points, points_avg, goals, assists, plus_minus, penalty_minutes, penalty_avg,
			created_at, updated_at
		) 
		SELECT $1, $2, $3, NULL, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		WHERE NOT EXISTS (
			SELECT 1 FROM player_statistics 
			WHERE player_id = $1 AND team_id = $2 AND tournament_id = $3 AND group_name IS NULL
		)`

	result, err := r.db.ExecContext(ctx, query,
		s.PlayerID, s.TeamID, s.TournamentID,
		s.Games, s.Points, s.PointsAvg, s.Goals, s.Assists, s.PlusMinus, s.PenaltyMinutes, s.PenaltyAvg,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		updateQuery := `
			UPDATE player_statistics SET
				games = $4, points = $5, points_avg = $6, goals = $7, assists = $8,
				plus_minus = $9, penalty_minutes = $10, penalty_avg = $11, updated_at = NOW()
			WHERE player_id = $1 AND team_id = $2 AND tournament_id = $3 AND group_name IS NULL`

		_, err = r.db.ExecContext(ctx, updateQuery,
			s.PlayerID, s.TeamID, s.TournamentID,
			s.Games, s.Points, s.PointsAvg, s.Goals, s.Assists, s.PlusMinus, s.PenaltyMinutes, s.PenaltyAvg,
		)
	}

	return err
}
