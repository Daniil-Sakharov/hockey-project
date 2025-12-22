package types

import (
	"testing"
)

func TestDomainError_Basic(t *testing.T) {
	err := NewDomainError(ErrorTypeParsingTemporary, "PARSING_ERROR", "test error")

	if err.Error() != "[parsing_temporary:PARSING_ERROR] test error" {
		t.Errorf("Expected '[parsing_temporary:PARSING_ERROR] test error', got %s", err.Error())
	}

	if !err.IsRetryable() {
		t.Error("Expected parsing temporary error to be retryable")
	}

	if err.Code != "PARSING_ERROR" {
		t.Errorf("Expected 'PARSING_ERROR', got %s", err.Code)
	}
}

func TestDomainError_WithContext(t *testing.T) {
	err := NewDomainError(ErrorTypeBusiness, "INVALID_DATA", "business rule violated")
	_ = err.WithContext("player_id", "123")

	if err.GetContext()["player_id"] != "123" {
		t.Error("Expected context to be set")
	}

	if err.IsRetryable() {
		t.Error("Expected business error to not be retryable")
	}
}

func TestDomainError_WithTraceID(t *testing.T) {
	err := NewDomainError(ErrorTypeNetwork, "NETWORK_ERROR", "connection failed")
	_ = err.WithTraceID("trace-123")

	if err.TraceID != "trace-123" {
		t.Errorf("Expected 'trace-123', got %s", err.TraceID)
	}
}
