package router

import (
	"context"
	"log"
	"strings"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleCallback обрабатывает callback query (нажатия на inline кнопки)
func (r *Router) handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		return
	}

	action := parts[0]

	// Отправляем ответ на callback (убирает "часики" загрузки)
	callbackResp := tgbotapi.NewCallback(query.ID, "")
	if _, err := bot.Request(callbackResp); err != nil {
		log.Printf("Error answering callback: %v", err)
		return
	}

	// Создаем logger для callback handlers
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	// Маршрутизация по типу действия
	switch action {
	case cb.ActionMenu:
		callback.HandleMenu(r, ctx, bot, query, parts)
	case cb.ActionFilter:
		callback.HandleFilter(r, ctx, bot, query, parts)
	case cb.ActionSearch:
		callback.HandleSearch(r, ctx, bot, query, parts)
	case cb.ActionPlayer:
		callback.HandleProfile(r, ctx, bot, query, parts, logger)
	case cb.ActionReport:
		callback.HandleReport(r, ctx, bot, query, logger)
	default:
		log.Printf("Unknown callback action: %s", action)
	}
}
