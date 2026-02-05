package pool

import (
	"context"
	"sync"
	"time"
)

// WorkerPool пул воркеров для параллельного парсинга сезонов
type WorkerPool struct {
	parser      SeasonParser // Интерфейс для парсинга
	tasks       chan SeasonTask
	results     chan WorkerPoolResult
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
}

// NewWorkerPool создает новый Worker Pool
func NewWorkerPool(ctx context.Context, parser SeasonParser, workerCount int) *WorkerPool {
	return &WorkerPool{
		parser:      parser,
		tasks:       make(chan SeasonTask, workerCount*2), // Буфер для задач
		results:     make(chan WorkerPoolResult, workerCount*2),
		workerCount: workerCount,
		ctx:         ctx,
	}
}

// Start запускает воркеры
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker функция каждого воркера
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.tasks {
		// Проверка контекста (отмена)
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		// Парсим турниры сезона через интерфейс
		tournaments, err := wp.parser.ParseSeasonTournaments(
			task.Domain,
			task.Season,
			task.AjaxURL,
		)

		// Отправляем результат
		wp.results <- WorkerPoolResult{
			Tournaments: tournaments,
			Error:       err,
			Task:        task,
		}

		// Rate limiting: 500ms между запросами (щадящий режим для серверов)
		time.Sleep(500 * time.Millisecond)
	}
}

// AddTask добавляет задачу в очередь
func (wp *WorkerPool) AddTask(task SeasonTask) {
	wp.tasks <- task
}

// Close закрывает очередь задач (больше задач не будет)
func (wp *WorkerPool) Close() {
	close(wp.tasks)
}

// Wait ждет завершения всех воркеров и закрывает канал результатов
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results возвращает канал результатов
func (wp *WorkerPool) Results() <-chan WorkerPoolResult {
	return wp.results
}
