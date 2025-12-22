package entities

import (
	"crypto/sha256"
	"fmt"
	"regexp"
)

// TeamID представляет уникальный идентификатор команды
type TeamID string

// NewTeamID создает новый ID из строки
func NewTeamID(id string) (TeamID, error) {
	if id == "" {
		return "", fmt.Errorf("team id cannot be empty")
	}
	return TeamID(id), nil
}

// String возвращает строковое представление
func (id TeamID) String() string {
	return string(id)
}

// IsEmpty проверяет пустой ли ID
func (id TeamID) IsEmpty() bool {
	return id == ""
}

// ExtractTeamIDFromURL извлекает ID из URL команды
func ExtractTeamIDFromURL(url string) TeamID {
	re := regexp.MustCompile(`_(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return TeamID(matches[1])
	}
	return TeamID(generateHashID(url))
}

// ExtractTeamIDFromURLLegacy возвращает string (совместимость)
func ExtractTeamIDFromURLLegacy(url string) string {
	return string(ExtractTeamIDFromURL(url))
}

func generateHashID(url string) string {
	hash := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x", hash)[:8]
}
