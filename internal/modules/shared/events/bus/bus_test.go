package bus

import (
	"context"
	"testing"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events"
)

func TestInMemoryEventBus_PublishSubscribe(t *testing.T) {
	bus := NewInMemoryEventBus()

	received := make(chan events.Event, 1)

	// Подписываемся на события
	bus.Subscribe("test.event", func(ctx context.Context, event events.Event) error {
		received <- event
		return nil
	})

	// Создаем и публикуем событие
	event := events.NewBaseEvent("test.event", "test-id", "test", "test-data", 1)

	err := bus.Publish(context.Background(), event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Проверяем что событие получено
	select {
	case receivedEvent := <-received:
		if receivedEvent.EventType() != "test.event" {
			t.Errorf("Expected 'test.event', got %s", receivedEvent.EventType())
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Event not received within timeout")
	}
}

func TestInMemoryEventBus_MultipleSubscribers(t *testing.T) {
	bus := NewInMemoryEventBus()

	received1 := make(chan bool, 1)
	received2 := make(chan bool, 1)

	// Два подписчика на одно событие
	bus.Subscribe("multi.event", func(ctx context.Context, event events.Event) error {
		received1 <- true
		return nil
	})

	bus.Subscribe("multi.event", func(ctx context.Context, event events.Event) error {
		received2 <- true
		return nil
	})

	// Публикуем событие
	event := events.NewBaseEvent("multi.event", "test-id", "test", "data", 1)
	err := bus.Publish(context.Background(), event)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Проверяем что оба подписчика получили событие
	timeout := time.After(100 * time.Millisecond)

	select {
	case <-received1:
	case <-timeout:
		t.Error("First subscriber did not receive event")
	}

	select {
	case <-received2:
	case <-timeout:
		t.Error("Second subscriber did not receive event")
	}
}
