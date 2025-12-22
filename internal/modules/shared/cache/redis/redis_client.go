package redis

import (
	pkgredis "github.com/Daniil-Sakharov/HockeyProject/pkg/cache/redis"
)

// Client is an alias to pkg/cache/redis client.
// Use pkg/cache/redis directly for full functionality.
type Client = pkgredis.Logger

// NewClient creates a redis client using pkg/cache/redis.
// This is a convenience wrapper - prefer using pkg/cache/redis directly.
