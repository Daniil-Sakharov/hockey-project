package domain

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
)

// PlayerRepository интерфейс для доступа к данным игроков
type PlayerRepository interface {
	SearchWithFilters(ctx context.Context, filters services.SearchFilters) ([]*services.PlayerWithTeam, int, error)
	GetByID(ctx context.Context, playerID string) (*services.PlayerProfile, error)
}

// SessionRepository интерфейс для хранения сессий пользователей
type SessionRepository interface {
	Get(userID int64) (*entities.UserSession, error)
	Save(session *entities.UserSession) error
	Delete(userID int64) error
}
