package repositories

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// TeamRepository интерфейс для работы с командами
type TeamRepository interface {
	Create(ctx context.Context, team *entities.Team) error
	CreateBatch(ctx context.Context, teams []*entities.Team) error
	GetByID(ctx context.Context, id string) (*entities.Team, error)
	GetByURL(ctx context.Context, url string) (*entities.Team, error)
	GetByName(ctx context.Context, name, source string) (*entities.Team, error)
	GetByNameAndCity(ctx context.Context, name, city, source string) (*entities.Team, error)
	List(ctx context.Context) ([]*entities.Team, error)
	Upsert(ctx context.Context, team *entities.Team) error
	UpdateCity(ctx context.Context, id, city string) error
	UpdateLogoURL(ctx context.Context, id, logoURL string) error
}
