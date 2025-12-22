package bot

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// RouterInterface интерфейс для router
type RouterInterface interface {
	Route(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)
}

// Bot представляет Telegram бота
type Bot struct {
	api    *tgbotapi.BotAPI
	router RouterInterface
	config *modules.TelegramConfig
}

// NewBot создает новый экземпляр бота
func NewBot(cfg *modules.TelegramConfig, router RouterInterface) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	api.Debug = cfg.EnableDebugMode

	ctx := context.Background()
	logger.Info(ctx, "✅ Authorized on Telegram account",
		zap.String("username", api.Self.UserName),
		zap.Int64("bot_id", api.Self.ID))

	return &Bot{api: api, router: router, config: cfg}, nil
}

// Start запускает бота в режиме long polling
func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.config.UpdateTimeout

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
			go b.handleUpdate(ctx, update)
		}
	}
}

// handleUpdate обрабатывает обновление с recovery
func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(ctx, "💥 Panic in update handler",
				zap.Any("panic", r),
				zap.Int("update_id", update.UpdateID))
		}
	}()
	b.router.Route(ctx, b.api, update)
}

// Stop останавливает бота
func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
	logger.Info(context.Background(), "🛑 Bot stopped")
}

// logUpdate логирует информацию о полученном обновлении
func (b *Bot) logUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message != nil {
		logger.Info(ctx, "📨 Received message",
			zap.Int("update_id", update.UpdateID),
			zap.Int64("user_id", update.Message.From.ID),
			zap.String("text", update.Message.Text))
	} else if update.CallbackQuery != nil {
		logger.Info(ctx, "🔘 Received callback",
			zap.Int("update_id", update.UpdateID),
			zap.Int64("user_id", update.CallbackQuery.From.ID),
			zap.String("data", update.CallbackQuery.Data))
	}
}
