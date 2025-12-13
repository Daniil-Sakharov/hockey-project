package tournament

import (
	"regexp"
	"time"
)

// Tournament представляет турнир
type Tournament struct {
	ID        string     `db:"id"`         // ID извлеченный из URL
	URL       string     `db:"url"`        // /tournaments/pervenstvo-tsfo-18171615-let-16756891/
	Name      string     `db:"name"`       // Первенство ЦФО 18/17/16/15 лет
	Domain    string     `db:"domain"`     // https://cfo.fhr.ru
	Season    string     `db:"season"`     // 2025/2026
	StartDate *time.Time `db:"start_date"` // Дата начала турнира
	EndDate   *time.Time `db:"end_date"`   // Дата окончания турнира (NULL если активен)
	IsEnded   bool       `db:"is_ended"`   // Флаг завершенности турнира
	CreatedAt time.Time  `db:"created_at"` // Дата создания
}

// ExtractIDFromURL извлекает ID турнира из URL
// Пример: /tournaments/pervenstvo-tsfo-18171615-let-16756891/ -> "16756891"
func ExtractIDFromURL(url string) string {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
