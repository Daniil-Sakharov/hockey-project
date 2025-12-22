package domain

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events"
)

// PlayerCreated событие создания игрока
type PlayerCreated struct {
	*events.BaseEvent
	PlayerData PlayerData `json:"player_data"`
}

// PlayerUpdated событие обновления игрока
type PlayerUpdated struct {
	*events.BaseEvent
	OldData PlayerData `json:"old_data"`
	NewData PlayerData `json:"new_data"`
}

// PlayerData данные игрока
type PlayerData struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Source   string `json:"source"`
}

// NewPlayerCreated создает событие создания игрока
func NewPlayerCreated(playerID string, data PlayerData) *PlayerCreated {
	baseEvent := events.NewBaseEvent(
		"player.created",
		playerID,
		"player",
		data,
		1,
	)

	return &PlayerCreated{
		BaseEvent:  baseEvent,
		PlayerData: data,
	}
}

// NewPlayerUpdated создает событие обновления игрока
func NewPlayerUpdated(playerID string, oldData, newData PlayerData) *PlayerUpdated {
	baseEvent := events.NewBaseEvent(
		"player.updated",
		playerID,
		"player",
		map[string]interface{}{
			"old_data": oldData,
			"new_data": newData,
		},
		1,
	)

	return &PlayerUpdated{
		BaseEvent: baseEvent,
		OldData:   oldData,
		NewData:   newData,
	}
}
