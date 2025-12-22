package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/scheduler/infrastructure"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/google/uuid"
)

// SchedulerService сервис планировщика
type SchedulerService struct {
	config     *modules.SchedulerConfig
	cron       *infrastructure.CronAdapter
	lockRepo   *infrastructure.LockRepository
	metrics    *infrastructure.SchedulerMetrics
	instanceID string
	handlers   map[string]func() error
	mu         sync.RWMutex
	running    bool
}

// NewSchedulerService создаёт новый сервис планировщика
func NewSchedulerService(
	config *modules.SchedulerConfig,
	lockRepo *infrastructure.LockRepository,
	metrics *infrastructure.SchedulerMetrics,
) *SchedulerService {
	return &SchedulerService{
		config:     config,
		cron:       infrastructure.NewCronAdapter(),
		lockRepo:   lockRepo,
		metrics:    metrics,
		instanceID: uuid.New().String()[:8],
		handlers:   make(map[string]func() error),
	}
}

// RegisterHandler регистрирует обработчик для задачи
func (s *SchedulerService) RegisterHandler(jobName string, handler func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[jobName] = handler
}

// GetHandlers возвращает все зарегистрированные handlers
func (s *SchedulerService) GetHandlers() map[string]func() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]func() error)
	for k, v := range s.handlers {
		result[k] = v
	}
	return result
}

// Start запускает планировщик
func (s *SchedulerService) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	for name, jobCfg := range s.config.EnabledJobs() {
		handler, ok := s.handlers[name]
		if !ok {
			logger.Warn(ctx, "No handler for job: "+name)
			continue
		}

		if err := s.scheduleJob(ctx, name, jobCfg, handler); err != nil {
			return err
		}
	}

	s.cron.Start()
	s.running = true

	logger.Info(ctx, fmt.Sprintf("Scheduler started (instance: %s, jobs: %d)", s.instanceID, len(s.config.EnabledJobs())))
	return nil
}

// Stop останавливает планировщик
func (s *SchedulerService) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.cron.Stop()
	s.running = false

	// Освобождаем все локи этого инстанса
	if err := s.lockRepo.ReleaseAll(ctx, s.instanceID); err != nil {
		logger.Error(ctx, "Failed to release locks: "+err.Error())
	}

	logger.Info(ctx, "Scheduler stopped (instance: "+s.instanceID+")")
	return nil
}

// IsBootstrapMode возвращает true если включён режим первого запуска
func (s *SchedulerService) IsBootstrapMode() bool {
	return s.config.BootstrapMode
}

func (s *SchedulerService) scheduleJob(ctx context.Context, name string, cfg modules.JobConfig, handler func() error) error {
	wrappedHandler := s.wrapHandler(name, cfg.Timeout, handler)

	_, err := s.cron.AddJob(cfg.Cron, wrappedHandler)
	if err != nil {
		return err
	}

	logger.Info(ctx, fmt.Sprintf("Job scheduled: %s (cron: %s, timeout: %s)", name, cfg.Cron, cfg.Timeout))
	return nil
}

func (s *SchedulerService) wrapHandler(jobName string, timeout time.Duration, handler func() error) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Пытаемся получить блокировку
		acquired, err := s.lockRepo.TryAcquire(ctx, jobName, timeout, s.instanceID)
		if err != nil {
			logger.Error(ctx, "Failed to acquire lock: "+err.Error())
			if s.metrics != nil {
				s.metrics.RecordError(ctx, jobName, "lock_failed")
			}
			return
		}
		if !acquired {
			logger.Info(ctx, "Job already running, skipping: "+jobName)
			return
		}
		defer func() { _ = s.lockRepo.Release(ctx, jobName, s.instanceID) }()

		startedAt := time.Now()
		logger.Info(ctx, "Job started: "+jobName)

		err = handler()

		duration := time.Since(startedAt)
		success := err == nil

		// Записываем метрики
		if s.metrics != nil {
			s.metrics.RecordJobExecution(ctx, jobName, duration, success)
		}

		if success {
			logger.Info(ctx, "Job completed: "+jobName+" ("+duration.String()+")")
		} else {
			logger.Error(ctx, "Job failed: "+jobName+" ("+duration.String()+"): "+err.Error())
		}
	}
}

// GetMetrics возвращает метрики scheduler
func (s *SchedulerService) GetMetrics() *infrastructure.SchedulerMetrics {
	return s.metrics
}
