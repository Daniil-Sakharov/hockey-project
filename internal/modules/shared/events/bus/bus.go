package bus

import (
	"context"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events"
)

// EventHandler функция обработки события
type EventHandler func(ctx context.Context, event events.Event) error

// EventBus интерфейс для публикации/подписки на события
type EventBus interface {
	Subscribe(eventType string, handler EventHandler)
	Publish(ctx context.Context, event events.Event) error
	PublishAsync(ctx context.Context, event events.Event)
}

// InMemoryEventBus реализация EventBus в памяти
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventBus создает новый EventBus
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe подписывается на события типа
func (b *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish публикует событие синхронно
func (b *InMemoryEventBus) Publish(ctx context.Context, event events.Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.EventType()]
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// PublishAsync публикует событие асинхронно
func (b *InMemoryEventBus) PublishAsync(ctx context.Context, event events.Event) {
	go func() {
		_ = b.Publish(ctx, event)
	}()
}
