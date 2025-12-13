package router

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleMessage обрабатывает обычные текстовые сообщения
func (r *Router) handleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	userID := msg.From.ID
	state := r.stateManager.GetState(userID)

	// Проверяем ожидание ввода ФИО
	if strings.HasPrefix(state.WaitingForInput, "fio_") {
		if err := r.fioInputHandler.HandleFioInput(ctx, bot, msg); err != nil {
			// Ошибка уже залогирована внутри
		}
		return
	}

	// По умолчанию - предлагаем использовать /start
	defaultMsg := tgbotapi.NewMessage(msg.Chat.ID, "Используйте /start для начала работы")
	bot.Send(defaultMsg)
}
