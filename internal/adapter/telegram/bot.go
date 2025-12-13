package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot представляет Telegram бота
type Bot struct {
	api    *tgbotapi.BotAPI
	router *router.Router
	config config.TelegramConfig
}

// NewBot создает новый экземпляр бота
func NewBot(
	cfg config.TelegramConfig,
	r *router.Router,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken())
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = cfg.Debug()

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:    api,
		router: r,
		config: cfg,
	}, nil
}

// Start запускает бота в режиме long polling
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Bot started. Waiting for updates...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Bot stopped")
			b.api.StopReceivingUpdates()
			return ctx.Err()

		case update := <-updates:
			// Обрабатываем обновление в отдельной горутине
			go b.router.Route(ctx, b.api, update)
		}
	}
}

// Stop останавливает бота
func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
	log.Println("Bot stopped")
}
