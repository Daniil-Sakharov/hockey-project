package filter

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleHeightSelect открывает меню выбора роста
func (h *FilterHandler) HandleHeightSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		"📏 **Рост (см)**\n\nВыберите диапазон роста:",
	)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.HeightSelect()
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

// HandleHeightValue обрабатывает выбор роста
func (h *FilterHandler) HandleHeightValue(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if value == "any" {
		state.Filters.Height = nil
	} else {
		heightRange := parseHeightRange(value)
		state.Filters.Height = heightRange
	}

	h.stateManager.UpdateFilters(userID, state.Filters)
	return h.HandleFilterMenu(ctx, botAPI, query)
}
