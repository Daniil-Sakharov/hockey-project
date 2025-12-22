package handlers

import (
	"context"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/errors/types"
)

// ErrorHandler интерфейс для обработки ошибок
type ErrorHandler interface {
	Handle(ctx context.Context, err error) error
	HandleWithRetry(ctx context.Context, err error, operation RetryableOperation) error
	ShouldRetry(err error) bool
	GetRetryDelay(err error, attempt int) time.Duration
}

// RetryableOperation операция которую можно повторить
type RetryableOperation func(ctx context.Context) error

// Config конфигурация для ErrorHandler
type Config struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Multiplier float64
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		MaxRetries: 3,
		BaseDelay:  time.Second,
		MaxDelay:   time.Minute * 5,
		Multiplier: 2.0,
	}
}

// ShouldRetry определяет нужно ли повторить операцию
func ShouldRetry(err error) bool {
	if domainErr, ok := err.(*types.DomainError); ok {
		return domainErr.IsRetryable()
	}
	return false
}

// GetRetryDelay вычисляет задержку перед повторной попыткой
func GetRetryDelay(baseDelay, maxDelay time.Duration, multiplier float64, attempt int) time.Duration {
	delay := time.Duration(float64(baseDelay) * pow(multiplier, float64(attempt-1)))

	if delay > maxDelay {
		delay = maxDelay
	}

	// Добавляем jitter (±25%)
	jitter := time.Duration(float64(delay) * 0.25)
	jitterMultiplier := float64(2*time.Now().UnixNano()%2 - 1)
	delay = delay + time.Duration(float64(jitter)*jitterMultiplier)

	return delay
}

// pow простая реализация возведения в степень для float64
func pow(base, exp float64) float64 {
	if exp == 0 {
		return 1
	}
	result := base
	for i := 1; i < int(exp); i++ {
		result *= base
	}
	return result
}
