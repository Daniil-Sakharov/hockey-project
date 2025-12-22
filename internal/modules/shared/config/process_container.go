package config

import (
	"context"
	"fmt"
)

// ProcessType тип процесса
type ProcessType string

const (
	ProcessTypeBot    ProcessType = "bot"
	ProcessTypeParser ProcessType = "parser"
	ProcessTypeRetry  ProcessType = "retry"
	ProcessTypeWeb    ProcessType = "web"
)

// ProcessSpecificContainer контейнер для процесс-специфичных конфигураций
type ProcessSpecificContainer struct {
	*Container
	processType ProcessType
}

// NewProcessContainer создает контейнер для конкретного процесса
func NewProcessContainer(processType ProcessType) *ProcessSpecificContainer {
	return &ProcessSpecificContainer{
		Container:   NewContainer(),
		processType: processType,
	}
}

// GetRequiredConfigs возвращает список необходимых конфигураций для процесса
func (c *ProcessSpecificContainer) GetRequiredConfigs() []string {
	switch c.processType {
	case ProcessTypeBot:
		return []string{"database", "telegram", "redis"}
	case ProcessTypeParser:
		return []string{"database", "parsing"}
	case ProcessTypeRetry:
		return []string{"database", "parsing"}
	case ProcessTypeWeb:
		return []string{"database"}
	default:
		return []string{"database"}
	}
}

// ValidateRequired проверяет что все необходимые конфигурации загружены
func (c *ProcessSpecificContainer) ValidateRequired(ctx context.Context) error {
	required := c.GetRequiredConfigs()

	for _, configName := range required {
		if err := c.validateConfig(ctx, configName); err != nil {
			return err
		}
	}
	return nil
}

func (c *ProcessSpecificContainer) validateConfig(ctx context.Context, name string) error {
	switch name {
	case "database":
		if _, err := c.Database(ctx); err != nil {
			return fmt.Errorf("database config: %w", err)
		}
	case "parsing":
		if _, err := c.Parsing(ctx); err != nil {
			return fmt.Errorf("parsing config: %w", err)
		}
	case "telegram":
		if _, err := c.Telegram(ctx); err != nil {
			return fmt.Errorf("telegram config: %w", err)
		}
	case "redis":
		if _, err := c.Redis(ctx); err != nil {
			return fmt.Errorf("redis config: %w", err)
		}
	}
	return nil
}
