package entities

import "time"

// Tournament представляет турнир (Domain Entity)
type Tournament struct {
	ID                  string     `db:"id"`
	URL                 string     `db:"url"`
	Name                string     `db:"name"`
	Domain              string     `db:"domain"`
	Season              string     `db:"season"`
	StartDate           *time.Time `db:"start_date"`
	EndDate             *time.Time `db:"end_date"`
	IsEnded             bool       `db:"is_ended"`
	ExternalID          *string    `db:"external_id"`
	BirthYear           *int       `db:"birth_year"`
	GroupName           *string    `db:"group_name"`
	Region              *string    `db:"region"`
	Source              string     `db:"source"`
	CreatedAt           time.Time  `db:"created_at"`
	LastPlayersParsedAt *time.Time `db:"last_players_parsed_at"`
	LastStatsParsedAt   *time.Time `db:"last_stats_parsed_at"`
}

// ExtractIDFromURL извлекает ID турнира из URL (deprecated, use id.ExtractTournamentIDFromURL)
func ExtractTournamentIDFromURLLegacy(url string) string {
	id, _ := ExtractTournamentIDFromURL(url)
	return id.String()
}
