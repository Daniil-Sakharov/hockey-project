package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// MatchTeamStatsRepository интерфейс для работы со статистикой команд за матч
type MatchTeamStatsRepository interface {
	Create(ctx context.Context, stats *entities.MatchTeamStats) error
	CreateBatch(ctx context.Context, stats []*entities.MatchTeamStats) error
	Upsert(ctx context.Context, stats *entities.MatchTeamStats) error
	UpsertBatch(ctx context.Context, stats []*entities.MatchTeamStats) error

	GetByMatch(ctx context.Context, matchID string) ([]*entities.MatchTeamStats, error)
	GetByMatchAndTeam(ctx context.Context, matchID, teamID string) (*entities.MatchTeamStats, error)

	DeleteByMatch(ctx context.Context, matchID string) error
}
