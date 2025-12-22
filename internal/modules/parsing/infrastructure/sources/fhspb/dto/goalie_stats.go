package dto

type GoalieStatsDTO struct {
	PlayerID        string
	TeamID          string
	PlayerName      string
	Number          int
	BirthDate       string
	TeamName        string
	Games           int
	Minutes         int
	GoalsAgainst    int
	ShotsAgainst    int
	SavePercentage  float64
	GoalsAgainstAvg float64
	Wins            int
	Shutouts        int
	Assists         int
	PenaltyMinutes  int
}
