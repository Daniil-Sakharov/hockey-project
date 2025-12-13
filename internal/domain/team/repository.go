package team

import "context"

// Repository определяет интерфейс для работы с хранилищем команд
type Repository interface {
	// Create создает новую команду
	Create(ctx context.Context, team *Team) error

	// Upsert создает или обновляет команду (ON CONFLICT DO UPDATE)
	Upsert(ctx context.Context, team *Team) (*Team, error)

	// CreateBatch создает несколько команд за одну транзакцию
	CreateBatch(ctx context.Context, teams []*Team) error

	// GetByID возвращает команду по ID
	GetByID(ctx context.Context, id string) (*Team, error)

	// GetByURL возвращает команду по URL (для дедупликации)
	GetByURL(ctx context.Context, url string) (*Team, error)

	// List возвращает список всех команд
	List(ctx context.Context, limit, offset int) ([]*Team, error)
}
