package store

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events"
)

// EventStore interface for event storage.
type EventStore interface {
	SaveEvent(ctx context.Context, event events.Event) error
	SaveEvents(ctx context.Context, events []events.Event) error
	GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]events.Event, error)
	GetAllEvents(ctx context.Context, aggregateType string, limit, offset int) ([]events.Event, error)
}

// PostgresEventStore implements EventStore for PostgreSQL.
type PostgresEventStore struct {
	db Database
}

// Database interface for DB operations.
type Database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
}

// Rows interface for query results.
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

// NewPostgresEventStore creates a new PostgreSQL EventStore.
func NewPostgresEventStore(db Database) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

// SaveEvent saves a single event.
func (s *PostgresEventStore) SaveEvent(ctx context.Context, event events.Event) error {
	query := `
		INSERT INTO player_events (player_id, event_type, event_data, source, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	return s.db.ExecContext(ctx, query,
		event.AggregateID(),
		event.EventType(),
		event.EventData(),
		"system",
		event.OccurredAt(),
	)
}

// SaveEvents saves multiple events.
func (s *PostgresEventStore) SaveEvents(ctx context.Context, evts []events.Event) error {
	for _, event := range evts {
		if err := s.SaveEvent(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// GetEvents retrieves events for an aggregate.
func (s *PostgresEventStore) GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]events.Event, error) {
	query := `
		SELECT player_id, event_type, event_data, created_at
		FROM player_events 
		WHERE player_id = $1 
		ORDER BY created_at ASC`

	rows, err := s.db.QueryContext(ctx, query, aggregateID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []events.Event
	for rows.Next() {
		var playerID, eventType string
		var eventData interface{}
		var createdAt interface{}

		if err := rows.Scan(&playerID, &eventType, &eventData, &createdAt); err != nil {
			return nil, err
		}

		event := events.NewBaseEvent(eventType, playerID, "player", eventData, 1)
		result = append(result, event)
	}

	return result, nil
}

// GetAllEvents retrieves all events of a type.
func (s *PostgresEventStore) GetAllEvents(ctx context.Context, aggregateType string, limit, offset int) ([]events.Event, error) {
	query := `
		SELECT player_id, event_type, event_data, created_at
		FROM player_events 
		ORDER BY created_at ASC 
		LIMIT $1 OFFSET $2`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []events.Event
	for rows.Next() {
		var playerID, eventType string
		var eventData interface{}
		var createdAt interface{}

		if err := rows.Scan(&playerID, &eventType, &eventData, &createdAt); err != nil {
			return nil, err
		}

		event := events.NewBaseEvent(eventType, playerID, "player", eventData, 1)
		result = append(result, event)
	}

	return result, nil
}
