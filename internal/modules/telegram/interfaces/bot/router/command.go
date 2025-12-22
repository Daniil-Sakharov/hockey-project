package router

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleCommand обрабатывает команды
func (r *Router) handleCommand(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	cmd := msg.Command()

	logger.Info(ctx, "🎯 handleCommand called",
		zap.String("command", cmd),
		zap.Int64("user_id", msg.From.ID),
		zap.String("username", msg.From.UserName))

	switch cmd {
	case "start":
		logger.Debug(ctx, "➡️ Calling startHandler.HandleStart")
		if err := r.startHandler.HandleStart(ctx, bot, msg); err != nil {
			logger.Error(ctx, "❌ Error handling /start", zap.Error(err))
		}
	default:
		logger.Warn(ctx, "⚠️ Unknown command", zap.String("command", cmd))
	}
}
