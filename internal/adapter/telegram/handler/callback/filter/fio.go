package filter

import (
	"context"
	"log"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *FilterHandler) HandleFioSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if state.Filters.FirstName != nil {
		state.TempFioFilters.FirstName = *state.Filters.FirstName
	} else {
		state.TempFioFilters.FirstName = ""
	}

	if state.Filters.LastName != nil {
		state.TempFioFilters.LastName = *state.Filters.LastName
	} else {
		state.TempFioFilters.LastName = ""
	}

	state.TempFioFilters.Patronymic = ""

	h.stateManager.SetCurrentView(userID, "fio_menu")

	text, err := h.msgPresenter.RenderFioMenu(state.TempFioFilters)
	if err != nil {
		log.Printf("Error rendering FIO menu: %v", err)
		return err
	}

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.FioMenu(state.TempFioFilters)
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

func (h *FilterHandler) HandleFioField(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error {
	userID := query.From.ID

	waitingFor := "fio_" + field
	h.stateManager.SetWaitingForInput(userID, waitingFor)

	text := h.msgPresenter.RenderFioInputRequest(field)

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	markup := h.keyboardPresenter.FioCancelButton()
	msg.ReplyMarkup = markup

	sent, err := botAPI.Send(msg)
	if err != nil {
		log.Printf("Error sending input request: %v", err)
		return err
	}

	h.stateManager.SetLastMsgID(userID, sent.MessageID)
	return nil
}

func (h *FilterHandler) HandleFioClear(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	switch field {
	case cb.FioClearLast:
		state.TempFioFilters.LastName = ""
	case cb.FioClearFirst:
		state.TempFioFilters.FirstName = ""
	case cb.FioClearPatr:
		state.TempFioFilters.Patronymic = ""
	}

	return h.HandleFioSelect(ctx, botAPI, query)
}

func (h *FilterHandler) HandleFioApply(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	lastName := state.TempFioFilters.LastName
	firstName := state.TempFioFilters.FirstName

	if lastName != "" {
		state.Filters.LastName = &lastName
	} else {
		state.Filters.LastName = nil
	}

	if firstName != "" {
		state.Filters.FirstName = &firstName
	} else {
		state.Filters.FirstName = nil
	}

	state.TempFioFilters = domainBot.TempFioData{}

	return h.HandleFilterMenu(ctx, botAPI, query)
}

func (h *FilterHandler) HandleFioBack(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	state.TempFioFilters = domainBot.TempFioData{}

	return h.HandleFilterMenu(ctx, botAPI, query)
}
