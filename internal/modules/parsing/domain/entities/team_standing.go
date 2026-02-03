package entities

import "time"

// TeamStanding представляет позицию команды в турнирной таблице
type TeamStanding struct {
	ID           string `db:"id"`
	TournamentID string `db:"tournament_id"`
	TeamID       string `db:"team_id"`

	// Позиция и очки
	Position *int `db:"position"`
	Points   int  `db:"points"`

	// Матчи
	Games    int `db:"games"`
	Wins     int `db:"wins"`
	WinsOT   int `db:"wins_ot"`
	WinsSO   int `db:"wins_so"`
	LossesSO int `db:"losses_so"`
	LossesOT int `db:"losses_ot"`
	Losses   int `db:"losses"`
	Draws    int `db:"draws"`

	// Шайбы
	GoalsFor       int `db:"goals_for"`
	GoalsAgainst   int `db:"goals_against"`
	GoalDifference int `db:"goal_difference"`

	// Метаданные
	GroupName *string   `db:"group_name"`
	BirthYear *int      `db:"birth_year"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// CalculateGoalDifference вычисляет разницу шайб
func (s *TeamStanding) CalculateGoalDifference() {
	s.GoalDifference = s.GoalsFor - s.GoalsAgainst
}

// GetTotalWins возвращает общее количество побед
func (s *TeamStanding) GetTotalWins() int {
	return s.Wins + s.WinsOT + s.WinsSO
}

// GetTotalLosses возвращает общее количество поражений
func (s *TeamStanding) GetTotalLosses() int {
	return s.Losses + s.LossesOT + s.LossesSO
}

// GetWinPercentage возвращает процент побед
func (s *TeamStanding) GetWinPercentage() float64 {
	if s.Games == 0 {
		return 0
	}
	return float64(s.GetTotalWins()) / float64(s.Games) * 100
}
