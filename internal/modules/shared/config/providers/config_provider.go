package providers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// ConfigProvider интерфейс для загрузки и валидации конфигурации
type ConfigProvider interface {
	Load(ctx context.Context, target interface{}) error
	Validate(ctx context.Context, config interface{}) error
	Reload(ctx context.Context) error
}

// configProvider реализация ConfigProvider
type configProvider struct {
	validator *validator.Validate
	sources   []ConfigSource
}

// ConfigSource интерфейс для источников конфигурации
type ConfigSource interface {
	Load(ctx context.Context, target interface{}) error
	Name() string
}

// NewConfigProvider создает новый ConfigProvider
func NewConfigProvider(sources ...ConfigSource) ConfigProvider {
	return &configProvider{
		validator: validator.New(),
		sources:   sources,
	}
}

// Load загружает конфигурацию из всех источников
func (p *configProvider) Load(ctx context.Context, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	// Проверяем что target это указатель на структуру
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	// Загружаем из всех источников по порядку
	for _, source := range p.sources {
		if err := source.Load(ctx, target); err != nil {
			return fmt.Errorf("failed to load from %s: %w", source.Name(), err)
		}
	}

	// Валидируем загруженную конфигурацию
	return p.Validate(ctx, target)
}

// Validate валидирует конфигурацию
func (p *configProvider) Validate(ctx context.Context, config interface{}) error {
	if err := p.validator.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

// Reload перезагружает конфигурацию
func (p *configProvider) Reload(ctx context.Context) error {
	// TODO: Реализовать hot reload через fsnotify
	return fmt.Errorf("hot reload not implemented yet")
}
