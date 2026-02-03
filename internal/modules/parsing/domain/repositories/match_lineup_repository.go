package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// MatchLineupRepository интерфейс для работы с составами матчей
type MatchLineupRepository interface {
	Create(ctx context.Context, lineup *entities.MatchLineup) error
	CreateBatch(ctx context.Context, lineups []*entities.MatchLineup) error
	Upsert(ctx context.Context, lineup *entities.MatchLineup) error

	GetByMatchID(ctx context.Context, matchID string) ([]*entities.MatchLineup, error)
	GetByPlayerID(ctx context.Context, playerID string) ([]*entities.MatchLineup, error)
	GetByMatchAndTeam(ctx context.Context, matchID, teamID string) ([]*entities.MatchLineup, error)
	GetByMatchAndJersey(ctx context.Context, matchID string, jerseyNumber int) (*entities.MatchLineup, error)

	DeleteByMatchID(ctx context.Context, matchID string) error
}
