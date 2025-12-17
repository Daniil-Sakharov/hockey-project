package player_statistics

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// CreateBatch создает несколько записей статистики за одну транзакцию
// Возвращает количество реально вставленных/обновленных записей
func (r *repository) CreateBatch(ctx context.Context, stats []*player_statistics.PlayerStatistic) (int, error) {
	if len(stats) == 0 {
		return 0, nil
	}

	// КРИТИЧНО: Валидируем все записи ДО начала транзакции
	// Иначе FK error внутри транзакции переведет её в "aborted" состояние
	validStats, _ := r.filterValidStats(ctx, stats)

	if len(validStats) == 0 {
		return 0, nil
	}

	// Логирование пропущенных записей происходит в service layer через statsLogger

	query := `
		INSERT INTO player_statistics (
			tournament_id, player_id, team_id, group_name, birth_year,
			games, goals, assists, points, plus, minus, plus_minus, penalty_minutes,
			goals_even_strength, goals_power_play, goals_short_handed,
			goals_period_1, goals_period_2, goals_period_3, goals_overtime,
			hat_tricks, game_winning_goals,
			goals_per_game, points_per_game, penalty_minutes_per_game
		) VALUES (
			:tournament_id, :player_id, :team_id, :group_name, :birth_year,
			:games, :goals, :assists, :points, :plus, :minus, :plus_minus, :penalty_minutes,
			:goals_even_strength, :goals_power_play, :goals_short_handed,
			:goals_period_1, :goals_period_2, :goals_period_3, :goals_overtime,
			:hat_tricks, :game_winning_goals,
			:goals_per_game, :points_per_game, :penalty_minutes_per_game
		)
		ON CONFLICT (tournament_id, player_id, group_name, birth_year)
		DO UPDATE SET
			team_id = EXCLUDED.team_id,
			games = EXCLUDED.games,
			goals = EXCLUDED.goals,
			assists = EXCLUDED.assists,
			points = EXCLUDED.points,
			plus = EXCLUDED.plus,
			minus = EXCLUDED.minus,
			plus_minus = EXCLUDED.plus_minus,
			penalty_minutes = EXCLUDED.penalty_minutes,
			goals_even_strength = EXCLUDED.goals_even_strength,
			goals_power_play = EXCLUDED.goals_power_play,
			goals_short_handed = EXCLUDED.goals_short_handed,
			goals_period_1 = EXCLUDED.goals_period_1,
			goals_period_2 = EXCLUDED.goals_period_2,
			goals_period_3 = EXCLUDED.goals_period_3,
			goals_overtime = EXCLUDED.goals_overtime,
			hat_tricks = EXCLUDED.hat_tricks,
			game_winning_goals = EXCLUDED.game_winning_goals,
			goals_per_game = EXCLUDED.goals_per_game,
			points_per_game = EXCLUDED.points_per_game,
			penalty_minutes_per_game = EXCLUDED.penalty_minutes_per_game,
			updated_at = NOW()
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Считаем реально вставленные/обновленные строки
	totalAffected := int64(0)

	// Все записи уже валидированы, можно вставлять без проверок
	for _, stat := range validStats {
		result, err := tx.NamedExecContext(ctx, query, stat)
		if err != nil {
			return 0, fmt.Errorf("failed to insert statistic: %w", err)
		}

		// Подсчитываем количество затронутых строк (вставлено или обновлено)
		if affected, err := result.RowsAffected(); err == nil {
			totalAffected += affected
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(totalAffected), nil
}

// filterValidStats фильтрует только те записи, для которых существуют FK в БД
func (r *repository) filterValidStats(ctx context.Context, stats []*player_statistics.PlayerStatistic) ([]*player_statistics.PlayerStatistic, int) {
	if len(stats) == 0 {
		return stats, 0
	}

	// Собираем уникальные player_id и team_id
	playerIDs := make(map[string]bool)
	teamIDs := make(map[string]bool)
	for _, stat := range stats {
		playerIDs[stat.PlayerID] = false // false = не проверено
		teamIDs[stat.TeamID] = false
	}

	// Проверяем существование player_id
	playerIDsList := make([]string, 0, len(playerIDs))
	for id := range playerIDs {
		playerIDsList = append(playerIDsList, id)
	}

	query := `SELECT id FROM players WHERE id = ANY($1)`
	var existingPlayers []string
	if err := r.db.SelectContext(ctx, &existingPlayers, query, playerIDsList); err == nil {
		for _, id := range existingPlayers {
			playerIDs[id] = true // помечаем как существующий
		}
	}

	// Проверяем существование team_id
	teamIDsList := make([]string, 0, len(teamIDs))
	for id := range teamIDs {
		teamIDsList = append(teamIDsList, id)
	}

	query = `SELECT id FROM teams WHERE id = ANY($1)`
	var existingTeams []string
	if err := r.db.SelectContext(ctx, &existingTeams, query, teamIDsList); err == nil {
		for _, id := range existingTeams {
			teamIDs[id] = true
		}
	}

	// Фильтруем только валидные записи
	validStats := make([]*player_statistics.PlayerStatistic, 0, len(stats))
	skippedCount := 0

	for _, stat := range stats {
		playerExists := playerIDs[stat.PlayerID]
		teamExists := teamIDs[stat.TeamID]

		if !playerExists {
			// Логирование происходит в service layer через statsLogger.LogFKConstraintSkip
			skippedCount++
			continue
		}
		if !teamExists {
			// Логирование происходит в service layer через statsLogger.LogFKConstraintSkip
			skippedCount++
			continue
		}

		validStats = append(validStats, stat)
	}

	return validStats, skippedCount
}
