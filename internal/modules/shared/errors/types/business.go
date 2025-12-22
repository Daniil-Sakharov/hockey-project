package types

import "fmt"

// Business error codes
const (
	CodeBusinessValidationFailed = "BUSINESS_VALIDATION_FAILED"
	CodeBusinessRuleViolation    = "BUSINESS_RULE_VIOLATION"
	CodeBusinessDataInconsistent = "BUSINESS_DATA_INCONSISTENT"
	CodeBusinessResourceNotFound = "BUSINESS_RESOURCE_NOT_FOUND"
)

// NewNotFoundError создает ошибку "не найдено"
func NewNotFoundError(resource, identifier string) *DomainError {
	return NewDomainError(
		ErrorTypeParsingPermanent,
		CodeParsingNotFound,
		fmt.Sprintf("%s not found: %s", resource, identifier),
	).WithContext("resource", resource).WithContext("identifier", identifier)
}

// NewInvalidFormatError создает ошибку неверного формата
func NewInvalidFormatError(field, expected, actual string) *DomainError {
	return NewDomainError(
		ErrorTypeParsingPermanent,
		CodeParsingInvalidFormat,
		fmt.Sprintf("Invalid format for %s: expected %s, got %s", field, expected, actual),
	).WithContext("field", field).WithContext("expected", expected).WithContext("actual", actual)
}

// NewAccessDeniedError создает ошибку отказа в доступе
func NewAccessDeniedError(resource string) *DomainError {
	return NewDomainError(
		ErrorTypeParsingPermanent,
		CodeParsingAccessDenied,
		fmt.Sprintf("Access denied to %s", resource),
	).WithContext("resource", resource)
}

// NewValidationError создает ошибку валидации
func NewValidationError(field, reason string) *DomainError {
	return NewDomainError(
		ErrorTypeValidation,
		CodeBusinessValidationFailed,
		fmt.Sprintf("Validation failed for %s: %s", field, reason),
	).WithContext("field", field).WithContext("reason", reason)
}

// NewBusinessRuleError создает ошибку нарушения бизнес-правила
func NewBusinessRuleError(rule, reason string) *DomainError {
	return NewDomainError(
		ErrorTypeBusiness,
		CodeBusinessRuleViolation,
		fmt.Sprintf("Business rule '%s' violated: %s", rule, reason),
	).WithContext("rule", rule).WithContext("reason", reason)
}
