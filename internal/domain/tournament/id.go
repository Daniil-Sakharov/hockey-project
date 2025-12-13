package tournament

import (
	"fmt"
	"regexp"
)

// ID представляет уникальный идентификатор турнира (Value Object)
type ID string

// NewID создает новый ID из строки
func NewID(id string) (ID, error) {
	if id == "" {
		return "", fmt.Errorf("tournament id cannot be empty")
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

// ExtractFromURL извлекает ID турнира из URL
// Пример: /tournaments/pervenstvo-tsfo-18171615-let-16756891/ -> "16756891"
func ExtractFromURL(url string) (ID, error) {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return NewID(matches[1])
	}
	return "", fmt.Errorf("cannot extract id from url: %s", url)
}
