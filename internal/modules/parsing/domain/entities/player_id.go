package entities

import (
	"fmt"
	"regexp"
)

// ID представляет уникальный идентификатор игрока (Value Object)
type PlayerID string

// NewID создает новый ID из строки
func NewPlayerID(id string) PlayerID {
	return PlayerID(id)
}

// String возвращает строковое представление
func (id PlayerID) String() string {
	return string(id)
}

// IsEmpty проверяет пустой ли ID
func (id PlayerID) IsEmpty() bool {
	return id == ""
}

// ExtractFromURL извлекает ID из URL профиля
// Пример: /player/abdrashitov-daniyar-2008-05-13-924040/ -> "924040"
func ExtractPlayerIDFromURL(profileURL string) PlayerID {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(profileURL)
	if len(matches) > 1 {
		return PlayerID(matches[1])
	}
	return ""
}

// ExtractIDFromURL - deprecated, use ExtractFromURL
// Сохранено для обратной совместимости
func ExtractPlayerIDFromURLLegacy(profileURL string) string {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(profileURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Validate проверяет валидность ID
func (id PlayerID) Validate() error {
	if id == "" {
		return fmt.Errorf("player id cannot be empty")
	}
	return nil
}
