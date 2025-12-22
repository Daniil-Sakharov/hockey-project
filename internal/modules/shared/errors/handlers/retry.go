package handlers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// retryHandler реализация retry логики
type retryHandler struct {
	logger      *zap.Logger
	config      Config
	errorLogger *errorLogger
}

// NewRetryHandler создает новый retry handler
func NewRetryHandler(logger *zap.Logger, config Config) *retryHandler {
	return &retryHandler{
		logger:      logger,
		config:      config,
		errorLogger: NewErrorLogger(logger),
	}
}

// HandleWithRetry обрабатывает ошибку с повторными попытками
func (h *retryHandler) HandleWithRetry(ctx context.Context, err error, operation RetryableOperation) error {
	if err == nil {
		return nil
	}

	if !ShouldRetry(err) {
		return h.errorLogger.Handle(ctx, err)
	}

	// Выполняем повторные попытки
	lastErr := err
	for attempt := 1; attempt <= h.config.MaxRetries; attempt++ {
		delay := GetRetryDelay(h.config.BaseDelay, h.config.MaxDelay, h.config.Multiplier, attempt)

		h.logger.Info("Retrying operation",
			zap.Int("attempt", attempt),
			zap.Int("max_retries", h.config.MaxRetries),
			zap.Duration("delay", delay),
			zap.String("error", lastErr.Error()),
		)

		// Ждем перед повторной попыткой
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		// Выполняем операцию
		if retryErr := operation(ctx); retryErr == nil {
			h.logger.Info("Operation succeeded after retry",
				zap.Int("attempt", attempt),
			)
			return nil
		} else {
			lastErr = retryErr
		}
	}

	// Все попытки исчерпаны
	h.logger.Error("All retry attempts failed",
		zap.Int("max_retries", h.config.MaxRetries),
		zap.String("final_error", lastErr.Error()),
	)

	return h.errorLogger.Handle(ctx, fmt.Errorf("operation failed after %d retries: %w", h.config.MaxRetries, lastErr))
}

// ShouldRetry определяет нужно ли повторить операцию
func (h *retryHandler) ShouldRetry(err error) bool {
	return ShouldRetry(err)
}

// GetRetryDelay вычисляет задержку перед повторной попыткой
func (h *retryHandler) GetRetryDelay(err error, attempt int) time.Duration {
	return GetRetryDelay(h.config.BaseDelay, h.config.MaxDelay, h.config.Multiplier, attempt)
}
