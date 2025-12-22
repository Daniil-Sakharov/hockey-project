package handlers

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/errors/types"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// errorLogger обработчик для логирования ошибок
type errorLogger struct {
	logger *zap.Logger
}

// NewErrorLogger создает новый error logger
func NewErrorLogger(logger *zap.Logger) *errorLogger {
	return &errorLogger{
		logger: logger,
	}
}

// Handle обрабатывает ошибку без повторных попыток
func (h *errorLogger) Handle(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// Получаем trace ID из контекста
	traceID := getTraceIDFromContext(ctx)

	// Если это уже DomainError, добавляем trace ID и логируем
	if domainErr, ok := err.(*types.DomainError); ok {
		_ = domainErr.WithTraceID(traceID)
		h.logDomainError(ctx, domainErr)
		return domainErr
	}

	// Оборачиваем обычную ошибку в DomainError
	domainErr := types.WrapError(err, types.ErrorTypeInfrastructure, "UNKNOWN_ERROR", "Unexpected error occurred")
	_ = domainErr.WithTraceID(traceID)
	h.logDomainError(ctx, domainErr)

	return domainErr
}

// logDomainError логирует доменную ошибку
func (h *errorLogger) logDomainError(ctx context.Context, err *types.DomainError) {
	fields := []zap.Field{
		zap.String("error_type", string(err.Type)),
		zap.String("error_code", err.Code),
		zap.String("error_message", err.Message),
		zap.Bool("retryable", err.Retryable),
		zap.Time("timestamp", err.Timestamp),
	}

	if err.TraceID != "" {
		fields = append(fields, zap.String("trace_id", err.TraceID))
	}

	if len(err.Context) > 0 {
		fields = append(fields, zap.Any("context", err.Context))
	}

	if err.Cause != nil {
		fields = append(fields, zap.Error(err.Cause))
	}

	// Выбираем уровень логирования в зависимости от типа ошибки
	switch err.Type {
	case types.ErrorTypeParsingTemporary, types.ErrorTypeNetwork:
		h.logger.Warn("Temporary error occurred", fields...)
	case types.ErrorTypeParsingPermanent, types.ErrorTypeBusiness, types.ErrorTypeValidation:
		h.logger.Error("Permanent error occurred", fields...)
	default:
		h.logger.Error("Error occurred", fields...)
	}
}

// getTraceIDFromContext извлекает trace ID из контекста
func getTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	return ""
}
