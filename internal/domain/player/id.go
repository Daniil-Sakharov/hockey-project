package player

import (
	"fmt"
	"regexp"
)

// ID представляет уникальный идентификатор игрока (Value Object)
type ID string

// NewID создает новый ID из строки
func NewID(id string) ID {
	return ID(id)
}

// String возвращает строковое представление
func (id ID) String() string {
	return string(id)
}

// IsEmpty проверяет пустой ли ID
func (id ID) IsEmpty() bool {
	return id == ""
}

// ExtractFromURL извлекает ID из URL профиля
// Пример: /player/abdrashitov-daniyar-2008-05-13-924040/ -> "924040"
func ExtractFromURL(profileURL string) ID {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(profileURL)
	if len(matches) > 1 {
		return ID(matches[1])
	}
	return ""
}

// ExtractIDFromURL - deprecated, use ExtractFromURL
// Сохранено для обратной совместимости
func ExtractIDFromURL(profileURL string) string {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(profileURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Validate проверяет валидность ID
func (id ID) Validate() error {
	if id == "" {
		return fmt.Errorf("player id cannot be empty")
	}
	return nil
}
