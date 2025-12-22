package parser

import (
	"context"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// DomainTask задача для парсинга домена
type DomainTask struct {
	Domain string
	Index  int
	Total  int
}

// DomainResult результат парсинга домена
type DomainResult struct {
	Domain      string
	Tournaments []*entities.Tournament
	IsDuplicate bool
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
	globalDedup  *sync.Map
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
		select {
		case <-wp.ctx.Done():
			return
		default:
		}

		result := wp.processDomain(workerID, task)
		wp.results <- result
	}
}

// AddTask добавляет задачу в очередь
func (wp *DomainWorkerPool) AddTask(task DomainTask) {
	wp.tasks <- task
}

// Close закрывает очередь задач
func (wp *DomainWorkerPool) Close() {
	close(wp.tasks)
}

// Wait ждет завершения всех воркеров
func (wp *DomainWorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results возвращает канал результатов
func (wp *DomainWorkerPool) Results() <-chan DomainResult {
	return wp.results
}
