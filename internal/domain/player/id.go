package player

import (
	"fmt"
	"regexp"
)

// ID представляет уникальный идентификатор игрока (Value Object)
type ID string

// NewID создает новый ID из строки
func NewID(id string) (ID, error) {
	if id == "" {
		return "", fmt.Errorf("player id cannot be empty")
	}
	return ID(id), nil
}

// MustNewID создает ID или паникует
func MustNewID(id string) ID {
	pid, err := NewID(id)
	if err != nil {
		panic(err)
	}
	return pid
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
// Пример: /player/abdrashitov-daniyar-rustamovich-2008-05-13-924040/ -> "924040"
func ExtractFromURL(profileURL string) (ID, error) {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(profileURL)
	if len(matches) > 1 {
		return NewID(matches[1])
	}
	return "", fmt.Errorf("cannot extract id from url: %s", profileURL)
}
