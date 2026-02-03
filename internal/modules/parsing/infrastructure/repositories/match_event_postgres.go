package repositories

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// StringSlice - обёртка для []string для сохранения в JSONB
type StringSlice []string

// Value implements driver.Valuer
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// Scan implements sql.Scanner
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, s)
}

type MatchEventPostgres struct {
	db *sqlx.DB
}

func NewMatchEventPostgres(db *sqlx.DB) *MatchEventPostgres {
	return &MatchEventPostgres{db: db}
}

func (r *MatchEventPostgres) Create(ctx context.Context, e *entities.MatchEvent) error {
	// Генерируем UUID если ID не задан
	if e.ID == "" {
		e.ID = uuid.New().String()
	}

	query := `
		INSERT INTO match_events (id, match_id, event_type, period, time_minutes, time_seconds,
			scorer_player_id, assist1_player_id, assist2_player_id, team_id, goal_type,
			goalie_player_id, home_players_on_ice, away_players_on_ice,
			penalty_player_id, penalty_minutes, penalty_reason, penalty_reason_code, is_home,
			score_home, score_away, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, NOW())`

	// Конвертируем слайсы в JSONB
	var homeOnIce, awayOnIce StringSlice
	if len(e.HomePlayersOnIce) > 0 {
		homeOnIce = StringSlice(e.HomePlayersOnIce)
	}
	if len(e.AwayPlayersOnIce) > 0 {
		awayOnIce = StringSlice(e.AwayPlayersOnIce)
	}

	_, err := r.db.ExecContext(ctx, query,
		e.ID, e.MatchID, e.EventType, e.Period, e.TimeMinutes, e.TimeSeconds,
		e.ScorerPlayerID, e.Assist1PlayerID, e.Assist2PlayerID, e.TeamID, e.GoalType,
		e.GoaliePlayerID, homeOnIce, awayOnIce,
		e.PenaltyPlayerID, e.PenaltyMinutes, e.PenaltyReason, e.PenaltyReasonCode, e.IsHome,
		e.ScoreHome, e.ScoreAway, e.Source,
	)
	if err != nil {
		return fmt.Errorf("create match event: %w", err)
	}
	return nil
}

func (r *MatchEventPostgres) CreateBatch(ctx context.Context, events []*entities.MatchEvent) error {
	if len(events) == 0 {
		return nil
	}

	query := `
		INSERT INTO match_events (id, match_id, event_type, period, time_minutes, time_seconds,
			scorer_player_id, assist1_player_id, assist2_player_id, team_id, goal_type,
			penalty_player_id, penalty_minutes, penalty_reason, is_home, source, created_at)
		VALUES (:id, :match_id, :event_type, :period, :time_minutes, :time_seconds,
			:scorer_player_id, :assist1_player_id, :assist2_player_id, :team_id, :goal_type,
			:penalty_player_id, :penalty_minutes, :penalty_reason, :is_home, :source, NOW())
		ON CONFLICT (id) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, query, events)
	if err != nil {
		return fmt.Errorf("create match events batch: %w", err)
	}
	return nil
}

func (r *MatchEventPostgres) GetByMatchID(ctx context.Context, matchID string) ([]*entities.MatchEvent, error) {
	var events []*entities.MatchEvent
	err := r.db.SelectContext(ctx, &events,
		`SELECT * FROM match_events WHERE match_id = $1 ORDER BY period, time_minutes, time_seconds`, matchID)
	if err != nil {
		return nil, fmt.Errorf("get events by match: %w", err)
	}
	return events, nil
}

func (r *MatchEventPostgres) GetGoalsByMatchID(ctx context.Context, matchID string) ([]*entities.MatchEvent, error) {
	var events []*entities.MatchEvent
	err := r.db.SelectContext(ctx, &events,
		`SELECT * FROM match_events WHERE match_id = $1 AND event_type = 'goal'
		ORDER BY period, time_minutes, time_seconds`, matchID)
	if err != nil {
		return nil, fmt.Errorf("get goals by match: %w", err)
	}
	return events, nil
}

func (r *MatchEventPostgres) GetPenaltiesByMatchID(ctx context.Context, matchID string) ([]*entities.MatchEvent, error) {
	var events []*entities.MatchEvent
	err := r.db.SelectContext(ctx, &events,
		`SELECT * FROM match_events WHERE match_id = $1 AND event_type = 'penalty'
		ORDER BY period, time_minutes, time_seconds`, matchID)
	if err != nil {
		return nil, fmt.Errorf("get penalties by match: %w", err)
	}
	return events, nil
}

func (r *MatchEventPostgres) GetByPlayerID(ctx context.Context, playerID string) ([]*entities.MatchEvent, error) {
	var events []*entities.MatchEvent
	err := r.db.SelectContext(ctx, &events,
		`SELECT * FROM match_events
		WHERE scorer_player_id = $1 OR assist1_player_id = $1 OR assist2_player_id = $1 OR penalty_player_id = $1
		ORDER BY created_at DESC`, playerID)
	if err != nil {
		return nil, fmt.Errorf("get events by player: %w", err)
	}
	return events, nil
}

func (r *MatchEventPostgres) DeleteByMatchID(ctx context.Context, matchID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM match_events WHERE match_id = $1`, matchID)
	if err != nil {
		return fmt.Errorf("delete events by match: %w", err)
	}
	return nil
}
