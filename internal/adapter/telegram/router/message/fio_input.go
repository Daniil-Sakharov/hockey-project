package message

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/filter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FioInputHandler struct {
	stateManager  bot.StateManager
	filterHandler *filter.FilterHandler
}

func NewFioInputHandler(stateManager bot.StateManager, filterHandler *filter.FilterHandler) *FioInputHandler {
	return &FioInputHandler{
		stateManager:  stateManager,
		filterHandler: filterHandler,
	}
}

func (h *FioInputHandler) HandleFioInput(ctx context.Context, botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	userID := message.From.ID
	state := h.stateManager.GetState(userID)

	if !strings.HasPrefix(state.WaitingForInput, "fio_") {
		return nil
	}

	text := strings.TrimSpace(message.Text)

	if !h.validateFioInput(text) {
		errorMsg := tgbotapi.NewMessage(message.Chat.ID,
			"❌ Ошибка: введите минимум 2 символа кириллицей.\n\nПопробуйте еще раз:")
		if _, err := botAPI.Send(errorMsg); err != nil {
			log.Printf("Error sending validation error: %v", err)
		}
		return nil
	}

	switch state.WaitingForInput {
	case "fio_last_name":
		state.TempFioFilters.LastName = text
	case "fio_first_name":
		state.TempFioFilters.FirstName = text
	case "fio_patronymic":
		state.TempFioFilters.Patronymic = text
	}

	h.stateManager.SetWaitingForInput(userID, "")

	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	if _, err := botAPI.Request(deleteMsg); err != nil {
		log.Printf("Error deleting user message: %v", err)
	}

	deleteReq := tgbotapi.NewDeleteMessage(message.Chat.ID, state.LastMsgID)
	if _, err := botAPI.Request(deleteReq); err != nil {
		log.Printf("Error deleting request message: %v", err)
	}

	fioText, _ := h.filterHandler.RenderFioMenuText(state.TempFioFilters)
	msg := tgbotapi.NewMessage(message.Chat.ID, fioText)
	msg.ParseMode = "Markdown"

	markup := h.filterHandler.GetFioKeyboard(state.TempFioFilters)
	msg.ReplyMarkup = markup

	sent, err := botAPI.Send(msg)
	if err != nil {
		log.Printf("Error sending FIO menu: %v", err)
		return err
	}

	h.stateManager.SetLastMsgID(userID, sent.MessageID)
	return nil
}

func (h *FioInputHandler) validateFioInput(text string) bool {
	if len(text) < 2 {
		return false
	}
	matched, _ := regexp.MatchString(`^[А-Яа-яЁё\s\-]+$`, text)
	return matched
}
