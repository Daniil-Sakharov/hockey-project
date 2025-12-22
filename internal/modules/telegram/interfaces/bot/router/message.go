package router

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleMessage обрабатывает текстовые сообщения
func (r *Router) handleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	logger.Info(ctx, "📝 handleMessage called",
		zap.Int64("user_id", msg.From.ID),
		zap.String("text", msg.Text),
		zap.Bool("is_command", msg.IsCommand()))

	if msg.IsCommand() {
		logger.Debug(ctx, "➡️ Message is command, routing to handleCommand")
		r.handleCommand(ctx, bot, msg)
		return
	}

	// Текстовый ввод (FIO)
	logger.Debug(ctx, "➡️ Routing to filterHandler.HandleTextInput")
	if err := r.filterHandler.HandleTextInput(ctx, bot, msg); err != nil {
		logger.Error(ctx, "❌ Error handling text input", zap.Error(err))
	}
}
