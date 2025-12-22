package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// PlayerStatisticsRepository интерфейс для работы со статистикой игроков
type PlayerStatisticsRepository interface {
	CreateBatch(ctx context.Context, stats []*entities.PlayerStatistic) (int, error)
	DeleteByTournament(ctx context.Context, tournamentID string) error
	DeleteAll(ctx context.Context) error
	GetByPlayerID(ctx context.Context, playerID string) ([]*entities.PlayerStatistic, error)
	GetByTournament(ctx context.Context, tournamentID string) ([]*entities.PlayerStatistic, error)
	CountAll(ctx context.Context) (int, error)
}

// PlayerTeamRepository интерфейс для работы со связями игрок-команда
type PlayerTeamRepository interface {
	Create(ctx context.Context, link *entities.PlayerTeam) error
	CreateBatch(ctx context.Context, links []*entities.PlayerTeam) error
	GetByPlayerID(ctx context.Context, playerID string) ([]*entities.PlayerTeam, error)
	GetByTeamID(ctx context.Context, teamID string) ([]*entities.PlayerTeam, error)
	Exists(ctx context.Context, playerID, teamID, tournamentID string) (bool, error)
	Upsert(ctx context.Context, link *entities.PlayerTeam) error
}
