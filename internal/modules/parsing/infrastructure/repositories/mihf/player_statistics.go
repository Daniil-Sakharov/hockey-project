package mihf

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type PlayerStatistics struct {
	ID                int       `db:"id"`
	PlayerID          string    `db:"player_id"`
	TeamID            string    `db:"team_id"`
	TournamentID      string    `db:"tournament_id"`
	Games             int       `db:"games"`
	Goals             int       `db:"goals"`
	Assists           int       `db:"assists"`
	Points            int       `db:"points"`
	PenaltyMinutes    int       `db:"penalty_minutes"`
	GoalsPowerPlay    *int      `db:"goals_power_play"`
	GoalsShortHanded  *int      `db:"goals_short_handed"`
	GoalsEvenStrength *int      `db:"goals_even_strength"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
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
			games, goals, assists, points, penalty_minutes,
			goals_power_play, goals_short_handed, goals_even_strength,
			created_at, updated_at
		)
		SELECT $1, $2, $3, NULL, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		WHERE NOT EXISTS (
			SELECT 1 FROM player_statistics
			WHERE player_id = $1 AND team_id = $2 AND tournament_id = $3 AND group_name IS NULL
		)`

	result, err := r.db.ExecContext(ctx, query,
		s.PlayerID, s.TeamID, s.TournamentID,
		s.Games, s.Goals, s.Assists, s.Points, s.PenaltyMinutes,
		s.GoalsPowerPlay, s.GoalsShortHanded, s.GoalsEvenStrength,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		updateQuery := `
			UPDATE player_statistics SET
				games = $4, goals = $5, assists = $6, points = $7, penalty_minutes = $8,
				goals_power_play = $9, goals_short_handed = $10, goals_even_strength = $11,
				updated_at = NOW()
			WHERE player_id = $1 AND team_id = $2 AND tournament_id = $3 AND group_name IS NULL`

		_, err = r.db.ExecContext(ctx, updateQuery,
			s.PlayerID, s.TeamID, s.TournamentID,
			s.Games, s.Goals, s.Assists, s.Points, s.PenaltyMinutes,
			s.GoalsPowerPlay, s.GoalsShortHanded, s.GoalsEvenStrength,
		)
	}

	return err
}
