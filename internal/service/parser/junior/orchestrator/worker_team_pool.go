package orchestrator

import (
	"context"
	"fmt"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// TeamTask задача для парсинга команды
type TeamTask struct {
	Team       *team.Team
	Tournament *tournament.Tournament
	Index      int // Для логирования (1/27, 2/27...)
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
	defer wp.wg.Done()

	for task := range wp.tasks {
		// Проверка контекста (отмена)
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		// Обрабатываем команду
		result := wp.processTeam(workerID, task)
		wp.results <- result
	}
}

// processTeam обрабатывает одну команду (парсинг игроков + сохранение)
func (wp *TeamWorkerPool) processTeam(workerID int, task TeamTask) TeamResult {
	ctx := wp.ctx
	team := task.Team
	tournament := task.Tournament

	logger.Info(ctx, fmt.Sprintf("    🏒 Worker %d: Team %d/%d: %s",
		workerID, task.Index, task.Total, team.Name))

	// Парсим игроков
	err := wp.orchestrator.SavePlayers(ctx, team.URL, team.ID, tournament.ID, tournament)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("      ⚠️  Worker %d failed: %v", workerID, err))
		return TeamResult{
			TeamName: team.Name,
			Error:    err,
		}
	}

	return TeamResult{
		TeamName: team.Name,
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
