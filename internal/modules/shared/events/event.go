package events

import (
	"time"

	"github.com/google/uuid"
)

// Event базовый интерфейс для всех событий
type Event interface {
	EventID() string
	EventType() string
	AggregateID() string
	AggregateType() string
	EventData() interface{}
	OccurredAt() time.Time
	Version() int
}

// BaseEvent базовая реализация Event
type BaseEvent struct {
	ID       string      `json:"event_id"`
	Type     string      `json:"event_type"`
	AggID    string      `json:"aggregate_id"`
	AggType  string      `json:"aggregate_type"`
	Data     interface{} `json:"event_data"`
	Occurred time.Time   `json:"occurred_at"`
	Ver      int         `json:"version"`
}

// NewBaseEvent создает новое базовое событие
func NewBaseEvent(eventType, aggregateID, aggregateType string, data interface{}, version int) *BaseEvent {
	return &BaseEvent{
		ID:       uuid.New().String(),
		Type:     eventType,
		AggID:    aggregateID,
		AggType:  aggregateType,
		Data:     data,
		Occurred: time.Now(),
		Ver:      version,
	}
}

func (e *BaseEvent) EventID() string        { return e.ID }
func (e *BaseEvent) EventType() string      { return e.Type }
func (e *BaseEvent) AggregateID() string    { return e.AggID }
func (e *BaseEvent) AggregateType() string  { return e.AggType }
func (e *BaseEvent) EventData() interface{} { return e.Data }
func (e *BaseEvent) OccurredAt() time.Time  { return e.Occurred }
func (e *BaseEvent) Version() int           { return e.Ver }
