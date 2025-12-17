package telegram

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
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

	ctx := context.Background()
	logger.Info(ctx, "✅ Authorized on Telegram account",
		zap.String("username", api.Self.UserName),
		zap.Int64("bot_id", api.Self.ID))

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

	logger.Info(ctx, "🚀 Bot started. Waiting for updates...")

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "🛑 Bot stopped by context")
			b.api.StopReceivingUpdates()
			return ctx.Err()

		case update := <-updates:
			b.logUpdate(ctx, update)

			// Обрабатываем обновление в отдельной горутине с восстановлением от паники
			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Error(ctx, "💥 Panic in update handler",
							zap.Any("panic", r),
							zap.Int("update_id", update.UpdateID))
					}
				}()
				b.router.Route(ctx, b.api, update)
			}()
		}
	}
}

// Stop останавливает бота
func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
	ctx := context.Background()
	logger.Info(ctx, "🛑 Bot stopped")
}

// logUpdate логирует информацию о полученном обновлении
func (b *Bot) logUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message != nil {
		logger.Info(ctx, "📨 Received message",
			zap.Int("update_id", update.UpdateID),
			zap.Int64("user_id", update.Message.From.ID),
			zap.String("username", update.Message.From.UserName),
			zap.String("text", update.Message.Text))
	} else if update.CallbackQuery != nil {
		logger.Info(ctx, "🔘 Received callback",
			zap.Int("update_id", update.UpdateID),
			zap.Int64("user_id", update.CallbackQuery.From.ID),
			zap.String("username", update.CallbackQuery.From.UserName),
			zap.String("data", update.CallbackQuery.Data))
	}
}
