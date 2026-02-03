package entities

import "time"

// Match представляет матч (Domain Entity)
type Match struct {
	ID           string  `db:"id"`
	ExternalID   string  `db:"external_id"`
	TournamentID *string `db:"tournament_id"`

	HomeTeamID *string `db:"home_team_id"`
	AwayTeamID *string `db:"away_team_id"`

	HomeScore   *int `db:"home_score"`
	AwayScore   *int `db:"away_score"`
	HomeScoreP1 *int `db:"home_score_p1"`
	AwayScoreP1 *int `db:"away_score_p1"`
	HomeScoreP2 *int `db:"home_score_p2"`
	AwayScoreP2 *int `db:"away_score_p2"`
	HomeScoreP3 *int `db:"home_score_p3"`
	AwayScoreP3 *int `db:"away_score_p3"`
	HomeScoreOT *int `db:"home_score_ot"`
	AwayScoreOT *int `db:"away_score_ot"`

	MatchNumber  *int       `db:"match_number"`
	ScheduledAt  *time.Time `db:"scheduled_at"`
	Status       string     `db:"status"`
	ResultType   *string    `db:"result_type"`
	Venue        *string    `db:"venue"`
	GroupName    *string    `db:"group_name"`
	BirthYear    *int       `db:"birth_year"`
	VideoURL     *string    `db:"video_url"`

	Source        string    `db:"source"`
	Domain        *string   `db:"domain"`
	DetailsParsed bool      `db:"details_parsed"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// Match status constants
const (
	MatchStatusScheduled  = "scheduled"
	MatchStatusInProgress = "in_progress"
	MatchStatusFinished   = "finished"
	MatchStatusCancelled  = "cancelled"
)

// Result type constants
const (
	ResultTypeRegular  = "regular"
	ResultTypeOT       = "OT"
	ResultTypeShootout = "SO"
)

// IsFinished проверяет завершён ли матч
func (m *Match) IsFinished() bool {
	return m.Status == MatchStatusFinished
}

// NeedsDetailsParsing проверяет нужен ли парсинг деталей
func (m *Match) NeedsDetailsParsing() bool {
	return m.IsFinished() && !m.DetailsParsed
}

// IsScheduled проверяет запланирован ли матч
func (m *Match) IsScheduled() bool {
	return m.Status == MatchStatusScheduled
}

// GetTotalScore возвращает общий счёт матча
func (m *Match) GetTotalScore() (home, away int) {
	if m.HomeScore != nil {
		home = *m.HomeScore
	}
	if m.AwayScore != nil {
		away = *m.AwayScore
	}
	return
}
