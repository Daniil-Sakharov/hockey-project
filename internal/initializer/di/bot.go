package di

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
)

// BotContainer содержит зависимости для Telegram бота
type BotContainer struct {
	*BaseContainer
	telegram *Telegram
}

func NewBotContainer(cfg *config.Config) *BotContainer {
	base := NewBaseContainer(cfg)
	telegram := NewTelegram(cfg, base.service)

	return &BotContainer{
		BaseContainer: base,
		telegram:      telegram,
	}
}

func (bc *BotContainer) Telegram() *Telegram {
	return bc.telegram
}
