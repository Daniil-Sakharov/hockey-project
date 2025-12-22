package modules

import (
	"fmt"
	"time"
)

// ParsingConfig конфигурация для парсеров
type ParsingConfig struct {
	// Junior parser settings
	Junior JuniorConfig `env:"JUNIOR"`

	// FHSPB parser settings
	FHSPB FHSPBConfig `env:"FHSPB"`

	// General parsing settings
	RequestTimeout  time.Duration `env:"PARSING_REQUEST_TIMEOUT" default:"30s"`
	MaxRetries      int           `env:"PARSING_MAX_RETRIES" validate:"min=0,max=10" default:"3"`
	RetryDelay      time.Duration `env:"PARSING_RETRY_DELAY" default:"5s"`
	UserAgent       string        `env:"PARSING_USER_AGENT" default:"HockeyBot/1.0"`
	EnableRetryJobs bool          `env:"PARSING_ENABLE_RETRY_JOBS" default:"true"`
}

// JuniorConfig конфигурация для Junior парсера
type JuniorConfig struct {
	BaseURL          string        `env:"JUNIOR_BASE_URL" validate:"url" default:"https://junior.fhr.ru"`
	RequestDelay     time.Duration `env:"JUNIOR_REQUEST_DELAY" default:"100ms"`
	WorkerCount      int           `env:"JUNIOR_WORKER_COUNT" validate:"min=1,max=50" default:"10"`
	DomainWorkers    int           `env:"JUNIOR_DOMAIN_WORKERS" validate:"min=1,max=20" default:"5"`
	MinBirthYear     int           `env:"JUNIOR_MIN_BIRTH_YEAR" validate:"min=2000,max=2020" default:"2008"`
	BatchSize        int           `env:"JUNIOR_BATCH_SIZE" validate:"min=1,max=1000" default:"100"`
	EnableAllSeasons bool          `env:"JUNIOR_ENABLE_ALL_SEASONS" default:"true"`
}

// FHSPBConfig конфигурация для FHSPB парсера
type FHSPBConfig struct {
	BaseURL           string        `env:"FHSPB_BASE_URL" validate:"url" default:"https://fhspb.ru"`
	RequestDelay      time.Duration `env:"FHSPB_REQUEST_DELAY" default:"150ms"`
	MaxBirthYear      int           `env:"FHSPB_MAX_BIRTH_YEAR" validate:"min=2000,max=2020" default:"2008"`
	TournamentWorkers int           `env:"FHSPB_TOURNAMENT_WORKERS" validate:"min=1,max=10" default:"3"`
	TeamWorkers       int           `env:"FHSPB_TEAM_WORKERS" validate:"min=1,max=20" default:"5"`
	PlayerWorkers     int           `env:"FHSPB_PLAYER_WORKERS" validate:"min=1,max=50" default:"10"`
	StatisticsWorkers int           `env:"FHSPB_STATISTICS_WORKERS" validate:"min=1,max=10" default:"3"`
	Mode              string        `env:"FHSPB_MODE" validate:"oneof=FULL INCREMENTAL" default:"INCREMENTAL"`
	RetryEnabled      bool          `env:"FHSPB_RETRY_ENABLED" default:"true"`
	RetryMaxAttempts  int           `env:"FHSPB_RETRY_MAX_ATTEMPTS" validate:"min=1,max=10" default:"3"`
	RetryDelay        time.Duration `env:"FHSPB_RETRY_DELAY" default:"5m"`
}

// IsValid проверяет валидность конфигурации парсинга
func (c *ParsingConfig) IsValid() error {
	if err := c.Junior.IsValid(); err != nil {
		return err
	}
	if err := c.FHSPB.IsValid(); err != nil {
		return err
	}
	return nil
}

// IsValid проверяет валидность конфигурации Junior парсера
func (c *JuniorConfig) IsValid() error {
	if c.BaseURL == "" {
		return fmt.Errorf("junior base URL is required")
	}
	return nil
}

// IsValid проверяет валидность конфигурации FHSPB парсера
func (c *FHSPBConfig) IsValid() error {
	if c.BaseURL == "" {
		return fmt.Errorf("fhspb base URL is required")
	}
	return nil
}
