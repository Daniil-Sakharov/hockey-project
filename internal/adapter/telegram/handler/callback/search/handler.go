package search

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает поиск игроков
type Handler struct {
	msgPresenter      *message.MessagePresenter
	keyboardPresenter *keyboard.KeyboardPresenter
	stateManager      bot.StateManager
	searchService     bot.SearchPlayerService
}

// NewHandler создает новый Handler
func NewHandler(
	msgPresenter *message.MessagePresenter,
	keyboardPresenter *keyboard.KeyboardPresenter,
	stateManager bot.StateManager,
	searchService bot.SearchPlayerService,
) *Handler {
	return &Handler{
		msgPresenter:      msgPresenter,
		keyboardPresenter: keyboardPresenter,
		stateManager:      stateManager,
		searchService:     searchService,
	}
}

// HandleSearch выполняет поиск и показывает результаты
func (h *Handler) HandleSearch(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	// Сбрасываем на первую страницу
	state.CurrentPage = 1

	// Удаляем сообщение с меню фильтров
	deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	if _, err := botAPI.Request(deleteMsg); err != nil {
		log.Printf("Error deleting filter menu message: %v", err)
	}

	return h.showSearchResults(ctx, botAPI, query.Message.Chat.ID, userID)
}

// HandlePageNext следующая страница
func (h *Handler) HandlePageNext(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)
	state.CurrentPage++

	return h.showSearchResults(ctx, botAPI, query.Message.Chat.ID, userID)
}

// HandlePagePrev предыдущая страница
func (h *Handler) HandlePagePrev(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	if state.CurrentPage > 1 {
		state.CurrentPage--
	}

	return h.showSearchResults(ctx, botAPI, query.Message.Chat.ID, userID)
}

// HandleCleanup удаляет результаты поиска (используется при возврате в меню фильтров)
func (h *Handler) HandleCleanup(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	// Удаляем все сообщения результатов поиска
	h.deleteOldResults(botAPI, query.Message.Chat.ID, state)

	return nil
}

// HandleBackToResults возвращает к результатам поиска из профиля игрока
func (h *Handler) HandleBackToResults(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID

	// Удаляем сообщение с профилем игрока
	deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
	if _, err := botAPI.Request(deleteMsg); err != nil {
		log.Printf("Error deleting profile message: %v", err)
	}

	// Показываем результаты поиска с текущими фильтрами
	return h.showSearchResults(ctx, botAPI, query.Message.Chat.ID, userID)
}
