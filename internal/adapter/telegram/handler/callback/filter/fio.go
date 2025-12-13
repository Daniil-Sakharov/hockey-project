package filter

import (
	"context"
	"log"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleFioSelect открывает панель ввода ФИО
func (h *FilterHandler) HandleFioSelect(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	// Инициализируем TempFioFilters текущими значениями из filters
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

	// Отчество пока не используется в domain, но готовим
	state.TempFioFilters.Patronymic = ""

	h.stateManager.SetCurrentView(userID, "fio_menu")

	// Рендерим панель ФИО
	text, err := h.msgPresenter.RenderFioMenu(state.TempFioFilters)
	if err != nil {
		log.Printf("Error rendering FIO menu: %v", err)
		return err
	}

	// Обновляем сообщение
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

// HandleFioField запускает ожидание ввода поля (фамилия/имя/отчество)
func (h *FilterHandler) HandleFioField(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error {
	userID := query.From.ID

	// Устанавливаем состояние ожидания
	waitingFor := "fio_" + field
	h.stateManager.SetWaitingForInput(userID, waitingFor)

	// Рендерим сообщение запроса
	text := h.msgPresenter.RenderFioInputRequest(field)

	// Отправляем новое сообщение с кнопкой отмены
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	markup := h.keyboardPresenter.FioCancelButton()
	msg.ReplyMarkup = markup

	sent, err := botAPI.Send(msg)
	if err != nil {
		log.Printf("Error sending input request: %v", err)
		return err
	}

	// Сохраняем ID сообщения запроса
	h.stateManager.SetLastMsgID(userID, sent.MessageID)

	return nil
}

// HandleFioClear очищает поле ФИО
func (h *FilterHandler) HandleFioClear(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	// Очищаем соответствующее поле
	switch field {
	case cb.FioClearLast:
		state.TempFioFilters.LastName = ""
	case cb.FioClearFirst:
		state.TempFioFilters.FirstName = ""
	case cb.FioClearPatr:
		state.TempFioFilters.Patronymic = ""
	}

	// Возвращаемся на панель ФИО
	return h.HandleFioSelect(ctx, botAPI, query)
}

// HandleFioApply применяет изменения и возвращается в главное меню фильтров
func (h *FilterHandler) HandleFioApply(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	// Создаём копии строк перед сохранением (важно! иначе указатели будут на обнулённые поля)
	lastName := state.TempFioFilters.LastName
	firstName := state.TempFioFilters.FirstName

	// Копируем в Filters
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

	// Очищаем временное хранилище
	state.TempFioFilters = domainBot.TempFioData{}

	// Возвращаемся в главное меню фильтров
	return h.HandleFilterMenu(ctx, botAPI, query)
}

// HandleFioBack отменяет изменения и возвращается в главное меню фильтров
func (h *FilterHandler) HandleFioBack(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	// Очищаем временное хранилище без применения
	state.TempFioFilters = domainBot.TempFioData{}

	// Возвращаемся в главное меню фильтров
	return h.HandleFilterMenu(ctx, botAPI, query)
}
