package filter

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleWeightSelect открывает меню выбора веса
func (h *FilterHandler) HandleWeightSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		"⚖️ **Вес (кг)**\n\nВыберите диапазон веса:",
	)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.WeightSelect()
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

// HandleWeightValue обрабатывает выбор веса
func (h *FilterHandler) HandleWeightValue(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if value == "any" {
		state.Filters.Weight = nil
	} else {
		weightRange := parseWeightRange(value)
		state.Filters.Weight = weightRange
	}

	h.stateManager.UpdateFilters(userID, state.Filters)
	return h.HandleFilterMenu(ctx, botAPI, query)
}
