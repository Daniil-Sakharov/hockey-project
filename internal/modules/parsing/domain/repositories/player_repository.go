package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// PlayerRepository интерфейс для работы с игроками
type PlayerRepository interface {
	Create(ctx context.Context, player *entities.Player) error
	CreateBatch(ctx context.Context, players []*entities.Player) error
	GetByID(ctx context.Context, id string) (*entities.Player, error)
	GetByProfileURL(ctx context.Context, url string) (*entities.Player, error)
	GetByExternalID(ctx context.Context, externalID, source string) (*entities.Player, error)
	ExistsByExternalID(ctx context.Context, externalID, source string) (bool, error)
	Update(ctx context.Context, player *entities.Player) error
	Upsert(ctx context.Context, player *entities.Player) error
}

// PlayerSearchFilters фильтры для поиска игроков
type PlayerSearchFilters struct {
	Name      string
	FirstName string
	LastName  string
	Position  string
	BirthYear *int
	MinHeight *int
	MaxHeight *int
	MinWeight *int
	MaxWeight *int
	Region    string
	Limit     int
	Offset    int
}
