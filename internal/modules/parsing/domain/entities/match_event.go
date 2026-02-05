package entities

import "time"

// MatchEvent представляет событие матча (гол, штраф)
type MatchEvent struct {
	ID      string `db:"id"`
	MatchID string `db:"match_id"`

	EventType   string `db:"event_type"`
	Period      *int   `db:"period"`
	TimeMinutes *int   `db:"time_minutes"`
	TimeSeconds *int   `db:"time_seconds"`

	// Для голов
	ScorerPlayerID   *string  `db:"scorer_player_id"`
	Assist1PlayerID  *string  `db:"assist1_player_id"`
	Assist2PlayerID  *string  `db:"assist2_player_id"`
	TeamID           *string  `db:"team_id"`
	GoalType         *string  `db:"goal_type"`
	GoaliePlayerID   *string  `db:"goalie_player_id"`    // Вратарь (NULL если пустые ворота)
	HomePlayersOnIce []string `db:"home_players_on_ice"` // Игроки дома на льду
	AwayPlayersOnIce []string `db:"away_players_on_ice"` // Игроки гостей на льду
	ScoreHome        *int     `db:"score_home"`          // Счёт после гола (домашние)
	ScoreAway        *int     `db:"score_away"`          // Счёт после гола (гости)

	// Для штрафов
	PenaltyPlayerID   *string `db:"penalty_player_id"`
	PenaltyMinutes    *int    `db:"penalty_minutes"`
	PenaltyReason     *string `db:"penalty_reason"`
	PenaltyReasonCode *string `db:"penalty_reason_code"`

	// Для событий команды (вратарь, пустые ворота, тайм-аут)
	IsHome *bool `db:"is_home"` // true = домашняя команда, false = гости

	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}

// Event type constants
const (
	EventTypeGoal         = "goal"
	EventTypePenalty      = "penalty"
	EventTypeGoalieChange = "goalie_change" // Смена вратаря
	EventTypeEmptyNet     = "empty_net"     // Пустые ворота
	EventTypeTimeout      = "timeout"       // Тайм-аут
)

// Goal type constants
const (
	GoalTypeEven      = "even"
	GoalTypePowerPlay = "pp"
	GoalTypeShortHand = "sh"
	GoalTypeEmptyNet  = "en"
)

// IsGoal проверяет является ли событие голом
func (e *MatchEvent) IsGoal() bool {
	return e.EventType == EventTypeGoal
}

// IsPenalty проверяет является ли событие штрафом
func (e *MatchEvent) IsPenalty() bool {
	return e.EventType == EventTypePenalty
}

// GetTimeString возвращает строковое представление времени события
func (e *MatchEvent) GetTimeString() string {
	if e.TimeMinutes == nil {
		return ""
	}
	minutes := *e.TimeMinutes
	seconds := 0
	if e.TimeSeconds != nil {
		seconds = *e.TimeSeconds
	}
	return formatTime(minutes, seconds)
}

func formatTime(minutes, seconds int) string {
	if seconds < 10 {
		return string(rune('0'+minutes/10)) + string(rune('0'+minutes%10)) + ":0" + string(rune('0'+seconds))
	}
	return string(rune('0'+minutes/10)) + string(rune('0'+minutes%10)) + ":" + string(rune('0'+seconds/10)) + string(rune('0'+seconds%10))
}
