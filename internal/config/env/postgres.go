package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port          string `env:"POSTGRES_PORT" envDefault:"5432"`
	User          string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password      string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	Database      string `env:"POSTGRES_DB" envDefault:"hockey_stats"`
	SSLMode       string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	MigrationsDir string `env:"MIGRATIONS_DIR" envDefault:"migrations"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (c *postgresConfig) Host() string {
	return c.raw.Host
}

func (c *postgresConfig) Port() string {
	return c.raw.Port
}

func (c *postgresConfig) User() string {
	return c.raw.User
}

func (c *postgresConfig) Password() string {
	return c.raw.Password
}

func (c *postgresConfig) Database() string {
	return c.raw.Database
}

func (c *postgresConfig) SSLMode() string {
	return c.raw.SSLMode
}

func (c *postgresConfig) MigrationsDir() string {
	return c.raw.MigrationsDir
}

func (c *postgresConfig) URI() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.raw.Host, c.raw.Port, c.raw.User, c.raw.Password, c.raw.Database, c.raw.SSLMode,
	)
}
