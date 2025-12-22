package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// insertBatch вставляет статистики в БД
func (r *StatisticsPostgres) insertBatch(ctx context.Context, stats []*entities.PlayerStatistic) (int, error) {
	query := `
		INSERT INTO player_statistics (
			tournament_id, player_id, team_id, group_name, birth_year,
			games, goals, assists, points, plus, minus, plus_minus, penalty_minutes,
			goals_even_strength, goals_power_play, goals_short_handed,
			goals_period_1, goals_period_2, goals_period_3, goals_overtime,
			hat_tricks, game_winning_goals,
			goals_per_game, points_per_game, penalty_minutes_per_game,
			points_avg, penalty_avg,
			created_at, updated_at
		) VALUES (
			:tournament_id, :player_id, :team_id, :group_name, :birth_year,
			:games, :goals, :assists, :points, :plus, :minus, :plus_minus, :penalty_minutes,
			:goals_even_strength, :goals_power_play, :goals_short_handed,
			:goals_period_1, :goals_period_2, :goals_period_3, :goals_overtime,
			:hat_tricks, :game_winning_goals,
			:goals_per_game, :points_per_game, :penalty_minutes_per_game,
			:points_avg, :penalty_avg,
			:created_at, :updated_at
		)
		ON CONFLICT (tournament_id, player_id, group_name, birth_year) DO UPDATE SET
			team_id = EXCLUDED.team_id,
			games = EXCLUDED.games, goals = EXCLUDED.goals, assists = EXCLUDED.assists,
			points = EXCLUDED.points, plus = EXCLUDED.plus, minus = EXCLUDED.minus,
			plus_minus = EXCLUDED.plus_minus, penalty_minutes = EXCLUDED.penalty_minutes,
			goals_even_strength = EXCLUDED.goals_even_strength,
			goals_power_play = EXCLUDED.goals_power_play,
			goals_short_handed = EXCLUDED.goals_short_handed,
			goals_period_1 = EXCLUDED.goals_period_1, goals_period_2 = EXCLUDED.goals_period_2,
			goals_period_3 = EXCLUDED.goals_period_3, goals_overtime = EXCLUDED.goals_overtime,
			hat_tricks = EXCLUDED.hat_tricks, game_winning_goals = EXCLUDED.game_winning_goals,
			goals_per_game = EXCLUDED.goals_per_game, points_per_game = EXCLUDED.points_per_game,
			penalty_minutes_per_game = EXCLUDED.penalty_minutes_per_game,
			points_avg = EXCLUDED.points_avg, penalty_avg = EXCLUDED.penalty_avg,
			updated_at = EXCLUDED.updated_at
	`

	result, err := r.db.NamedExecContext(ctx, query, stats)
	if err != nil {
		return 0, fmt.Errorf("failed to create statistics batch: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

// GetByPlayerID получает статистики игрока
func (r *StatisticsPostgres) GetByPlayerID(ctx context.Context, playerID string) ([]*entities.PlayerStatistic, error) {
	query := `
		SELECT id, tournament_id, player_id, team_id, group_name, birth_year,
			   games, goals, assists, points, plus, minus, plus_minus, penalty_minutes,
			   goals_even_strength, goals_power_play, goals_short_handed,
			   goals_period_1, goals_period_2, goals_period_3, goals_overtime,
			   hat_tricks, game_winning_goals,
			   goals_per_game, points_per_game, penalty_minutes_per_game,
			   points_avg, penalty_avg, created_at, updated_at
		FROM player_statistics WHERE player_id = $1
	`

	var stats []*entities.PlayerStatistic
	err := r.db.SelectContext(ctx, &stats, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by player ID: %w", err)
	}

	return stats, nil
}

// GetByTournament получает статистики турнира
func (r *StatisticsPostgres) GetByTournament(ctx context.Context, tournamentID string) ([]*entities.PlayerStatistic, error) {
	query := `
		SELECT id, tournament_id, player_id, team_id, group_name, birth_year,
			   games, goals, assists, points, plus, minus, plus_minus, penalty_minutes,
			   goals_even_strength, goals_power_play, goals_short_handed,
			   goals_period_1, goals_period_2, goals_period_3, goals_overtime,
			   hat_tricks, game_winning_goals,
			   goals_per_game, points_per_game, penalty_minutes_per_game,
			   points_avg, penalty_avg, created_at, updated_at
		FROM player_statistics WHERE tournament_id = $1
	`

	var stats []*entities.PlayerStatistic
	err := r.db.SelectContext(ctx, &stats, query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by tournament: %w", err)
	}

	return stats, nil
}

// DeleteByTournament удаляет статистики турнира
func (r *StatisticsPostgres) DeleteByTournament(ctx context.Context, tournamentID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM player_statistics WHERE tournament_id = $1`, tournamentID)
	if err != nil {
		return fmt.Errorf("failed to delete statistics by tournament: %w", err)
	}
	return nil
}

// DeleteAll удаляет все статистики
func (r *StatisticsPostgres) DeleteAll(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM player_statistics`)
	if err != nil {
		return fmt.Errorf("failed to delete all statistics: %w", err)
	}
	return nil
}

// CountAll возвращает общее количество статистик
func (r *StatisticsPostgres) CountAll(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM player_statistics`)
	if err != nil {
		return 0, fmt.Errorf("failed to count statistics: %w", err)
	}
	return count, nil
}
