package orchestrator

import (
	"context"
	"fmt"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// RunJuniorParsing парсит junior.fhr.ru (ВСЕ домены → турниры → команды → игроки)
func (s *orchestratorService) RunJuniorParsing(ctx context.Context) error {
	// ЭТАП 1: Находим ВСЕ домены *.fhr.ru
	logger.Info(ctx, "📊 STAGE 1: Discovering all domains...")
	domains, err := s.juniorService.ParseDomains(ctx)
	if err != nil {
		return fmt.Errorf("failed to discover domains: %w", err)
	}
	logger.Info(ctx, fmt.Sprintf("  ✅ Found %d domains", len(domains)))

	// ЭТАП 2: Глобальная дедупликация турниров по ID
	var globalTournamentDedup sync.Map

	// ЭТАП 3: Worker Pool для доменов
	pool := NewDomainWorkerPool(ctx, s, s.config.DomainWorkers(), &globalTournamentDedup)
	pool.Start()

	// ЭТАП 4: Добавляем домены в очередь
	go func() {
		for idx, domain := range domains {
			pool.AddTask(DomainTask{
				Domain: domain,
				Index:  idx + 1,
				Total:  len(domains),
			})
		}
		pool.Close()
	}()

	// ЭТАП 5: Ждем завершения (в отдельной goroutine)
	go pool.Wait()

	// ЭТАП 6: Собираем результаты
	var allTournaments []*tournament.Tournament
	stats := DomainParsingStats{
		skippedDomains: make([]string, 0),
	}

	for result := range pool.Results() {
		stats.total++

		if result.Error != nil {
			logger.Warn(ctx, fmt.Sprintf("  ⚠️  Domain error: %s - %v", result.Domain, result.Error))
			stats.errors++
			continue
		}

		if result.IsDuplicate {
			stats.skippedDomains = append(stats.skippedDomains, result.Domain)
			continue
		}

		stats.parsedDomains++
		allTournaments = append(allTournaments, result.Tournaments...)
	}

	// ЭТАП 7: Финальная статистика доменов
	s.logDomainStats(ctx, stats, len(domains), len(allTournaments))

	// ЭТАП 8: Продолжаем с турнирами (ParseTeams → ParsePlayers)
	logger.Info(ctx, "")
	logger.Info(ctx, "================================================================================")
	logger.Info(ctx, "🚀 STARTING TOURNAMENT & PLAYER PROCESSING...")
	logger.Info(ctx, fmt.Sprintf("  Total tournaments to process: %d", len(allTournaments)))
	logger.Info(ctx, "================================================================================")
	
	return s.processTournaments(ctx, allTournaments)
}

// DomainParsingStats статистика парсинга доменов
type DomainParsingStats struct {
	total          int
	parsedDomains  int
	skippedDomains []string
	errors         int
}

// logDomainStats выводит статистику парсинга доменов
func (s *orchestratorService) logDomainStats(
	ctx context.Context,
	stats DomainParsingStats,
	totalDomains int,
	tournamentsCount int,
) {
	logger.Info(ctx, "")
	logger.Info(ctx, "================================================================================")
	logger.Info(ctx, "📊 DOMAIN PARSING STATISTICS:")
	logger.Info(ctx, fmt.Sprintf("  Total domains: %d", totalDomains))
	logger.Info(ctx, fmt.Sprintf("  Parsed domains: %d", stats.parsedDomains))
	logger.Info(ctx, fmt.Sprintf("  Skipped (duplicates): %d", len(stats.skippedDomains)))
	logger.Info(ctx, fmt.Sprintf("  Errors: %d", stats.errors))
	logger.Info(ctx, fmt.Sprintf("  Total unique tournaments: %d", tournamentsCount))
	logger.Info(ctx, "================================================================================")
}
