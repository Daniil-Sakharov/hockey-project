package dto

type PlayerStatsDTO struct {
	PlayerID       string
	TeamID         string
	PlayerName     string
	Number         int
	Role           string // К, А или пусто
	BirthDate      string
	TeamName       string
	Position       string // Нп, Зщ, Вр
	Games          int
	Points         int
	PointsAvg      float64
	Goals          int
	Assists        int
	PlusMinus      int
	PenaltyMinutes int
	PenaltyAvg     float64
}
