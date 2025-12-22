package pool

import (
	"context"
	"time"
)

// Worker воркер для выполнения задач
type Worker struct {
	id      int
	tasks   <-chan Task
	results chan<- Result
	metrics *PoolMetrics
}

// NewWorker создает нового воркера
func NewWorker(id int, tasks <-chan Task, results chan<- Result, metrics *PoolMetrics) *Worker {
	return &Worker{
		id:      id,
		tasks:   tasks,
		results: results,
		metrics: metrics,
	}
}

// Run запускает воркера
func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case task, ok := <-w.tasks:
			if !ok {
				return // Канал закрыт
			}

			w.processTask(ctx, task)

		case <-ctx.Done():
			return
		}
	}
}

// processTask обрабатывает задачу
func (w *Worker) processTask(ctx context.Context, task Task) {
	start := time.Now()

	// Уменьшаем счетчик очереди
	w.metrics.QueuedTasks.Add(ctx, -1)

	// Выполняем задачу
	result := task.Execute(ctx)

	// Записываем метрики
	duration := time.Since(start)
	w.metrics.ProcessedTasks.Add(ctx, 1)

	if w.metrics.TaskDuration != nil {
		w.metrics.TaskDuration.Record(ctx, duration.Seconds())
	}

	// Отправляем результат
	select {
	case w.results <- result:
	case <-ctx.Done():
	}
}
