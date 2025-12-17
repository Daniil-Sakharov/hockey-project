package env

import (
	"github.com/caarlos0/env/v11"
)

type telegramEnvConfig struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	Debug    bool   `env:"TELEGRAM_DEBUG" envDefault:"false"`
}

type telegramConfig struct {
	raw telegramEnvConfig
}

func NewTelegramConfig() (*telegramConfig, error) {
	var raw telegramEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &telegramConfig{raw: raw}, nil
}

func (c *telegramConfig) BotToken() string {
	return c.raw.BotToken
}

func (c *telegramConfig) Debug() bool {
	return c.raw.Debug
}
