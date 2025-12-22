package entities

import "time"

// GoalieStatistic представляет статистику вратаря в турнире
type GoalieStatistic struct {
	ID              int       `db:"id"`
	PlayerID        string    `db:"player_id"`
	TeamID          string    `db:"team_id"`
	TournamentID    string    `db:"tournament_id"`
	Games           int       `db:"games"`
	Minutes         int       `db:"minutes"`
	GoalsAgainst    int       `db:"goals_against"`
	ShotsAgainst    int       `db:"shots_against"`
	SavePercentage  *float64  `db:"save_percentage"`
	GoalsAgainstAvg *float64  `db:"goals_against_avg"`
	Wins            int       `db:"wins"`
	Shutouts        int       `db:"shutouts"`
	Assists         int       `db:"assists"`
	PenaltyMinutes  int       `db:"penalty_minutes"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
