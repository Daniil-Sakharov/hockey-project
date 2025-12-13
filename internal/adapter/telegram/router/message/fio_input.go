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

// FioInputHandler обрабатывает ввод ФИО
type FioInputHandler struct {
	stateManager  bot.StateManager
	filterHandler *filter.FilterHandler
}

// NewFioInputHandler создает новый обработчик
func NewFioInputHandler(stateManager bot.StateManager, filterHandler *filter.FilterHandler) *FioInputHandler {
	return &FioInputHandler{
		stateManager:  stateManager,
		filterHandler: filterHandler,
	}
}

// HandleFioInput обрабатывает текстовый ввод ФИО
func (h *FioInputHandler) HandleFioInput(ctx context.Context, botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	userID := message.From.ID
	state := h.stateManager.GetState(userID)

	// Проверяем что мы ждем ввода ФИО
	if !strings.HasPrefix(state.WaitingForInput, "fio_") {
		return nil // Не наше
	}

	text := strings.TrimSpace(message.Text)

	// Валидация
	if !h.validateFioInput(text) {
		// Отправляем сообщение об ошибке
		errorMsg := tgbotapi.NewMessage(message.Chat.ID,
			"❌ Ошибка: введите минимум 2 символа кириллицей.\n\nПопробуйте еще раз:")
		if _, err := botAPI.Send(errorMsg); err != nil {
			log.Printf("Error sending validation error: %v", err)
		}
		return nil
	}

	// Сохраняем в TempFioFilters
	switch state.WaitingForInput {
	case "fio_last_name":
		state.TempFioFilters.LastName = text
	case "fio_first_name":
		state.TempFioFilters.FirstName = text
	case "fio_patronymic":
		state.TempFioFilters.Patronymic = text
	}

	// Сбрасываем ожидание ввода
	h.stateManager.SetWaitingForInput(userID, "")

	// Удаляем сообщение пользователя
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	if _, err := botAPI.Request(deleteMsg); err != nil {
		log.Printf("Error deleting user message: %v", err)
	}

	// Удаляем сообщение запроса
	deleteReq := tgbotapi.NewDeleteMessage(message.Chat.ID, state.LastMsgID)
	if _, err := botAPI.Request(deleteReq); err != nil {
		log.Printf("Error deleting request message: %v", err)
	}

	// Отправляем обновленную панель ФИО
	text = ""
	if state.TempFioFilters.LastName != "" {
		text += "Фамилия: " + state.TempFioFilters.LastName + "\n"
	}
	if state.TempFioFilters.FirstName != "" {
		text += "Имя: " + state.TempFioFilters.FirstName + "\n"
	}
	if state.TempFioFilters.Patronymic != "" {
		text += "Отчество: " + state.TempFioFilters.Patronymic + "\n"
	}

	// Рендерим панель ФИО через presenter
	fioText, _ := h.filterHandler.RenderFioMenuText(state.TempFioFilters)
	msg := tgbotapi.NewMessage(message.Chat.ID, fioText)
	msg.ParseMode = "Markdown"

	// Создаем клавиатуру через handler
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

// validateFioInput проверяет корректность ввода
func (h *FioInputHandler) validateFioInput(text string) bool {
	// Минимум 2 символа
	if len(text) < 2 {
		return false
	}

	// Только кириллица, пробелы и дефисы
	matched, _ := regexp.MatchString(`^[А-Яа-яЁё\s\-]+$`, text)
	return matched
}
