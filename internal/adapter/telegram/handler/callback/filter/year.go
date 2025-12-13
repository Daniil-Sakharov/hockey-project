package filter

import (
	"context"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleYearSelect открывает меню выбора года
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

// HandleYearValue обрабатывает выбор года
func (h *FilterHandler) HandleYearValue(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if value == "any" {
		// Сбрасываем фильтр года
		state.Filters.Year = nil
	} else {
		// Парсим год
		year, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid year value: %v", err)
			return err
		}
		state.Filters.Year = &year
	}

	// Обновляем состояние
	h.stateManager.UpdateFilters(userID, state.Filters)

	// Возвращаемся к меню фильтров
	return h.HandleFilterMenu(ctx, botAPI, query)
}
