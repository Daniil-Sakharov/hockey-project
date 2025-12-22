package types

import (
	"fmt"
	"time"
)

// ErrorType тип ошибки для категоризации
type ErrorType string

const (
	// Parsing errors
	ErrorTypeParsingTemporary ErrorType = "parsing_temporary"
	ErrorTypeParsingPermanent ErrorType = "parsing_permanent"

	// Business logic errors
	ErrorTypeBusiness   ErrorType = "business"
	ErrorTypeValidation ErrorType = "validation"

	// Infrastructure errors
	ErrorTypeInfrastructure ErrorType = "infrastructure"
	ErrorTypeDatabase       ErrorType = "database"
	ErrorTypeNetwork        ErrorType = "network"
	ErrorTypeExternal       ErrorType = "external"
)

// DomainError базовая структура для всех доменных ошибок
type DomainError struct {
	Type      ErrorType              `json:"type"`
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Cause     error                  `json:"-"`
	Timestamp time.Time              `json:"timestamp"`
	Retryable bool                   `json:"retryable"`
	TraceID   string                 `json:"trace_id,omitempty"`
}

// Error реализует интерфейс error
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Type, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)
}

// Unwrap возвращает причину ошибки для errors.Unwrap
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// Is проверяет является ли ошибка определенного типа
func (e *DomainError) Is(target error) bool {
	if t, ok := target.(*DomainError); ok {
		return e.Type == t.Type && e.Code == t.Code
	}
	return false
}

// WithContext добавляет контекст к ошибке
func (e *DomainError) WithContext(key string, value interface{}) *DomainError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// WithTraceID добавляет trace ID к ошибке
func (e *DomainError) WithTraceID(traceID string) *DomainError {
	e.TraceID = traceID
	return e
}

// IsRetryable проверяет можно ли повторить операцию
func (e *DomainError) IsRetryable() bool {
	return e.Retryable
}

// GetContext возвращает контекст ошибки
func (e *DomainError) GetContext() map[string]interface{} {
	return e.Context
}
