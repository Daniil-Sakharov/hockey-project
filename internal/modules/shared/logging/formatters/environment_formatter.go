package formatters

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// EnvironmentFormatter форматтер зависящий от окружения
type EnvironmentFormatter struct {
	environment string
}

// NewEnvironmentFormatter создает форматтер для окружения
func NewEnvironmentFormatter(env string) *EnvironmentFormatter {
	return &EnvironmentFormatter{environment: env}
}

// FormatMessage форматирует сообщение для конкретного окружения
func (f *EnvironmentFormatter) FormatMessage(ctx context.Context, msg string, fields ...zap.Field) {
	// Добавляем информацию об окружении
	envField := zap.String("environment", f.environment)
	allFields := append(fields, envField)

	// Используем pkg/logger для вывода
	logger.Info(ctx, msg, allFields...)
}
