package parser

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

const taskTimeout = 1 * time.Minute // Таймаут на обработку одной команды (уменьшен для тестирования)

// TeamTask задача для парсинга команды
type TeamTask struct {
	Team       *entities.Team
	Tournament *entities.Tournament
	BirthYear  *int
	GroupName  *string
	Index      int
	Total      int
}

// TeamResult результат парсинга команды
type TeamResult struct {
	TeamName     string
	PlayersCount int
	Error        error
}

// TeamWorkerPool Worker Pool для параллельного парсинга команд
type TeamWorkerPool struct {
	workerCount  int
	tasks        chan TeamTask
	results      chan TeamResult
	wg           sync.WaitGroup
	ctx          context.Context
	orchestrator *orchestratorService
}

// NewTeamWorkerPool создает Worker Pool для команд
func NewTeamWorkerPool(
	ctx context.Context,
	orchestrator *orchestratorService,
	workerCount int,
) *TeamWorkerPool {
	return &TeamWorkerPool{
		workerCount:  workerCount,
		tasks:        make(chan TeamTask, workerCount*2),
		results:      make(chan TeamResult, workerCount*2),
		ctx:          ctx,
		orchestrator: orchestrator,
	}
}

// Start запускает воркеры
func (wp *TeamWorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker функция каждого воркера
func (wp *TeamWorkerPool) worker(workerID int) {
	defer func() {
		logger.Debug(wp.ctx, fmt.Sprintf("    👷 Worker %d: exiting", workerID))
		wp.wg.Done()
	}()

	for task := range wp.tasks {
		// Проверка контекста (отмена)
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		// Обрабатываем команду с таймаутом
		result := wp.processTeamWithTimeout(workerID, task)
		logger.Debug(wp.ctx, fmt.Sprintf("    👷 Worker %d: sending result for %s", workerID, task.Team.Name))
		wp.results <- result
		logger.Debug(wp.ctx, fmt.Sprintf("    👷 Worker %d: result sent for %s", workerID, task.Team.Name))
	}
}

// processTeamWithTimeout обрабатывает команду с таймаутом
func (wp *TeamWorkerPool) processTeamWithTimeout(workerID int, task TeamTask) TeamResult {
	resultCh := make(chan TeamResult, 1)

	go func() {
		resultCh <- wp.processTeam(workerID, task)
	}()

	select {
	case result := <-resultCh:
		return result
	case <-time.After(taskTimeout):
		logger.Warn(wp.ctx, fmt.Sprintf("    ⏱️  Worker %d: TIMEOUT processing team %s", workerID, task.Team.Name))
		return TeamResult{
			TeamName: task.Team.Name,
			Error:    fmt.Errorf("timeout after %v", taskTimeout),
		}
	case <-wp.ctx.Done():
		return TeamResult{
			TeamName: task.Team.Name,
			Error:    wp.ctx.Err(),
		}
	}
}

// processTeam обрабатывает одну команду (парсинг игроков + сохранение)
func (wp *TeamWorkerPool) processTeam(workerID int, task TeamTask) TeamResult {
	ctx := wp.ctx
	t := task.Team
	tournament := task.Tournament

	logger.Info(ctx, fmt.Sprintf("    🏒 Worker %d: Team %d/%d: %s",
		workerID, task.Index, task.Total, t.Name))

	// Парсим игроков (передаем домен турнира и контекст year/group)
	err := wp.orchestrator.SavePlayers(ctx, tournament.Domain, t.URL, t.ID, tournament.ID, tournament, task.BirthYear, task.GroupName)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("      ⚠️  Worker %d failed: %v", workerID, err))
		return TeamResult{
			TeamName: t.Name,
			Error:    err,
		}
	}

	logger.Info(ctx, fmt.Sprintf("    ✅ Worker %d: Team %s DONE", workerID, t.Name))
	return TeamResult{
		TeamName: t.Name,
	}
}

// AddTask добавляет задачу в очередь
func (wp *TeamWorkerPool) AddTask(task TeamTask) {
	wp.tasks <- task
}

// Close закрывает очередь задач (больше задач не будет)
func (wp *TeamWorkerPool) Close() {
	close(wp.tasks)
}

// Wait ждет завершения всех воркеров и закрывает канал результатов
func (wp *TeamWorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results возвращает канал результатов
func (wp *TeamWorkerPool) Results() <-chan TeamResult {
	return wp.results
}
