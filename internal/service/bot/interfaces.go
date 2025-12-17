package bot

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// StateManager управляет состоянием пользователей
type StateManager interface {
	GetState(userID int64) *bot.UserState
	UpdateFilters(userID int64, filters bot.SearchFilters)
	SetLastMsgID(userID int64, msgID int)
	SetCurrentView(userID int64, view string)
	SetWaitingForInput(userID int64, input string)
	ResetFilters(userID int64)
	ClearState(userID int64)
}

// ProfileService сервис для работы с профилями игроков
type ProfileService interface {
	GetPlayerProfile(ctx context.Context, playerID string) (*player.Profile, error)
}
