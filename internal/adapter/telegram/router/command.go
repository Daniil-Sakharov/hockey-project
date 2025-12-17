package router

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleCommand обрабатывает команды
func (r *Router) handleCommand(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	logger.Info(ctx, "⚡ Routing command",
		zap.String("command", message.Command()),
		zap.Int64("user_id", message.From.ID))

	switch message.Command() {
	case "start":
		if err := r.startHandler.Handle(ctx, bot, message); err != nil {
			logger.Error(ctx, "❌ Error handling /start", zap.Error(err))
		}
	default:
		logger.Info(ctx, "❓ Unknown command", zap.String("command", message.Command()))
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Используйте /start")
		if _, err := bot.Send(msg); err != nil {
			logger.Error(ctx, "❌ Error sending unknown command message", zap.Error(err))
		}
	}
}
