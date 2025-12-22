package strategies

import (
	"context"
	"time"
)

// CacheStrategy interface for caching strategies.
type CacheStrategy interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// CacheClient interface for cache client.
type CacheClient interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}) error
	SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

// Tier1Strategy - fast cache for frequently used data (TTL: 5 min).
type Tier1Strategy struct {
	c CacheClient
}

// NewTier1Strategy creates Tier 1 strategy.
func NewTier1Strategy(client CacheClient) *Tier1Strategy {
	return &Tier1Strategy{c: client}
}

// Get retrieves value from cache.
func (t *Tier1Strategy) Get(ctx context.Context, key string) ([]byte, error) {
	return t.c.Get(ctx, key)
}

// Set stores value in cache.
func (t *Tier1Strategy) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = 5 * time.Minute
	}
	return t.c.SetWithTTL(ctx, key, value, ttl)
}

// Delete removes value from cache.
func (t *Tier1Strategy) Delete(ctx context.Context, key string) error {
	return t.c.Del(ctx, key)
}
