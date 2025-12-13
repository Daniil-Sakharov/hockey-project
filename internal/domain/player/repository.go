package player

import "context"

// Repository определяет интерфейс для работы с хранилищем игроков
type Repository interface {
	// Create создает нового игрока
	Create(ctx context.Context, player *Player) error

	// CreateBatch создает несколько игроков за одну транзакцию
	CreateBatch(ctx context.Context, players []*Player) error

	// GetByID возвращает игрока по ID
	GetByID(ctx context.Context, id string) (*Player, error)

	// GetByProfileURL возвращает игрока по URL профиля (для дедупликации)
	GetByProfileURL(ctx context.Context, url string) (*Player, error)

	// GetByExternalID возвращает игрока по внешнему ID и источнику
	GetByExternalID(ctx context.Context, externalID, source string) (*Player, error)

	// ExistsByExternalID проверяет существование игрока по внешнему ID и источнику
	ExistsByExternalID(ctx context.Context, externalID, source string) (bool, error)

	// Search ищет игроков по фильтрам
	Search(ctx context.Context, filters SearchFilters) ([]*Player, error)

	// SearchWithTeam ищет игроков с информацией о команде (для бота)
	SearchWithTeam(ctx context.Context, filters SearchFilters) ([]*PlayerWithTeam, int, error)

	// Update обновляет данные игрока
	Update(ctx context.Context, player *Player) error

	// Upsert создает или обновляет игрока по external_id и source
	Upsert(ctx context.Context, player *Player) error
}

// SearchFilters фильтры для поиска игроков
type SearchFilters struct {
	Name      string // Поиск по имени (ILIKE)
	FirstName string // Поиск по имени (ILIKE)
	LastName  string // Поиск по фамилии (ILIKE)
	Position  string // Фильтр по позиции
	BirthYear *int   // Фильтр по году рождения
	MinHeight *int   // Минимальный рост
	MaxHeight *int   // Максимальный рост
	MinWeight *int   // Минимальный вес
	MaxWeight *int   // Максимальный вес
	Limit     int    // Лимит результатов
	Offset    int    // Смещение для пагинации
}

// PlayerWithTeam игрок с информацией о команде
type PlayerWithTeam struct {
	Player   *Player
	TeamName string
	TeamCity string
}
