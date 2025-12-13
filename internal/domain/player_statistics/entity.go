package player_statistics

import (
	"fmt"
	"time"
)

// PlayerStatistic представляет статистику игрока в турнире
type PlayerStatistic struct {
	ID           int       `db:"id"`
	TournamentID string    `db:"tournament_id"`
	PlayerID     string    `db:"player_id"`
	TeamID       string    `db:"team_id"`
	GroupName    string    `db:"group_name"`
	BirthYear    int       `db:"birth_year"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`

	// Основная статистика
	Games          int `db:"games"`
	Goals          int `db:"goals"`
	Assists        int `db:"assists"`
	Points         int `db:"points"`
	Plus           int `db:"plus"`
	Minus          int `db:"minus"`
	PlusMinus      int `db:"plus_minus"`
	PenaltyMinutes int `db:"penalty_minutes"`

	// Детальная статистика голов
	GoalsEvenStrength int `db:"goals_even_strength"`
	GoalsPowerPlay    int `db:"goals_power_play"`
	GoalsShortHanded  int `db:"goals_short_handed"`
	GoalsPeriod1      int `db:"goals_period_1"`
	GoalsPeriod2      int `db:"goals_period_2"`
	GoalsPeriod3      int `db:"goals_period_3"`
	GoalsOvertime     int `db:"goals_overtime"`
	HatTricks         int `db:"hat_tricks"`
	GameWinningGoals  int `db:"game_winning_goals"`

	// Средние показатели
	GoalsPerGame          float64 `db:"goals_per_game"`
	PointsPerGame         float64 `db:"points_per_game"`
	PenaltyMinutesPerGame float64 `db:"penalty_minutes_per_game"`
}

// Validate проверяет валидность данных статистики
func (ps *PlayerStatistic) Validate() error {
	if ps.TournamentID == "" {
		return fmt.Errorf("tournament_id is required")
	}
	if ps.PlayerID == "" {
		return fmt.Errorf("player_id is required")
	}
	if ps.TeamID == "" {
		return fmt.Errorf("team_id is required")
	}
	if ps.GroupName == "" {
		return fmt.Errorf("group_name is required")
	}
	if ps.BirthYear != 0 && (ps.BirthYear < 2000 || ps.BirthYear > 2020) {
		return fmt.Errorf("birth_year must be 0 or between 2000 and 2020")
	}
	if ps.Games < 0 || ps.Goals < 0 || ps.Assists < 0 {
		return fmt.Errorf("stats cannot be negative")
	}
	return nil
}

// IsValid возвращает true если статистика валидна
func (ps *PlayerStatistic) IsValid() bool {
	return ps.Validate() == nil
}
