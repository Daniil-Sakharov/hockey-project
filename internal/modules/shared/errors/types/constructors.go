package types

import "time"

// NewDomainError создает новую доменную ошибку
func NewDomainError(errorType ErrorType, code, message string) *DomainError {
	return &DomainError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: isRetryableByType(errorType),
		Context:   make(map[string]interface{}),
	}
}

// WrapError оборачивает существующую ошибку в DomainError
func WrapError(err error, errorType ErrorType, code, message string) *DomainError {
	return &DomainError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Cause:     err,
		Timestamp: time.Now(),
		Retryable: isRetryableByType(errorType),
		Context:   make(map[string]interface{}),
	}
}

// isRetryableByType определяет можно ли повторить операцию по типу ошибки
func isRetryableByType(errorType ErrorType) bool {
	switch errorType {
	case ErrorTypeParsingTemporary, ErrorTypeNetwork, ErrorTypeExternal:
		return true
	case ErrorTypeParsingPermanent, ErrorTypeBusiness, ErrorTypeValidation:
		return false
	case ErrorTypeInfrastructure, ErrorTypeDatabase:
		return true // зависит от конкретной ошибки
	default:
		return false
	}
}
