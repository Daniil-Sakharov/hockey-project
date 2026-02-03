package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// MatchEventRepository интерфейс для работы с событиями матчей
type MatchEventRepository interface {
	Create(ctx context.Context, event *entities.MatchEvent) error
	CreateBatch(ctx context.Context, events []*entities.MatchEvent) error

	GetByMatchID(ctx context.Context, matchID string) ([]*entities.MatchEvent, error)
	GetGoalsByMatchID(ctx context.Context, matchID string) ([]*entities.MatchEvent, error)
	GetPenaltiesByMatchID(ctx context.Context, matchID string) ([]*entities.MatchEvent, error)
	GetByPlayerID(ctx context.Context, playerID string) ([]*entities.MatchEvent, error)

	DeleteByMatchID(ctx context.Context, matchID string) error
}
