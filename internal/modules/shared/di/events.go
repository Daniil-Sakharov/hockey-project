package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/events/store"
)

// EventStore возвращает Event Store
func (c *Container) EventStore(ctx context.Context) (*store.SQLXEventStore, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.eventStore != nil {
		return c.eventStore, nil
	}

	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}

	c.eventStore = store.NewSQLXEventStore(db)
	return c.eventStore, nil
}
