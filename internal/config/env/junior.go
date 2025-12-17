package env

import "github.com/caarlos0/env/v11"

type juniorEnvConfig struct {
	BaseURL       string `env:"JUNIOR_BASE_URL" envDefault:"https://junior.fhr.ru"`
	DomainWorkers int    `env:"JUNIOR_DOMAIN_WORKERS" envDefault:"5"`
	MinBirthYear  int    `env:"JUNIOR_MIN_BIRTH_YEAR" envDefault:"2008"`
}

type juniorConfig struct {
	raw juniorEnvConfig
}

func NewJuniorConfig() (*juniorConfig, error) {
	var raw juniorEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &juniorConfig{raw: raw}, nil
}

func (c *juniorConfig) BaseURL() string    { return c.raw.BaseURL }
func (c *juniorConfig) DomainWorkers() int { return c.raw.DomainWorkers }
func (c *juniorConfig) MinBirthYear() int  { return c.raw.MinBirthYear }
