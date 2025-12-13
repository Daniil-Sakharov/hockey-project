package filter

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandlePositionSelect открывает меню выбора позиции
func (h *FilterHandler) HandlePositionSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		"🏒 **Позиция**\n\nВыберите позицию игрока:",
	)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.PositionSelect()
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

// HandlePositionValue обрабатывает выбор позиции
func (h *FilterHandler) HandlePositionValue(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if value == "any" {
		state.Filters.Position = nil
	} else {
		// Преобразуем значение в русское название
		var position string
		switch value {
		case "forward":
			position = "Нападающий"
		case "defender":
			position = "Защитник"
		case "goalie":
			position = "Вратарь"
		default:
			position = value
		}
		state.Filters.Position = &position
	}

	h.stateManager.UpdateFilters(userID, state.Filters)
	return h.HandleFilterMenu(ctx, botAPI, query)
}
