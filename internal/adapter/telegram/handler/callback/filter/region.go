package filter

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *FilterHandler) HandleRegionSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		"🗺️ **Регион (Федеральный округ)**\n\nВыберите федеральный округ:",
	)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.RegionSelect()
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

func (h *FilterHandler) HandleRegionValue(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if value == "any" {
		state.Filters.Region = nil
	} else {
		state.Filters.Region = &value
	}

	h.stateManager.UpdateFilters(userID, state.Filters)
	return h.HandleFilterMenu(ctx, botAPI, query)
}
