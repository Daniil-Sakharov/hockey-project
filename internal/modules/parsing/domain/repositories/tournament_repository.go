package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// TournamentRepository интерфейс для работы с турнирами
type TournamentRepository interface {
	Create(ctx context.Context, tournament *entities.Tournament) error
	CreateBatch(ctx context.Context, tournaments []*entities.Tournament) error
	GetByID(ctx context.Context, id string) (*entities.Tournament, error)
	GetByURL(ctx context.Context, url string) (*entities.Tournament, error)
	GetBySource(ctx context.Context, source string) ([]*entities.Tournament, error)
	GetWithTeams(ctx context.Context, source string) ([]*entities.Tournament, error) // Турниры где есть команды
	List(ctx context.Context) ([]*entities.Tournament, error)
	Update(ctx context.Context, tournament *entities.Tournament) error
	UpdateBirthYearGroups(ctx context.Context, id string, birthYearGroupsJSON string) error
}
