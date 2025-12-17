package tournament

import "context"

// Repository определяет интерфейс для работы с хранилищем турниров
type Repository interface {
	Create(ctx context.Context, tournament *Tournament) error
	CreateBatch(ctx context.Context, tournaments []*Tournament) error
	GetByID(ctx context.Context, id string) (*Tournament, error)
	GetByURL(ctx context.Context, url string) (*Tournament, error)
	List(ctx context.Context) ([]*Tournament, error)
	Update(ctx context.Context, tournament *Tournament) error
}
