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
	RequestDelay      time.Duration `env:"FHSPB_REQUEST_DELAY" envDefault:"150ms"`
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
func (c *fhspbConfig) RequestDelay() time.Duration { return c.raw.RequestDelay }
