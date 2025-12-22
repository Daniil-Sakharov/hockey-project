package factory

import (
	"context"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/workers/pool"
)

// PoolFactory фабрика для создания пулов воркеров
type PoolFactory struct{}

// NewPoolFactory создает новую фабрику
func NewPoolFactory() *PoolFactory {
	return &PoolFactory{}
}

// CreateParsingPool создает пул для парсинга
func (f *PoolFactory) CreateParsingPool(ctx context.Context, workerCount int) (*pool.ConfigurablePool, error) {
	config := pool.Config{
		Name:           "parsing-pool",
		WorkerCount:    workerCount,
		MaxWorkers:     workerCount * 2,
		BufferSize:     workerCount * 4,
		TaskTimeout:    30 * time.Second,
		ScaleThreshold: 0.8,
		ScaleInterval:  10 * time.Second,
	}

	return pool.NewConfigurablePool(ctx, config)
}

// CreateTeamPool создает пул для обработки команд
func (f *PoolFactory) CreateTeamPool(ctx context.Context, workerCount int) (*pool.ConfigurablePool, error) {
	config := pool.Config{
		Name:           "team-pool",
		WorkerCount:    workerCount,
		MaxWorkers:     workerCount * 2,
		BufferSize:     workerCount * 3,
		TaskTimeout:    20 * time.Second,
		ScaleThreshold: 0.75,
		ScaleInterval:  15 * time.Second,
	}

	return pool.NewConfigurablePool(ctx, config)
}

// CreatePlayerPool создает пул для обработки игроков
func (f *PoolFactory) CreatePlayerPool(ctx context.Context, workerCount int) (*pool.ConfigurablePool, error) {
	config := pool.Config{
		Name:           "player-pool",
		WorkerCount:    workerCount,
		MaxWorkers:     workerCount * 3,
		BufferSize:     workerCount * 5,
		TaskTimeout:    15 * time.Second,
		ScaleThreshold: 0.85,
		ScaleInterval:  5 * time.Second,
	}

	return pool.NewConfigurablePool(ctx, config)
}

// CreateCustomPool создает пул с кастомной конфигурацией
func (f *PoolFactory) CreateCustomPool(ctx context.Context, config pool.Config) (*pool.ConfigurablePool, error) {
	return pool.NewConfigurablePool(ctx, config)
}
