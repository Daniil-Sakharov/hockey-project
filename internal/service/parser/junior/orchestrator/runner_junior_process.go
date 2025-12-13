package orchestrator

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

const (
	teamWorkers = 10 // Worker Pool для команд (10 параллельно)
)

// processTournaments обрабатывает турниры (команды → игроки) с Worker Pool
func (s *orchestratorService) processTournaments(
	ctx context.Context,
	tournaments []*tournament.Tournament,
) error {
	logger.Info(ctx, "")
	logger.Info(ctx, "📊 STAGE 2: Processing tournaments...")

	totalTeams := 0
	totalErrors := 0

	// Для КАЖДОГО турнира (пока последовательно)
	for idx, t := range tournaments {
		logger.Info(ctx, fmt.Sprintf("  🏆 Tournament %d/%d: %s (ID: %s, URL: %s)", 
			idx+1, len(tournaments), t.Name, t.ID, t.URL))

		// Парсим команды
		logger.Info(ctx, "    🔍 Parsing teams...")
		teamsDTO, err := s.juniorService.ParseTeams(ctx, t.Domain, t.URL)
		if err != nil {
			logger.Warn(ctx, fmt.Sprintf("    ⚠️  Failed to parse teams: %v", err))
			logger.Warn(ctx, "    ⏭️  SKIPPING tournament, continuing with next...")
			totalErrors++
			continue
		}

		logger.Info(ctx, fmt.Sprintf("    ✅ Found %d teams from page", len(teamsDTO)))

		// Сохраняем команды
		logger.Info(ctx, "    💾 Saving teams to database...")
		teams, err := s.SaveTeams(ctx, teamsDTO)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("    ❌ CRITICAL: Failed to save teams: %v", err))
			logger.Warn(ctx, "    ⏭️  SKIPPING tournament, continuing with next...")
			totalErrors++
			continue
		}

		logger.Info(ctx, fmt.Sprintf("    ✅ Saved %d teams to database", len(teams)))

		if len(teams) == 0 {
			logger.Info(ctx, "    ℹ️  No teams to process, moving to next tournament")
			continue
		}

		// НОВОЕ: Worker Pool для параллельного парсинга команд (10 воркеров)
		logger.Info(ctx, fmt.Sprintf("    🚀 Starting team worker pool (%d workers) to parse players...", teamWorkers))
		pool := NewTeamWorkerPool(ctx, s, teamWorkers)
		pool.Start()

		// Добавляем команды в очередь
		go func() {
			for teamIdx, team := range teams {
				pool.AddTask(TeamTask{
					Team:       team,
					Tournament: t,
					Index:      teamIdx + 1,
					Total:      len(teams),
				})
			}
			pool.Close()
		}()

		// Ждем завершения
		go pool.Wait()

		// Собираем результаты
		teamProcessed := 0
		teamErrors := 0

		for result := range pool.Results() {
			if result.Error != nil {
				logger.Warn(ctx, fmt.Sprintf("      ⚠️  Team error: %s - %v", result.TeamName, result.Error))
				teamErrors++
			} else {
				teamProcessed++
			}
		}

		totalTeams += teamProcessed
		totalErrors += teamErrors

		logger.Info(ctx, fmt.Sprintf("    📊 Tournament result: %d teams processed, %d errors", teamProcessed, teamErrors))
		logger.Info(ctx, fmt.Sprintf("    ✅ Tournament COMPLETED"))
	}

	logger.Info(ctx, "")
	logger.Info(ctx, "================================================================================")
	logger.Info(ctx, "📊 FINAL STATISTICS:")
	logger.Info(ctx, fmt.Sprintf("  Tournaments processed: %d", len(tournaments)))
	logger.Info(ctx, fmt.Sprintf("  Teams processed: %d", totalTeams))
	logger.Info(ctx, fmt.Sprintf("  Errors: %d", totalErrors))
	logger.Info(ctx, "================================================================================")

	return nil
}
