package entities

import (
	"fmt"
	"regexp"
)

// ID представляет уникальный идентификатор турнира (Value Object)
type TournamentID string

// NewID создает новый ID из строки
func NewTournamentID(id string) (TournamentID, error) {
	if id == "" {
		return "", fmt.Errorf("tournament id cannot be empty")
	}
	return TournamentID(id), nil
}

// String возвращает строковое представление
func (id TournamentID) String() string {
	return string(id)
}

// IsEmpty проверяет пустой ли ID
func (id TournamentID) IsEmpty() bool {
	return id == ""
}

// ExtractFromURL извлекает ID турнира из URL
// Пример: /tournaments/pervenstvo-tsfo-18171615-let-16756891/ -> "16756891"
func ExtractTournamentIDFromURL(url string) (TournamentID, error) {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return NewTournamentID(matches[1])
	}
	return "", fmt.Errorf("cannot extract id from url: %s", url)
}
