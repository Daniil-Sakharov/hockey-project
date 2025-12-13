package team

import "time"

// Team представляет команду (Domain Entity)
type Team struct {
	ID        string    `db:"id"`
	URL       string    `db:"url"`
	Name      string    `db:"name"`
	City      string    `db:"city"`
	CreatedAt time.Time `db:"created_at"`
}

// ExtractIDFromURL извлекает ID из URL команды (deprecated, use id.ExtractFromURL)
func ExtractIDFromURL(url string) string {
	return ExtractFromURL(url).String()
}
