package entities

import "time"

// Team представляет команду (Domain Entity)
type Team struct {
	ID           string    `db:"id"`
	URL          string    `db:"url"`
	Name         string    `db:"name"`
	City         string    `db:"city"`
	ExternalID   *string   `db:"external_id"`
	TournamentID *string   `db:"tournament_id"`
	Region       *string   `db:"region"`
	LogoURL      *string   `db:"logo_url"`
	Source       string    `db:"source"`
	CreatedAt    time.Time `db:"created_at"`
}
