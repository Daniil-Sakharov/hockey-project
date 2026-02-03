package entities

import "time"

// MatchTeamStats представляет статистику команды за матч (Domain Entity)
type MatchTeamStats struct {
	ID      string `db:"id"`
	MatchID string `db:"match_id"`
	TeamID  string `db:"team_id"`

	// Броски по периодам
	ShotsP1    int `db:"shots_p1"`
	ShotsP2    int `db:"shots_p2"`
	ShotsP3    int `db:"shots_p3"`
	ShotsOT    int `db:"shots_ot"`
	ShotsTotal int `db:"shots_total"`

	// Метаданные
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}

// CalculateTotal вычисляет общее количество бросков
func (s *MatchTeamStats) CalculateTotal() {
	s.ShotsTotal = s.ShotsP1 + s.ShotsP2 + s.ShotsP3 + s.ShotsOT
}

// GetShotsByPeriod возвращает броски для указанного периода (1-4)
func (s *MatchTeamStats) GetShotsByPeriod(period int) int {
	switch period {
	case 1:
		return s.ShotsP1
	case 2:
		return s.ShotsP2
	case 3:
		return s.ShotsP3
	case 4:
		return s.ShotsOT
	default:
		return 0
	}
}
