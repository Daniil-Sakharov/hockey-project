package entities

import "time"

// PlayerTeam представляет связь игрок-команда-турнир
type PlayerTeam struct {
	PlayerID     string     `db:"player_id"`
	TeamID       string     `db:"team_id"`
	TournamentID string     `db:"tournament_id"`
	Number       *string    `db:"number"`
	Season       string     `db:"season"`
	StartedAt    *time.Time `db:"started_at"`
	EndedAt      *time.Time `db:"ended_at"`
	IsActive     bool       `db:"is_active"`
	JerseyNumber *int       `db:"jersey_number"`
	Role         *string    `db:"role"`
	Position     *string    `db:"position"`
	Height       *int       `db:"height"`
	Weight       *int       `db:"weight"`
	PhotoURL     *string    `db:"photo_url"`
	BirthYear    *int       `db:"birth_year"`
	GroupName    *string    `db:"group_name"`
	Source       string     `db:"source"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
