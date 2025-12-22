package outputs

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// OutputManager управляет выводом логов
type OutputManager struct {
	outputs []Output
}

// Output интерфейс вывода
type Output interface {
	Write(ctx context.Context, msg string, fields ...zap.Field)
}

// ConsoleOutput вывод в консоль
type ConsoleOutput struct{}

// JSONOutput вывод в JSON формате
type JSONOutput struct{}

// OTELOutput вывод в OpenTelemetry
type OTELOutput struct{}

// NewOutputManager создает менеджер выводов
func NewOutputManager() *OutputManager {
	return &OutputManager{
		outputs: make([]Output, 0),
	}
}

// AddOutput добавляет вывод
func (m *OutputManager) AddOutput(output Output) {
	m.outputs = append(m.outputs, output)
}

// Write записывает во все выводы
func (m *OutputManager) Write(ctx context.Context, msg string, fields ...zap.Field) {
	for _, output := range m.outputs {
		output.Write(ctx, msg, fields...)
	}
}

// Write для ConsoleOutput
func (c *ConsoleOutput) Write(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(ctx, msg, fields...)
}

// Write для JSONOutput
func (j *JSONOutput) Write(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(ctx, msg, fields...)
}

// Write для OTELOutput
func (o *OTELOutput) Write(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(ctx, msg, fields...)
}
