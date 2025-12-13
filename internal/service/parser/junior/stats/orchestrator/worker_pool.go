package stats_orchestrator

import (
	"context"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// TournamentTask задача для парсинга статистики турнира
type TournamentTask struct {
	Tournament *tournament.Tournament
	Index      int // Для логирования (1/140, 2/140...)
	Total      int
}

// TournamentResult результат парсинга статистики турнира
type TournamentResult struct {
	Tournament *tournament.Tournament
	Count      int // Количество обработанных записей статистики
	Error      error
}

// TournamentWorkerPool Worker Pool для параллельного парсинга статистики турниров
type TournamentWorkerPool struct {
	workerCount int
	tasks       chan TournamentTask
	results     chan TournamentResult
	wg          sync.WaitGroup
	ctx         context.Context
	service     *service
}

// NewTournamentWorkerPool создает Worker Pool для турниров
func NewTournamentWorkerPool(
	ctx context.Context,
	service *service,
	workerCount int,
) *TournamentWorkerPool {
	return &TournamentWorkerPool{
		workerCount: workerCount,
		tasks:       make(chan TournamentTask, workerCount*2),
		results:     make(chan TournamentResult, workerCount*2),
		ctx:         ctx,
		service:     service,
	}
}

// Start запускает воркеры
func (wp *TournamentWorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker функция каждого воркера
func (wp *TournamentWorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		// Проверка контекста (отмена)
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		// Логируем начало обработки турнира
		wp.service.logger.Printf("\n[Worker %d] [%d/%d] Парсинг турнира: %s (ID: %s)",
			workerID,
			task.Index,
			task.Total,
			task.Tournament.Name,
			task.Tournament.ID)
		wp.service.logger.Printf("[Worker %d]         URL: %s%s",
			workerID,
			task.Tournament.Domain,
			task.Tournament.URL)

		// Парсим статистику турнира
		count, err := wp.service.parseTournamentStats(
			wp.ctx,
			task.Tournament,
		)

		// Логируем результат
		if err != nil {
			wp.service.logger.Printf("[Worker %d] ❌ Ошибка: %v", workerID, err)
		} else {
			wp.service.logger.Printf("[Worker %d] ✅ Обработано записей: %d", workerID, count)
		}

		// Отправляем результат
		wp.results <- TournamentResult{
			Tournament: task.Tournament,
			Count:      count,
			Error:      err,
		}
	}
}

// AddTask добавляет задачу в очередь
func (wp *TournamentWorkerPool) AddTask(task TournamentTask) {
	wp.tasks <- task
}

// Close закрывает очередь задач (больше задач не будет)
func (wp *TournamentWorkerPool) Close() {
	close(wp.tasks)
}

// Wait ждет завершения всех воркеров и закрывает канал результатов
func (wp *TournamentWorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results возвращает канал результатов
func (wp *TournamentWorkerPool) Results() <-chan TournamentResult {
	return wp.results
}
