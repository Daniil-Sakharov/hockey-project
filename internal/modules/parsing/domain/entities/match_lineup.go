package entities

import "time"

// MatchLineup представляет игрока в составе матча
type MatchLineup struct {
	ID       string `db:"id"`
	MatchID  string `db:"match_id"`
	PlayerID string `db:"player_id"`
	TeamID   string `db:"team_id"`

	JerseyNumber *int    `db:"jersey_number"`
	Position     *string `db:"position"`
	CaptainRole  *string `db:"captain_role"` // К = капитан, А = ассистент

	// Статистика за матч
	Goals          int `db:"goals"`
	Assists        int `db:"assists"`
	PenaltyMinutes int `db:"penalty_minutes"`
	PlusMinus      int `db:"plus_minus"`

	// Для вратарей
	Saves       *int `db:"saves"`
	GoalsAgainst *int `db:"goals_against"`
	TimeOnIce   *int `db:"time_on_ice"`

	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}

// Lineup position constants
const (
	LineupPositionGoalie   = "G"
	LineupPositionDefender = "D"
	LineupPositionForward  = "F"
)

// IsGoalie проверяет является ли игрок вратарём
func (l *MatchLineup) IsGoalie() bool {
	return l.Position != nil && *l.Position == LineupPositionGoalie
}

// GetPoints возвращает количество очков (голы + передачи)
func (l *MatchLineup) GetPoints() int {
	return l.Goals + l.Assists
}

// GetTimeOnIceMinutes возвращает время на льду в минутах
func (l *MatchLineup) GetTimeOnIceMinutes() int {
	if l.TimeOnIce == nil {
		return 0
	}
	return *l.TimeOnIce / 60
}
