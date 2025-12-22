package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/cache"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/cache/redis"
	redigo "github.com/gomodule/redigo/redis"
)

// RedisClient возвращает клиент Redis
func (c *Container) RedisClient(ctx context.Context) (cache.RedisClient, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.redisClient != nil {
		return c.redisClient, nil
	}

	redisConfig, err := c.configContainer.Redis(ctx)
	if err != nil {
		return nil, err
	}

	if !redisConfig.IsEnabled() {
		return nil, nil // Redis отключен
	}

	pool := &redigo.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		MaxActive:   redisConfig.MaxActive,
		IdleTimeout: redisConfig.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			conn, err := redigo.Dial("tcp", redisConfig.Address())
			if err != nil {
				return nil, err
			}
			if redisConfig.Password != "" {
				if _, err := conn.Do("AUTH", redisConfig.Password); err != nil {
					_ = conn.Close()
					return nil, err
				}
			}
			if redisConfig.DB != 0 {
				if _, err := conn.Do("SELECT", redisConfig.DB); err != nil {
					_ = conn.Close()
					return nil, err
				}
			}
			return conn, nil
		},
	}

	c.redisClient = redis.NewClient(pool, nil, redisConfig.IdleTimeout)
	return c.redisClient, nil
}

// HasRedis проверяет доступен ли Redis
func (c *Container) HasRedis(ctx context.Context) bool {
	redisConfig, err := c.configContainer.Redis(ctx)
	if err != nil {
		return false
	}
	return redisConfig.IsEnabled()
}
