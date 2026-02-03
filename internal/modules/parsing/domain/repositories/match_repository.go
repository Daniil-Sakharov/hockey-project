package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// MatchRepository интерфейс для работы с матчами
type MatchRepository interface {
	Create(ctx context.Context, match *entities.Match) error
	CreateBatch(ctx context.Context, matches []*entities.Match) error
	Update(ctx context.Context, match *entities.Match) error
	Upsert(ctx context.Context, match *entities.Match) error

	GetByID(ctx context.Context, id string) (*entities.Match, error)
	GetByExternalID(ctx context.Context, externalID, source string) (*entities.Match, error)
	GetByTournament(ctx context.Context, tournamentID string) ([]*entities.Match, error)
	GetByTeam(ctx context.Context, teamID string) ([]*entities.Match, error)
	GetUnparsedFinished(ctx context.Context, source string, limit int) ([]*entities.Match, error)

	MarkDetailsParsed(ctx context.Context, id string) error
}
