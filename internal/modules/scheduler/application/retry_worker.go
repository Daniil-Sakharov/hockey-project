package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// RetryWorker обрабатывает неудачные задачи парсинга
type RetryWorker struct {
	failedJobRepo *repositories.FailedJobRepository
	handlers      map[string]RetryHandler
}

// RetryHandler функция для повторной попытки
type RetryHandler func(ctx context.Context, job *repositories.FailedJob) error

// NewRetryWorker создаёт новый retry worker
func NewRetryWorker(failedJobRepo *repositories.FailedJobRepository) *RetryWorker {
	return &RetryWorker{
		failedJobRepo: failedJobRepo,
		handlers:      make(map[string]RetryHandler),
	}
}

// RegisterHandler регистрирует обработчик для типа задачи
func (w *RetryWorker) RegisterHandler(jobType string, handler RetryHandler) {
	w.handlers[jobType] = handler
}

// Process обрабатывает pending задачи
func (w *RetryWorker) Process(ctx context.Context) error {
	jobs, err := w.failedJobRepo.GetPendingRetries(ctx, 100)
	if err != nil {
		return fmt.Errorf("get pending retries: %w", err)
	}

	if len(jobs) == 0 {
		logger.Info(ctx, "No pending retry jobs")
		return nil
	}

	logger.Info(ctx, fmt.Sprintf("Processing %d retry jobs", len(jobs)))

	var processed, succeeded, failed int

	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		handler, ok := w.handlers[job.JobType]
		if !ok {
			logger.Warn(ctx, "No handler for job type: "+job.JobType)
			continue
		}

		processed++
		err := handler(ctx, job)

		if err == nil {
			// Успех - удаляем задачу
			if delErr := w.failedJobRepo.Delete(ctx, job.ID); delErr != nil {
				logger.Error(ctx, "Failed to delete job: "+delErr.Error())
			}
			succeeded++
		} else {
			// Ошибка - увеличиваем счётчик
			nextRetry := time.Now().Add(repositories.GetRetryInterval(job.RetryCount + 1))
			if updErr := w.failedJobRepo.IncrementRetry(ctx, job.ID, nextRetry, err.Error()); updErr != nil {
				logger.Error(ctx, "Failed to update job: "+updErr.Error())
			}
			failed++
		}
	}

	logger.Info(ctx, fmt.Sprintf("Retry complete: processed=%d, succeeded=%d, failed=%d", processed, succeeded, failed))
	return nil
}

// Run запускает обработку как задачу scheduler
func (w *RetryWorker) Run() error {
	ctx := context.Background()
	return w.Process(ctx)
}

// Cleanup удаляет старые записи
func (w *RetryWorker) Cleanup(ctx context.Context, olderThan time.Duration) error {
	deleted, err := w.failedJobRepo.CleanupOld(ctx, olderThan)
	if err != nil {
		return err
	}
	if deleted > 0 {
		logger.Info(ctx, fmt.Sprintf("Cleaned up %d old failed jobs", deleted))
	}
	return nil
}

// Stats возвращает статистику
func (w *RetryWorker) Stats(ctx context.Context) (pending, failed int, err error) {
	pending, err = w.failedJobRepo.CountPending(ctx)
	if err != nil {
		return 0, 0, err
	}
	failed, err = w.failedJobRepo.CountFailed(ctx)
	return pending, failed, err
}
