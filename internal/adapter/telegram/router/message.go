package router

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) handleMessage(ctx context.Context, botAPI *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	userID := msg.From.ID
	state := r.stateManager.GetState(userID)

	if strings.HasPrefix(state.WaitingForInput, "fio_") {
		if err := r.fioInputHandler.HandleFioInput(ctx, botAPI, msg); err != nil {
			// Ошибка уже залогирована внутри
		}
		return
	}

	defaultMsg := tgbotapi.NewMessage(msg.Chat.ID, "Используйте /start для начала работы")
	botAPI.Send(defaultMsg)
}
