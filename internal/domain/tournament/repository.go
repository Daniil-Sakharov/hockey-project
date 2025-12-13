package tournament

import "context"

// Repository определяет интерфейс для работы с хранилищем турниров
type Repository interface {
	// Create создает новый турнир
	Create(ctx context.Context, tournament *Tournament) error

	// CreateBatch создает несколько турниров за одну транзакцию
	CreateBatch(ctx context.Context, tournaments []*Tournament) error

	// GetByID возвращает турнир по ID
	GetByID(ctx context.Context, id string) (*Tournament, error)

	// GetByURL возвращает турнир по URL (для дедупликации)
	GetByURL(ctx context.Context, url string) (*Tournament, error)

	// Update обновляет турнир
	Update(ctx context.Context, tournament *Tournament) error

	// List возвращает список всех турниров
	List(ctx context.Context, limit, offset int) ([]*Tournament, error)
}
