package team

import "context"

// Repository определяет интерфейс для работы с хранилищем команд
type Repository interface {
	Create(ctx context.Context, team *Team) error
	CreateBatch(ctx context.Context, teams []*Team) error
	GetByID(ctx context.Context, id string) (*Team, error)
	GetByURL(ctx context.Context, url string) (*Team, error)
	List(ctx context.Context) ([]*Team, error)
	Upsert(ctx context.Context, team *Team) error
}
