package di

import (
	"context"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events/store"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/cache"
	"github.com/jmoiron/sqlx"
)

// Container модульный DI контейнер
type Container struct {
	configContainer *config.Container
	mu              sync.RWMutex

	// Infrastructure (lazy initialized)
	db          *sqlx.DB
	redisClient cache.RedisClient
	eventStore  *store.SQLXEventStore
}

// NewContainer создает новый модульный DI контейнер
func NewContainer() *Container {
	return &Container{
		configContainer: config.NewContainer(),
	}
}

// NewContainerWithConfig создает контейнер с существующим config container
func NewContainerWithConfig(cfg *config.Container) *Container {
	return &Container{
		configContainer: cfg,
	}
}

// Config возвращает контейнер конфигурации
func (c *Container) Config() *config.Container {
	return c.configContainer
}

// Close закрывает все ресурсы
func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// DB возвращает подключение к PostgreSQL
func (c *Container) DB(ctx context.Context) (*sqlx.DB, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.db != nil {
		return c.db, nil
	}

	dbConfig, err := c.configContainer.Database(ctx)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", dbConfig.URI())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(dbConfig.ConnMaxIdleTime)

	c.db = db
	return c.db, nil
}
