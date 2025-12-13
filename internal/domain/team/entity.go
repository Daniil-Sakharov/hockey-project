package team

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"time"
)

// Team представляет команду
type Team struct {
	ID        string    `db:"id"`         // ID (hash от URL или извлеченный из URL)
	URL       string    `db:"url"`        // /tournaments/.../buran_5136295/
	Name      string    `db:"name"`       // БУРАН
	City      string    `db:"city"`       // Воронеж
	CreatedAt time.Time `db:"created_at"` // Дата создания
}

// ExtractIDFromURL извлекает ID из URL команды
// Пример: /tournaments/pervenstvo/.../buran_5136295/ -> "5136295"
func ExtractIDFromURL(url string) string {
	re := regexp.MustCompile(`_(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	// Fallback: используем hash от URL
	return generateHashID(url)
}

// generateHashID генерирует MD5 hash от URL (первые 8 символов)
func generateHashID(url string) string {
	hash := md5.Sum([]byte(url))
	return fmt.Sprintf("%x", hash)[:8]
}
