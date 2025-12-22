package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/cache/strategies"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
)

const (
	sessionKeyPrefix = "telegram:session:"
	sessionTTL       = 30 * time.Minute
)

// SessionCache implements SessionRepository with caching.
type SessionCache struct {
	cache strategies.CacheStrategy
}

// NewSessionCache creates a new cached session repository.
func NewSessionCache(cache strategies.CacheStrategy) *SessionCache {
	return &SessionCache{cache: cache}
}

// Get retrieves session from cache.
func (c *SessionCache) Get(userID int64) (*entities.UserSession, error) {
	key := c.key(userID)
	data, err := c.cache.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("session not found")
	}

	var session entities.UserSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

// Save stores session in cache.
func (c *SessionCache) Save(session *entities.UserSession) error {
	key := c.key(session.UserID)
	return c.cache.Set(context.Background(), key, session, sessionTTL)
}

// Delete removes session from cache.
func (c *SessionCache) Delete(userID int64) error {
	key := c.key(userID)
	return c.cache.Delete(context.Background(), key)
}

func (c *SessionCache) key(userID int64) string {
	return fmt.Sprintf("%s%d", sessionKeyPrefix, userID)
}

// Ensure SessionCache implements SessionRepository.
var _ domain.SessionRepository = (*SessionCache)(nil)
