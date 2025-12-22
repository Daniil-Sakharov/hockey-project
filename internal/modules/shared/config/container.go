package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/providers"
)

// Container DI контейнер для конфигурации
type Container struct {
	mu       sync.RWMutex
	provider providers.ConfigProvider
	configs  map[string]interface{}
}

// NewContainer создает новый контейнер конфигурации
func NewContainer() *Container {
	envSource := providers.NewEnvSource("")
	provider := providers.NewConfigProvider(envSource)

	return &Container{
		provider: provider,
		configs:  make(map[string]interface{}),
	}
}

// Database возвращает конфигурацию базы данных
func (c *Container) Database(ctx context.Context) (*modules.DatabaseConfig, error) {
	config, err := c.getOrLoad(ctx, "database", &modules.DatabaseConfig{})
	if err != nil {
		return nil, err
	}
	return config.(*modules.DatabaseConfig), nil
}

// Redis возвращает конфигурацию Redis
func (c *Container) Redis(ctx context.Context) (*modules.RedisConfig, error) {
	config, err := c.getOrLoad(ctx, "redis", &modules.RedisConfig{})
	if err != nil {
		return nil, err
	}
	return config.(*modules.RedisConfig), nil
}

// Parsing возвращает конфигурацию парсинга
func (c *Container) Parsing(ctx context.Context) (*modules.ParsingConfig, error) {
	config, err := c.getOrLoad(ctx, "parsing", &modules.ParsingConfig{})
	if err != nil {
		return nil, err
	}
	return config.(*modules.ParsingConfig), nil
}

// Telegram возвращает конфигурацию Telegram
func (c *Container) Telegram(ctx context.Context) (*modules.TelegramConfig, error) {
	config, err := c.getOrLoad(ctx, "telegram", &modules.TelegramConfig{})
	if err != nil {
		return nil, err
	}
	return config.(*modules.TelegramConfig), nil
}

// getOrLoad получает конфигурацию из кэша или загружает новую
func (c *Container) getOrLoad(ctx context.Context, key string, target interface{}) (interface{}, error) {
	c.mu.RLock()
	if config, exists := c.configs[key]; exists {
		c.mu.RUnlock()
		return config, nil
	}
	c.mu.RUnlock()

	if err := c.provider.Load(ctx, target); err != nil {
		return nil, fmt.Errorf("failed to load %s config: %w", key, err)
	}

	c.mu.Lock()
	c.configs[key] = target
	c.mu.Unlock()

	return target, nil
}

// Reload перезагружает все конфигурации
func (c *Container) Reload(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.configs = make(map[string]interface{})
	return c.provider.Reload(ctx)
}
