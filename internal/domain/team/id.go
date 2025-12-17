package team

import (
	"crypto/sha256"
	"fmt"
	"regexp"
)

// ID представляет уникальный идентификатор команды (Value Object)
type ID string

// NewID создает новый ID из строки
func NewID(id string) (ID, error) {
	if id == "" {
		return "", fmt.Errorf("team id cannot be empty")
	}
	return ID(id), nil
}

// String возвращает строковое представление
func (id ID) String() string {
	return string(id)
}

// IsEmpty проверяет пустой ли ID
func (id ID) IsEmpty() bool {
	return id == ""
}

// ExtractFromURL извлекает ID из URL команды
// Пример: /tournaments/pervenstvo/.../buran_5136295/ -> "5136295"
func ExtractFromURL(url string) ID {
	re := regexp.MustCompile(`_(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return ID(matches[1])
	}
	// Fallback: используем hash от URL
	return ID(generateHashID(url))
}

func generateHashID(url string) string {
	hash := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", hash)[:8]
}
