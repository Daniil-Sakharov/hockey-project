package events

import (
	"testing"
	"time"
)

func TestBaseEvent(t *testing.T) {
	data := map[string]interface{}{
		"name": "Test Player",
		"age":  25,
	}

	event := NewBaseEvent("player.created", "player-123", "player", data, 1)

	if event.EventID() == "" {
		t.Error("Expected event ID, got empty string")
	}

	if event.EventType() != "player.created" {
		t.Errorf("Expected event type 'player.created', got %s", event.EventType())
	}

	if event.AggregateID() != "player-123" {
		t.Errorf("Expected aggregate ID 'player-123', got %s", event.AggregateID())
	}

	if event.AggregateType() != "player" {
		t.Errorf("Expected aggregate type 'player', got %s", event.AggregateType())
	}

	if event.Version() != 1 {
		t.Errorf("Expected version 1, got %d", event.Version())
	}

	if time.Since(event.OccurredAt()) > time.Second {
		t.Error("Event occurred at time is too old")
	}
}
