package filter

import (
	"context"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *FilterHandler) HandleYearSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		"🎂 **Год рождения**\n\nВыберите год рождения игрока:",
	)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.YearSelect()
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

func (h *FilterHandler) HandleYearValue(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if value == "any" {
		state.Filters.Year = nil
	} else {
		year, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid year value: %v", err)
			return err
		}
		state.Filters.Year = &year
	}

	h.stateManager.UpdateFilters(userID, state.Filters)
	return h.HandleFilterMenu(ctx, botAPI, query)
}
