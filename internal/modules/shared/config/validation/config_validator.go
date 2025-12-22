package validation

import (
	"fmt"
	"strings"
)

// ValidationError ошибка валидации конфигурации
type ValidationError struct {
	Field   string
	Message string
}

// Error реализует интерфейс error
func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation error: %s - %s", e.Field, e.Message)
}

// ConfigValidator валидатор конфигурации
type ConfigValidator struct {
	errors []ValidationError
}

// NewConfigValidator создает валидатор конфигурации
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		errors: make([]ValidationError, 0),
	}
}

// ValidateRequired проверяет обязательное поле
func (v *ConfigValidator) ValidateRequired(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "required field is empty",
		})
	}
}

// ValidatePort проверяет порт
func (v *ConfigValidator) ValidatePort(field string, port int) {
	if port <= 0 || port > 65535 {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "port must be between 1 and 65535",
		})
	}
}

// HasErrors проверяет наличие ошибок
func (v *ConfigValidator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors возвращает все ошибки
func (v *ConfigValidator) GetErrors() []ValidationError {
	return v.errors
}
