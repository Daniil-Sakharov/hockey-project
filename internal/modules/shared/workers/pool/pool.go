package pool

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// ConfigurablePool унифицированный пул воркеров с адаптивным масштабированием
type ConfigurablePool struct {
	name        string
	workerCount int
	maxWorkers  int
	tasks       chan Task
	results     chan Result
	workers     []*Worker
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	metrics     *PoolMetrics
	config      Config
}

// Config конфигурация пула
type Config struct {
	Name           string
	WorkerCount    int
	MaxWorkers     int
	BufferSize     int
	TaskTimeout    time.Duration
	ScaleThreshold float64 // 0.8 = масштабировать при 80% загрузке
	ScaleInterval  time.Duration
}

// PoolMetrics метрики пула для OTEL
type PoolMetrics struct {
	ActiveWorkers   metric.Int64UpDownCounter
	QueuedTasks     metric.Int64UpDownCounter
	ProcessedTasks  metric.Int64Counter
	TaskDuration    metric.Float64Histogram
	PoolUtilization metric.Float64Gauge
}

// NewConfigurablePool создает новый пул воркеров
func NewConfigurablePool(ctx context.Context, config Config) (*ConfigurablePool, error) {
	poolCtx, cancel := context.WithCancel(ctx)

	metrics, err := initMetrics(config.Name)
	if err != nil {
		cancel()
		return nil, err
	}

	pool := &ConfigurablePool{
		name:        config.Name,
		workerCount: config.WorkerCount,
		maxWorkers:  config.MaxWorkers,
		tasks:       make(chan Task, config.BufferSize),
		results:     make(chan Result, config.BufferSize),
		ctx:         poolCtx,
		cancel:      cancel,
		metrics:     metrics,
		config:      config,
	}

	return pool, nil
}

// Start запускает пул воркеров
func (p *ConfigurablePool) Start() {
	p.startWorkers(p.workerCount)

	// Запускаем адаптивное масштабирование
	if p.config.ScaleInterval > 0 {
		go p.adaptiveScaling()
	}
}

// startWorkers запускает указанное количество воркеров
func (p *ConfigurablePool) startWorkers(count int) {
	for i := 0; i < count; i++ {
		worker := NewWorker(i, p.tasks, p.results, p.metrics)
		p.workers = append(p.workers, worker)

		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Run(p.ctx)
		}(worker)
	}

	p.metrics.ActiveWorkers.Add(p.ctx, int64(count))
}

// Submit отправляет задачу в пул
func (p *ConfigurablePool) Submit(task Task) {
	select {
	case p.tasks <- task:
		p.metrics.QueuedTasks.Add(p.ctx, 1)
	case <-p.ctx.Done():
		// Пул закрыт
	}
}

// Results возвращает канал результатов
func (p *ConfigurablePool) Results() <-chan Result {
	return p.results
}

// Close закрывает пул и ждет завершения всех воркеров
func (p *ConfigurablePool) Close() {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
	p.cancel()
}

// initMetrics инициализирует OTEL метрики
func initMetrics(poolName string) (*PoolMetrics, error) {
	meter := otel.Meter("workers.pool")

	activeWorkers, err := meter.Int64UpDownCounter(
		"worker_pool.active_workers",
		metric.WithDescription("Number of active workers"),
	)
	if err != nil {
		return nil, err
	}

	queuedTasks, err := meter.Int64UpDownCounter(
		"worker_pool.queued_tasks",
		metric.WithDescription("Number of queued tasks"),
	)
	if err != nil {
		return nil, err
	}

	processedTasks, err := meter.Int64Counter(
		"worker_pool.processed_tasks",
		metric.WithDescription("Total processed tasks"),
	)
	if err != nil {
		return nil, err
	}

	return &PoolMetrics{
		ActiveWorkers:  activeWorkers,
		QueuedTasks:    queuedTasks,
		ProcessedTasks: processedTasks,
	}, nil
}
