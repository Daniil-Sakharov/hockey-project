package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events"
	"github.com/jmoiron/sqlx"
)

// SQLXEventStore реализация EventStore для PostgreSQL с sqlx
type SQLXEventStore struct {
	db *sqlx.DB
}

// NewSQLXEventStore создает новый EventStore с sqlx
func NewSQLXEventStore(db *sqlx.DB) *SQLXEventStore {
	return &SQLXEventStore{db: db}
}

// SaveEvent сохраняет событие в БД
func (s *SQLXEventStore) SaveEvent(ctx context.Context, event events.Event) error {
	eventData, err := json.Marshal(event.EventData())
	if err != nil {
		return err
	}

	query := `
		INSERT INTO player_events (player_id, event_type, event_data, source, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = s.db.ExecContext(ctx, query,
		event.AggregateID(),
		event.EventType(),
		eventData,
		"system",
		event.OccurredAt(),
	)
	return err
}

// SaveEvents сохраняет несколько событий
func (s *SQLXEventStore) SaveEvents(ctx context.Context, evts []events.Event) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO player_events (player_id, event_type, event_data, source, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	for _, event := range evts {
		eventData, err := json.Marshal(event.EventData())
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query,
			event.AggregateID(),
			event.EventType(),
			eventData,
			"system",
			event.OccurredAt(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// eventRow структура для чтения из БД
type eventRow struct {
	PlayerID  string    `db:"player_id"`
	EventType string    `db:"event_type"`
	EventData []byte    `db:"event_data"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}

// GetEvents получает события для агрегата
func (s *SQLXEventStore) GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]events.Event, error) {
	query := `
		SELECT player_id, event_type, event_data, source, created_at
		FROM player_events 
		WHERE player_id = $1 
		ORDER BY created_at ASC`

	var rows []eventRow
	if err := s.db.SelectContext(ctx, &rows, query, aggregateID); err != nil {
		return nil, err
	}

	result := make([]events.Event, 0, len(rows))
	for _, row := range rows {
		var data interface{}
		if err := json.Unmarshal(row.EventData, &data); err != nil {
			data = string(row.EventData)
		}
		event := events.NewBaseEvent(row.EventType, row.PlayerID, "player", data, 1)
		result = append(result, event)
	}

	return result, nil
}

// GetAllEvents получает все события определенного типа
func (s *SQLXEventStore) GetAllEvents(ctx context.Context, aggregateType string, limit, offset int) ([]events.Event, error) {
	query := `
		SELECT player_id, event_type, event_data, source, created_at
		FROM player_events 
		ORDER BY created_at ASC 
		LIMIT $1 OFFSET $2`

	var rows []eventRow
	if err := s.db.SelectContext(ctx, &rows, query, limit, offset); err != nil {
		return nil, err
	}

	result := make([]events.Event, 0, len(rows))
	for _, row := range rows {
		var data interface{}
		if err := json.Unmarshal(row.EventData, &data); err != nil {
			data = string(row.EventData)
		}
		event := events.NewBaseEvent(row.EventType, row.PlayerID, "player", data, 1)
		result = append(result, event)
	}

	return result, nil
}
