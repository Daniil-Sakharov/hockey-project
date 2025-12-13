package bot

import (
	"context"

	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// StateManager управляет состоянием пользователей
type StateManager interface {
	GetState(userID int64) *domainBot.UserState
	UpdateFilters(userID int64, filters domainBot.SearchFilters)
	SetLastMsgID(userID int64, msgID int)
	SetCurrentView(userID int64, view string)
	SetWaitingForInput(userID int64, input string)
	ResetFilters(userID int64)
	ClearState(userID int64)
}

// ProfileService сервис для работы с профилями игроков
// Реализация: internal/service/bot/profile/service.go
type ProfileService interface {
	// GetPlayerProfile получает полный профиль игрока для отображения
	GetPlayerProfile(ctx context.Context, playerID string) (*domainPlayer.Profile, error)
}
