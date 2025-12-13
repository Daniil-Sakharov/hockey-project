package tournament

import "time"

// Tournament представляет турнир (Domain Entity)
type Tournament struct {
	ID        string     `db:"id"`
	URL       string     `db:"url"`
	Name      string     `db:"name"`
	Domain    string     `db:"domain"`
	Season    string     `db:"season"`
	StartDate *time.Time `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
	IsEnded   bool       `db:"is_ended"`
	CreatedAt time.Time  `db:"created_at"`
}

// ExtractIDFromURL извлекает ID турнира из URL (deprecated, use id.ExtractFromURL)
func ExtractIDFromURL(url string) string {
	id, _ := ExtractFromURL(url)
	return id.String()
}
