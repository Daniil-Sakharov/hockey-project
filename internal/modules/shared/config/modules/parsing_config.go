package modules

import (
	"fmt"
	"time"
)

// ParsingConfig конфигурация для парсеров
type ParsingConfig struct {
	// Junior parser settings
	Junior JuniorConfig `env:"JUNIOR"`

	// Junior stats settings
	JuniorStats JuniorStatsConfig `env:"JUNIOR_STATS"`

	// FHSPB parser settings
	FHSPB FHSPBConfig `env:"FHSPB"`

	// FHSPB stats settings
	FHSPBStats FHSPBStatsConfig `env:"FHSPB_STATS"`

	// MIHF parser settings
	MIHF MIHFConfig `env:"MIHF"`

	// FHMoscow parser settings
	FHMoscow FHMoscowConfig `env:"FHMOSCOW"`

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
	RequestDelay     time.Duration `env:"JUNIOR_REQUEST_DELAY" default:"500ms"`
	WorkerCount      int           `env:"JUNIOR_WORKER_COUNT" validate:"min=1,max=50" default:"3"`
	DomainWorkers    int           `env:"JUNIOR_DOMAIN_WORKERS" validate:"min=1,max=20" default:"2"`
	MinBirthYear     int           `env:"JUNIOR_MIN_BIRTH_YEAR" validate:"min=2000,max=2020" default:"2008"`
	BatchSize        int           `env:"JUNIOR_BATCH_SIZE" validate:"min=1,max=1000" default:"100"`
	EnableAllSeasons bool          `env:"JUNIOR_ENABLE_ALL_SEASONS" default:"true"`
}

// JuniorStatsConfig конфигурация для Junior Stats парсера
type JuniorStatsConfig struct {
	RequestDelay      time.Duration `env:"JUNIOR_STATS_REQUEST_DELAY" default:"500ms"`
	TournamentWorkers int           `env:"JUNIOR_STATS_TOURNAMENT_WORKERS" validate:"min=1,max=10" default:"2"`
	BatchSize         int           `env:"JUNIOR_STATS_BATCH_SIZE" validate:"min=1,max=500" default:"100"`
	SkipExisting      bool          `env:"JUNIOR_STATS_SKIP_EXISTING" default:"false"`
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

// FHSPBStatsConfig конфигурация для FHSPB Stats парсера
type FHSPBStatsConfig struct {
	RequestDelay      time.Duration `env:"FHSPB_STATS_REQUEST_DELAY" default:"150ms"`
	TournamentWorkers int           `env:"FHSPB_STATS_TOURNAMENT_WORKERS" validate:"min=1,max=10" default:"3"`
	BatchSize         int           `env:"FHSPB_STATS_BATCH_SIZE" validate:"min=1,max=500" default:"100"`
	SkipExisting      bool          `env:"FHSPB_STATS_SKIP_EXISTING" default:"false"`
}

// MIHFConfig конфигурация для MIHF парсера (stats.mihf.ru)
type MIHFConfig struct {
	BaseURL           string        `env:"MIHF_BASE_URL" validate:"url" default:"https://stats.mihf.ru"`
	RequestDelay      time.Duration `env:"MIHF_REQUEST_DELAY" default:"150ms"`
	MinBirthYear      int           `env:"MIHF_MIN_BIRTH_YEAR" validate:"min=2000,max=2020" default:"2008"`
	MaxBirthYear      int           `env:"MIHF_MAX_BIRTH_YEAR" default:"0"`
	SeasonWorkers     int           `env:"MIHF_SEASON_WORKERS" validate:"min=1,max=5" default:"2"`
	TournamentWorkers int           `env:"MIHF_TOURNAMENT_WORKERS" validate:"min=1,max=10" default:"3"`
	TeamWorkers       int           `env:"MIHF_TEAM_WORKERS" validate:"min=1,max=20" default:"5"`
	PlayerWorkers     int           `env:"MIHF_PLAYER_WORKERS" validate:"min=1,max=50" default:"10"`
	RetryEnabled      bool          `env:"MIHF_RETRY_ENABLED" default:"true"`
	RetryMaxAttempts  int           `env:"MIHF_RETRY_MAX_ATTEMPTS" validate:"min=1,max=10" default:"3"`
	RetryDelay        time.Duration `env:"MIHF_RETRY_DELAY" default:"5m"`
	MaxSeasons        int           `env:"MIHF_MAX_SEASONS" default:"0"`
	TestSeason        string        `env:"MIHF_TEST_SEASON" default:""`
}

// FHMoscowConfig конфигурация для FHMoscow парсера (fhmoscow.com)
type FHMoscowConfig struct {
	BaseURL           string        `env:"FHMOSCOW_BASE_URL" validate:"url" default:"https://www.fhmoscow.com"`
	RequestDelay      time.Duration `env:"FHMOSCOW_REQUEST_DELAY" default:"150ms"`
	MinBirthYear      int           `env:"FHMOSCOW_MIN_BIRTH_YEAR" validate:"min=2000,max=2020" default:"2008"`
	SeasonWorkers     int           `env:"FHMOSCOW_SEASON_WORKERS" validate:"min=1,max=5" default:"2"`
	TournamentWorkers int           `env:"FHMOSCOW_TOURNAMENT_WORKERS" validate:"min=1,max=10" default:"3"`
	TeamWorkers       int           `env:"FHMOSCOW_TEAM_WORKERS" validate:"min=1,max=20" default:"5"`
	PlayerWorkers     int           `env:"FHMOSCOW_PLAYER_WORKERS" validate:"min=1,max=50" default:"10"`
	RetryEnabled      bool          `env:"FHMOSCOW_RETRY_ENABLED" default:"true"`
	RetryMaxAttempts  int           `env:"FHMOSCOW_RETRY_MAX_ATTEMPTS" validate:"min=1,max=10" default:"3"`
	RetryDelay        time.Duration `env:"FHMOSCOW_RETRY_DELAY" default:"5m"`
	MaxSeasons        int           `env:"FHMOSCOW_MAX_SEASONS" default:"0"`
	TestSeason        string        `env:"FHMOSCOW_TEST_SEASON" default:""`
	// Player scanning (since team roster pages are JavaScript-rendered)
	ScanPlayers bool `env:"FHMOSCOW_SCAN_PLAYERS" default:"true"`
	MaxPlayerID int  `env:"FHMOSCOW_MAX_PLAYER_ID" default:"15000"`
}

// IsValid проверяет валидность конфигурации парсинга
func (c *ParsingConfig) IsValid() error {
	if err := c.Junior.IsValid(); err != nil {
		return err
	}
	if err := c.FHSPB.IsValid(); err != nil {
		return err
	}
	if err := c.MIHF.IsValid(); err != nil {
		return err
	}
	if err := c.FHMoscow.IsValid(); err != nil {
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

// IsValid проверяет валидность конфигурации MIHF парсера
func (c *MIHFConfig) IsValid() error {
	if c.BaseURL == "" {
		return fmt.Errorf("mihf base URL is required")
	}
	return nil
}

// IsValid проверяет валидность конфигурации FHMoscow парсера
func (c *FHMoscowConfig) IsValid() error {
	if c.BaseURL == "" {
		return fmt.Errorf("fhmoscow base URL is required")
	}
	return nil
}
