package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// StandingRepository интерфейс для работы с турнирными таблицами
type StandingRepository interface {
	Create(ctx context.Context, standing *entities.TeamStanding) error
	CreateBatch(ctx context.Context, standings []*entities.TeamStanding) error
	Upsert(ctx context.Context, standing *entities.TeamStanding) error
	UpsertBatch(ctx context.Context, standings []*entities.TeamStanding) error

	GetByTournament(ctx context.Context, tournamentID string) ([]*entities.TeamStanding, error)
	GetByTeam(ctx context.Context, teamID string) ([]*entities.TeamStanding, error)
	GetByTournamentAndGroup(ctx context.Context, tournamentID, groupName string) ([]*entities.TeamStanding, error)

	DeleteByTournament(ctx context.Context, tournamentID string) error
}
