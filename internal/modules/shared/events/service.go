package events

import (
	"context"
)

// EventStore интерфейс для хранения событий (локальная копия)
type EventStore interface {
	SaveEvent(ctx context.Context, event Event) error
	GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
}

// EventBus интерфейс для публикации/подписки (локальная копия)
type EventBus interface {
	Subscribe(eventType string, handler func(ctx context.Context, event Event) error)
	Publish(ctx context.Context, event Event) error
	PublishAsync(ctx context.Context, event Event)
}

// EventService координирует EventStore и EventBus
type EventService struct {
	store EventStore
	bus   EventBus
}

// NewEventService создает новый EventService
func NewEventService(store EventStore, bus EventBus) *EventService {
	return &EventService{
		store: store,
		bus:   bus,
	}
}

// SaveAndPublish сохраняет событие и публикует его
func (s *EventService) SaveAndPublish(ctx context.Context, event Event) error {
	// Сохраняем в store
	if err := s.store.SaveEvent(ctx, event); err != nil {
		return err
	}

	// Публикуем в bus
	return s.bus.Publish(ctx, event)
}

// SaveAndPublishAsync сохраняет событие и публикует его асинхронно
func (s *EventService) SaveAndPublishAsync(ctx context.Context, event Event) error {
	// Сохраняем в store
	if err := s.store.SaveEvent(ctx, event); err != nil {
		return err
	}

	// Публикуем асинхронно
	s.bus.PublishAsync(ctx, event)
	return nil
}

// GetEvents получает события для агрегата
func (s *EventService) GetEvents(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error) {
	return s.store.GetEvents(ctx, aggregateID, fromVersion)
}

// Subscribe подписывается на события
func (s *EventService) Subscribe(eventType string, handler func(ctx context.Context, event Event) error) {
	s.bus.Subscribe(eventType, handler)
}
