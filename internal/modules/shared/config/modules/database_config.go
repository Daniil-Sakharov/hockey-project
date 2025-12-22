package modules

import (
	"fmt"
	"time"
)

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host            string        `env:"POSTGRES_HOST" validate:"required" default:"localhost"`
	Port            int           `env:"POSTGRES_PORT" validate:"min=1,max=65535" default:"5432"`
	User            string        `env:"POSTGRES_USER" validate:"required" default:"hockey"`
	Password        string        `env:"POSTGRES_PASSWORD" validate:"required"`
	Database        string        `env:"POSTGRES_DB" validate:"required" default:"hockey_stats"`
	SSLMode         string        `env:"POSTGRES_SSL_MODE" validate:"oneof=disable require verify-ca verify-full" default:"disable"`
	MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS" validate:"min=1" default:"25"`
	MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS" validate:"min=1" default:"5"`
	ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" default:"5m"`
	ConnMaxIdleTime time.Duration `env:"POSTGRES_CONN_MAX_IDLE_TIME" default:"1m"`
}

// URI возвращает строку подключения к PostgreSQL
func (c *DatabaseConfig) URI() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

// IsValid проверяет валидность конфигурации
func (c *DatabaseConfig) IsValid() error {
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}
