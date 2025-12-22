package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type GoalieStatisticsPostgres struct {
	db *sqlx.DB
}

func NewGoalieStatisticsPostgres(db *sqlx.DB) *GoalieStatisticsPostgres {
	return &GoalieStatisticsPostgres{db: db}
}

func (r *GoalieStatisticsPostgres) Upsert(ctx context.Context, s *entities.GoalieStatistic) error {
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

func (r *GoalieStatisticsPostgres) GetByTournament(ctx context.Context, tournamentID string) ([]*entities.GoalieStatistic, error) {
	query := `
		SELECT id, player_id, team_id, tournament_id, games, minutes, goals_against, shots_against, 
		       save_percentage, goals_against_avg, wins, shutouts, assists, penalty_minutes, created_at, updated_at 
		FROM goalie_statistics WHERE tournament_id = $1 ORDER BY save_percentage DESC`

	var stats []*entities.GoalieStatistic
	err := r.db.SelectContext(ctx, &stats, query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goalie statistics: %w", err)
	}
	return stats, nil
}
