package player

import "context"

// Repository определяет интерфейс для работы с хранилищем игроков
type Repository interface {
	Create(ctx context.Context, player *Player) error
	CreateBatch(ctx context.Context, players []*Player) error
	GetByID(ctx context.Context, id string) (*Player, error)
	GetByProfileURL(ctx context.Context, url string) (*Player, error)
	GetByExternalID(ctx context.Context, externalID, source string) (*Player, error)
	ExistsByExternalID(ctx context.Context, externalID, source string) (bool, error)
	Search(ctx context.Context, filters SearchFilters) ([]*Player, error)
	SearchWithTeam(ctx context.Context, filters SearchFilters) ([]*PlayerWithTeam, int, error)
	Update(ctx context.Context, player *Player) error
	Upsert(ctx context.Context, player *Player) error
}

// SearchFilters фильтры для поиска игроков
type SearchFilters struct {
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

// PlayerWithTeam игрок с информацией о команде
type PlayerWithTeam struct {
	Player   *Player
	TeamName string
	TeamCity string
}
