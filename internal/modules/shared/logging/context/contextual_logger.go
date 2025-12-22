package context

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// ContextualLogger wraps pkg/logger with module context.
type ContextualLogger struct {
	module string
	fields []zap.Field
}

// NewContextualLogger creates a logger for the given module.
func NewContextualLogger(module string) *ContextualLogger {
	return &ContextualLogger{
		module: module,
		fields: []zap.Field{zap.String("module", module)},
	}
}

// WithField adds a field to the logger.
func (l *ContextualLogger) WithField(key string, value interface{}) *ContextualLogger {
	newFields := make([]zap.Field, len(l.fields))
	copy(newFields, l.fields)
	newFields = append(newFields, zap.Any(key, value))
	return &ContextualLogger{module: l.module, fields: newFields}
}

// Info logs info message using pkg/logger.
func (l *ContextualLogger) Info(ctx context.Context, msg string) {
	logger.Info(ctx, msg, l.fields...)
}

// Error logs error message using pkg/logger.
func (l *ContextualLogger) Error(ctx context.Context, msg string, err error) {
	fields := append(l.fields, zap.Error(err))
	logger.Error(ctx, msg, fields...)
}

// Warn logs warning message using pkg/logger.
func (l *ContextualLogger) Warn(ctx context.Context, msg string) {
	logger.Warn(ctx, msg, l.fields...)
}
