package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// processDomain обрабатывает один домен (полный цикл)
func (wp *DomainWorkerPool) processDomain(workerID int, task DomainTask) DomainResult {
	ctx := wp.ctx
	domain := task.Domain

	logger.Info(ctx, "")
	logger.Info(ctx, fmt.Sprintf("🌍 Worker %d: Domain %d/%d: %s", workerID, task.Index, task.Total, domain))

	// ШАГ 1: Быстрая проверка на дубликат
	isDuplicate, err := wp.checkDuplicate(domain)
	if err != nil {
		return DomainResult{Domain: domain, Error: err}
	}

	if isDuplicate {
		logger.Info(ctx, "  🔁 DUPLICATE domain detected - SKIPPING")
		return DomainResult{Domain: domain, IsDuplicate: true}
	}

	// ШАГ 2: Парсим ВСЕ сезоны
	tournamentsDTO, err := wp.parseAllSeasons(domain)
	if err != nil {
		return DomainResult{Domain: domain, Error: err}
	}

	// ШАГ 3: Фильтруем дубликаты
	uniqueTournaments := wp.filterDuplicates(tournamentsDTO)

	logger.Info(ctx, fmt.Sprintf("  ✅ Found %d unique tournaments (all seasons)", len(uniqueTournaments)))

	// ШАГ 4: Сохраняем турниры
	saved, err := wp.saveTournaments(domain, uniqueTournaments)
	if err != nil {
		return DomainResult{Domain: domain, Error: err}
	}

	logger.Info(ctx, fmt.Sprintf("🎉 Domain %s COMPLETED! (%d tournaments)", domain, len(saved)))
	logger.Info(ctx, "================================================================================")

	return DomainResult{Domain: domain, Tournaments: saved, IsDuplicate: false}
}

// checkDuplicate проверяет домен на дубликат с retry
func (wp *DomainWorkerPool) checkDuplicate(domain string) (bool, error) {
	ctx := wp.ctx

	isDuplicate, err := wp.orchestrator.isDuplicateDomain(ctx, domain, wp.globalDedup)
	if err != nil {
		if strings.Contains(err.Error(), "500") {
			logger.Warn(ctx, "  ⚠️  HTTP 500, retry after 10s...")
			time.Sleep(10 * time.Second)

			isDuplicate, err = wp.orchestrator.isDuplicateDomain(ctx, domain, wp.globalDedup)
			if err != nil {
				return false, fmt.Errorf("quick check failed after retry: %w", err)
			}
		} else {
			return false, fmt.Errorf("quick check failed: %w", err)
		}
	}

	return isDuplicate, nil
}

// parseAllSeasons парсит все сезоны с retry
func (wp *DomainWorkerPool) parseAllSeasons(domain string) ([]junior.TournamentDTO, error) {
	ctx := wp.ctx

	logger.Info(ctx, "  🚀 Parsing ALL seasons...")
	tournamentsDTO, err := wp.orchestrator.juniorService.ParseAllSeasonsTournaments(ctx, domain)
	if err != nil {
		if strings.Contains(err.Error(), "500") {
			logger.Warn(ctx, "  ⚠️  HTTP 500, retry after 10s...")
			time.Sleep(10 * time.Second)

			tournamentsDTO, err = wp.orchestrator.juniorService.ParseAllSeasonsTournaments(ctx, domain)
			if err != nil {
				return nil, fmt.Errorf("failed to parse all seasons after retry: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to parse all seasons: %w", err)
		}
	}

	return tournamentsDTO, nil
}

// filterDuplicates фильтрует дубликаты турниров
func (wp *DomainWorkerPool) filterDuplicates(tournaments []junior.TournamentDTO) []junior.TournamentDTO {
	var unique []junior.TournamentDTO
	duplicateCount := 0

	for _, t := range tournaments {
		_, exists := wp.globalDedup.LoadOrStore(t.ID, true)
		if !exists {
			unique = append(unique, t)
		} else {
			duplicateCount++
		}
	}

	if duplicateCount > 0 {
		logger.Info(wp.ctx, fmt.Sprintf("  📊 Filtered %d duplicates", duplicateCount))
	}

	return unique
}

// saveTournaments сохраняет турниры в БД
func (wp *DomainWorkerPool) saveTournaments(domain string, tournaments []junior.TournamentDTO) ([]*entities.Tournament, error) {
	ctx := wp.ctx

	logger.Info(ctx, "  💾 Saving tournaments to database...")
	saved, err := wp.orchestrator.SaveTournamentsBatch(ctx, tournaments, 100)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("  ❌ Failed to save tournaments: %v", err))
		return nil, fmt.Errorf("failed to save tournaments: %w", err)
	}

	logger.Info(ctx, fmt.Sprintf("  ✅ Saved %d tournaments", len(saved)))
	return saved, nil
}
