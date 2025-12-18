package config

import "time"

// LoggerConfig конфигурация логгера
type LoggerConfig interface {
	Level() string
	AsJson() bool
}

// PostgresConfig конфигурация PostgreSQL
type PostgresConfig interface {
	Host() string
	Port() string
	User() string
	Password() string
	Database() string
	SSLMode() string
	MigrationsDir() string
	URI() string
}

// ParserConfig конфигурация парсера
type ParserConfig interface {
	BaseURL() string
	Timeout() int
}

// TelegramConfig конфигурация Telegram бота
type TelegramConfig interface {
	BotToken() string
	Debug() bool
}

// RegistryConfig конфигурация registrynew.fhr.ru
type RegistryConfig interface {
	URL() string
	Username() string
	Password() string
}

// FHSPBConfig конфигурация парсера fhspb.ru
type FHSPBConfig interface {
	MaxBirthYear() int
	TournamentWorkers() int
	TeamWorkers() int
	PlayerWorkers() int
	StatisticsWorkers() int
	RequestDelay() time.Duration
	Mode() string
	RetryEnabled() bool
	RetryMaxAttempts() int
	RetryDelay() time.Duration
}

// JuniorConfig конфигурация парсера junior.fhr.ru
type JuniorConfig interface {
	BaseURL() string
	DomainWorkers() int
	MinBirthYear() int
}
