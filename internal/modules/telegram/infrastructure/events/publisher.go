package events

import (
	"context"
	"strconv"

	sharedEvents "github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events/bus"
)

const (
	EventPlayerSearchRequested = "telegram.player_search_requested"
	EventFilterApplied         = "telegram.filter_applied"
	EventSessionStarted        = "telegram.session_started"
)

// Publisher publishes telegram events to the event bus.
type Publisher struct {
	bus bus.EventBus
}

// NewPublisher creates a new event publisher.
func NewPublisher(eventBus bus.EventBus) *Publisher {
	return &Publisher{bus: eventBus}
}

// PublishSearchRequested publishes player search event.
func (p *Publisher) PublishSearchRequested(ctx context.Context, userID int64, query string) {
	event := sharedEvents.NewBaseEvent(
		EventPlayerSearchRequested,
		strconv.FormatInt(userID, 10),
		"user",
		map[string]string{"query": query},
		1,
	)
	p.bus.PublishAsync(ctx, event)
}

// PublishFilterApplied publishes filter applied event.
func (p *Publisher) PublishFilterApplied(ctx context.Context, userID int64, filterType, value string) {
	event := sharedEvents.NewBaseEvent(
		EventFilterApplied,
		strconv.FormatInt(userID, 10),
		"user",
		map[string]string{"filter_type": filterType, "value": value},
		1,
	)
	p.bus.PublishAsync(ctx, event)
}

// PublishSessionStarted publishes session started event.
func (p *Publisher) PublishSessionStarted(ctx context.Context, userID int64) {
	event := sharedEvents.NewBaseEvent(
		EventSessionStarted,
		strconv.FormatInt(userID, 10),
		"user",
		nil,
		1,
	)
	p.bus.PublishAsync(ctx, event)
}
