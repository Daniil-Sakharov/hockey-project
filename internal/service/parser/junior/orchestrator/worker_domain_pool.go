package orchestrator

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// DomainTask задача для парсинга домена
type DomainTask struct {
	Domain string
	Index  int // Для логирования (1/23, 2/23...)
	Total  int
}

// DomainResult результат парсинга домена
type DomainResult struct {
	Domain      string
	Tournaments []*tournament.Tournament
	IsDuplicate bool // true если домен полностью дубликат
	Error       error
}

// DomainWorkerPool Worker Pool для параллельного парсинга доменов
type DomainWorkerPool struct {
	workerCount  int
	tasks        chan DomainTask
	results      chan DomainResult
	wg           sync.WaitGroup
	ctx          context.Context
	orchestrator *orchestratorService
	globalDedup  *sync.Map // Глобальная дедупликация турниров по ID
}

// NewDomainWorkerPool создает Worker Pool для доменов
func NewDomainWorkerPool(
	ctx context.Context,
	orchestrator *orchestratorService,
	workerCount int,
	globalDedup *sync.Map,
) *DomainWorkerPool {
	return &DomainWorkerPool{
		workerCount:  workerCount,
		tasks:        make(chan DomainTask, workerCount*2),
		results:      make(chan DomainResult, workerCount*2),
		ctx:          ctx,
		orchestrator: orchestrator,
		globalDedup:  globalDedup,
	}
}

// Start запускает воркеры
func (wp *DomainWorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker функция каждого воркера
func (wp *DomainWorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		// Проверка контекста (отмена)
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		// Обрабатываем домен
		result := wp.processDomain(workerID, task)
		wp.results <- result
	}
}

// processDomain обрабатывает один домен (полный цикл)
func (wp *DomainWorkerPool) processDomain(workerID int, task DomainTask) DomainResult {
	ctx := wp.ctx
	domain := task.Domain

	logger.Info(ctx, "")
	logger.Info(ctx, fmt.Sprintf("🌍 Worker %d: Domain %d/%d: %s",
		workerID, task.Index, task.Total, domain))

	// ШАГ 1: Быстрая проверка на дубликат (парсим только первый сезон)
	// С retry для HTTP 500
	isDuplicate, err := wp.orchestrator.isDuplicateDomain(ctx, domain, wp.globalDedup)
	if err != nil {
		// Проверяем HTTP 500 и делаем retry
		if strings.Contains(err.Error(), "500") || strings.Contains(err.Error(), "HTTP статус: 500") {
			logger.Warn(ctx, "  ⚠️  HTTP 500 detected, retry after 10 seconds...")
			time.Sleep(10 * time.Second)

			// Повторная попытка
			isDuplicate, err = wp.orchestrator.isDuplicateDomain(ctx, domain, wp.globalDedup)
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("  ⚠️  Quick check failed after retry: %v", err))
				return DomainResult{
					Domain: domain,
					Error:  fmt.Errorf("quick check failed after retry: %w", err),
				}
			}
		} else {
			logger.Warn(ctx, fmt.Sprintf("  ⚠️  Quick check failed: %v", err))
			return DomainResult{
				Domain: domain,
				Error:  fmt.Errorf("quick check failed: %w", err),
			}
		}
	}

	if isDuplicate {
		// Домен - полная копия, пропускаем
		logger.Info(ctx, "  🔁 DUPLICATE domain detected - SKIPPING")
		return DomainResult{
			Domain:      domain,
			IsDuplicate: true,
		}
	}

	// ШАГ 2: Парсим ВСЕ сезоны (Worker Pool внутри)
	// С retry для HTTP 500
	logger.Info(ctx, "  🚀 Parsing ALL seasons...")
	tournamentsDTO, err := wp.orchestrator.juniorService.ParseAllSeasonsTournaments(ctx, domain)
	if err != nil {
		// Проверяем HTTP 500 и делаем retry
		if strings.Contains(err.Error(), "500") || strings.Contains(err.Error(), "HTTP статус: 500") {
			logger.Warn(ctx, "  ⚠️  HTTP 500 detected, retry after 10 seconds...")
			time.Sleep(10 * time.Second)

			// Повторная попытка
			tournamentsDTO, err = wp.orchestrator.juniorService.ParseAllSeasonsTournaments(ctx, domain)
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("  ⚠️  Failed to parse all seasons after retry: %v", err))
				return DomainResult{
					Domain: domain,
					Error:  fmt.Errorf("failed to parse all seasons after retry: %w", err),
				}
			}
		} else {
			logger.Warn(ctx, fmt.Sprintf("  ⚠️  Failed to parse all seasons: %v", err))
			return DomainResult{
				Domain: domain,
				Error:  fmt.Errorf("failed to parse all seasons: %w", err),
			}
		}
	}

	// ШАГ 3: Фильтруем дубликаты (на всякий случай, хотя первый сезон уже зарегистрирован)
	var uniqueTournaments []junior.TournamentDTO
	duplicateCount := 0

	for _, t := range tournamentsDTO {
		_, exists := wp.globalDedup.LoadOrStore(t.ID, true)
		if !exists {
			uniqueTournaments = append(uniqueTournaments, t)
		} else {
			duplicateCount++
		}
	}

	if duplicateCount > 0 {
		logger.Info(ctx, fmt.Sprintf("  📊 Filtered %d duplicates from other seasons", duplicateCount))
	}

	logger.Info(ctx, fmt.Sprintf("  ✅ Found %d unique tournaments (all seasons)", len(uniqueTournaments)))

	// ШАГ 4: Сохраняем турниры батчами
	logger.Info(ctx, "  💾 Saving tournaments to database...")
	saved, err := wp.orchestrator.SaveTournamentsBatch(ctx, uniqueTournaments, 100)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("  ❌ CRITICAL: Failed to save tournaments: %v", err))
		logger.Warn(ctx, "  ⏭️  SKIPPING domain processing due to tournament save error")
		return DomainResult{
			Domain: domain,
			Error:  fmt.Errorf("failed to save tournaments: %w", err),
		}
	}

	logger.Info(ctx, fmt.Sprintf("  ✅ Saved %d tournaments to database", len(saved)))
	logger.Info(ctx, "")
	logger.Info(ctx, fmt.Sprintf("🎉 Domain %s COMPLETED SUCCESSFULLY! (%d tournaments ready for processing)",
		domain, len(saved)))
	logger.Info(ctx, "================================================================================")

	return DomainResult{
		Domain:      domain,
		Tournaments: saved,
		IsDuplicate: false,
	}
}

// AddTask добавляет задачу в очередь
func (wp *DomainWorkerPool) AddTask(task DomainTask) {
	wp.tasks <- task
}

// Close закрывает очередь задач (больше задач не будет)
func (wp *DomainWorkerPool) Close() {
	close(wp.tasks)
}

// Wait ждет завершения всех воркеров и закрывает канал результатов
func (wp *DomainWorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results возвращает канал результатов
func (wp *DomainWorkerPool) Results() <-chan DomainResult {
	return wp.results
}
