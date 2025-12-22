package handlers

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestRetryHandler_HandleWithRetry(t *testing.T) {
	logger := zap.NewNop()
	config := Config{MaxRetries: 3, BaseDelay: time.Millisecond}
	handler := NewRetryHandler(logger, config)

	// Тест успешного выполнения без ошибки
	err := handler.HandleWithRetry(context.Background(), nil, nil)
	if err != nil {
		t.Errorf("Expected nil for no error, got %v", err)
	}
}
