package router

import (
	"context"
	"strings"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/router/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleCallback обрабатывает callback query
func (r *Router) handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, ":")
	if len(parts) < 1 {
		return
	}

	action := parts[0]

	logger.Info(ctx, "Routing callback",
		zap.String("action", action),
		zap.Int64("user_id", query.From.ID))

	// Отвечаем на callback
	callbackResp := tgbotapi.NewCallback(query.ID, "")
	if _, err := bot.Request(callbackResp); err != nil {
		logger.Error(ctx, "Error answering callback", zap.Error(err))
		return
	}

	switch action {
	case cb.ActionMenu:
		callback.HandleMenu(r, ctx, bot, query, parts)
	case cb.ActionFilter:
		callback.HandleFilter(r, ctx, bot, query, parts)
	case cb.ActionSearch:
		callback.HandleSearch(r, ctx, bot, query, parts)
	case cb.ActionPlayer:
		callback.HandleProfile(r, ctx, bot, query, parts)
	case cb.ActionReport:
		callback.HandleReport(r, ctx, bot, query)
	default:
		logger.Warn(ctx, "Unknown callback action", zap.String("action", action))
	}
}
