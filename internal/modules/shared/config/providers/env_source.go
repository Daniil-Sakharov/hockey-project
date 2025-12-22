package providers

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// envSource загружает конфигурацию из переменных окружения
type envSource struct {
	prefix string
}

// NewEnvSource создает новый источник из переменных окружения
func NewEnvSource(prefix string) ConfigSource {
	return &envSource{prefix: prefix}
}

// Name возвращает имя источника
func (s *envSource) Name() string {
	return fmt.Sprintf("env(%s)", s.prefix)
}

// Load загружает конфигурацию из переменных окружения
func (s *envSource) Load(ctx context.Context, target interface{}) error {
	rv := reflect.ValueOf(target).Elem()
	rt := rv.Type()

	return s.loadStruct(rv, rt, s.prefix)
}

// loadStruct рекурсивно загружает структуру
func (s *envSource) loadStruct(rv reflect.Value, rt reflect.Type, prefix string) error {
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		// Пропускаем неэкспортируемые поля
		if !fieldValue.CanSet() {
			continue
		}

		// Получаем тег env или используем имя поля
		envTag := field.Tag.Get("env")
		if envTag == "-" {
			continue
		}

		var envKey string
		if envTag != "" {
			envKey = envTag
		} else {
			envKey = s.buildEnvKey(prefix, field.Name)
		}

		// Обрабатываем вложенные структуры
		if fieldValue.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Duration(0)) {
			if err := s.loadStruct(fieldValue, field.Type, envKey); err != nil {
				return err
			}
			continue
		}

		// Загружаем значение из переменной окружения
		if err := s.setFieldFromEnv(fieldValue, field.Type, envKey, field.Tag.Get("default")); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

// buildEnvKey строит ключ переменной окружения
func (s *envSource) buildEnvKey(prefix, fieldName string) string {
	key := strings.ToUpper(fieldName)
	if prefix != "" {
		return fmt.Sprintf("%s_%s", strings.ToUpper(prefix), key)
	}
	return key
}

// setFieldFromEnv устанавливает значение поля из переменной окружения
func (s *envSource) setFieldFromEnv(fieldValue reflect.Value, fieldType reflect.Type, envKey, defaultValue string) error {
	envValue := os.Getenv(envKey)
	if envValue == "" && defaultValue != "" {
		envValue = defaultValue
	}

	if envValue == "" {
		return nil // Оставляем zero value
	}

	// Проверяем time.Duration
	if fieldType == reflect.TypeOf(time.Duration(0)) {
		duration, err := time.ParseDuration(envValue)
		if err != nil {
			return fmt.Errorf("invalid duration value for %s: %s", envKey, envValue)
		}
		fieldValue.SetInt(int64(duration))
		return nil
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(envValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int value for %s: %s", envKey, envValue)
		}
		fieldValue.SetInt(intVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(envValue)
		if err != nil {
			return fmt.Errorf("invalid bool value for %s: %s", envKey, envValue)
		}
		fieldValue.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}

	return nil
}
