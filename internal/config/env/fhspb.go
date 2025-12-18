package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type fhspbEnvConfig struct {
	MaxBirthYear      int           `env:"FHSPB_MAX_BIRTH_YEAR" envDefault:"2008"`
	TournamentWorkers int           `env:"FHSPB_TOURNAMENT_WORKERS" envDefault:"3"`
	TeamWorkers       int           `env:"FHSPB_TEAM_WORKERS" envDefault:"5"`
	PlayerWorkers     int           `env:"FHSPB_PLAYER_WORKERS" envDefault:"10"`
	StatisticsWorkers int           `env:"FHSPB_STATISTICS_WORKERS" envDefault:"3"`
	RequestDelay      time.Duration `env:"FHSPB_REQUEST_DELAY" envDefault:"150ms"`
	Mode              string        `env:"FHSPB_MODE" envDefault:"INCREMENTAL"`
	RetryEnabled      bool          `env:"FHSPB_RETRY_ENABLED" envDefault:"true"`
	RetryMaxAttempts  int           `env:"FHSPB_RETRY_MAX_ATTEMPTS" envDefault:"3"`
	RetryDelay        time.Duration `env:"FHSPB_RETRY_DELAY" envDefault:"5m"`
}

type fhspbConfig struct {
	raw fhspbEnvConfig
}

func NewFHSPBConfig() (*fhspbConfig, error) {
	var raw fhspbEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &fhspbConfig{raw: raw}, nil
}

func (c *fhspbConfig) MaxBirthYear() int           { return c.raw.MaxBirthYear }
func (c *fhspbConfig) TournamentWorkers() int      { return c.raw.TournamentWorkers }
func (c *fhspbConfig) TeamWorkers() int            { return c.raw.TeamWorkers }
func (c *fhspbConfig) PlayerWorkers() int          { return c.raw.PlayerWorkers }
func (c *fhspbConfig) StatisticsWorkers() int      { return c.raw.StatisticsWorkers }
func (c *fhspbConfig) RequestDelay() time.Duration { return c.raw.RequestDelay }
func (c *fhspbConfig) Mode() string                { return c.raw.Mode }
func (c *fhspbConfig) RetryEnabled() bool          { return c.raw.RetryEnabled }
func (c *fhspbConfig) RetryMaxAttempts() int       { return c.raw.RetryMaxAttempts }
func (c *fhspbConfig) RetryDelay() time.Duration   { return c.raw.RetryDelay }
