package player

import "context"

type UserService interface {
	Create(ctx context.Context, player *Player) error
	GetByID(ctx context.Context, id int) (*Player, error)
	Update(ctx context.Context)
}
