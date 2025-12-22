package modules

import (
	"fmt"
	"time"
)

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Host        string        `env:"REDIS_HOST" default:"localhost"`
	Port        int           `env:"REDIS_PORT" default:"6379"`
	Password    string        `env:"REDIS_PASSWORD"`
	DB          int           `env:"REDIS_DB" default:"0"`
	MaxIdle     int           `env:"REDIS_MAX_IDLE" default:"10"`
	MaxActive   int           `env:"REDIS_MAX_ACTIVE" default:"100"`
	IdleTimeout time.Duration `env:"REDIS_IDLE_TIMEOUT" default:"240s"`
	Enabled     bool          `env:"REDIS_ENABLED" default:"false"`
}

// Address возвращает адрес Redis
func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsEnabled проверяет включен ли Redis
func (c *RedisConfig) IsEnabled() bool {
	return c.Enabled
}

// IsValid проверяет валидность конфигурации
func (c *RedisConfig) IsValid() error {
	if c.Host == "" {
		return fmt.Errorf("redis host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid redis port: %d", c.Port)
	}
	return nil
}
