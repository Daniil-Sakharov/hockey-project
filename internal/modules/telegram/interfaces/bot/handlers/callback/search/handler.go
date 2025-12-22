package search

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает поиск
type Handler struct {
	presenter     *presenter.Presenter
	keyboard      *keyboard.KeyboardPresenter
	stateService  *services.UserStateService
	searchService *services.PlayerSearchService
}

// NewHandler создает новый Handler
func NewHandler(
	presenter *presenter.Presenter,
	keyboard *keyboard.KeyboardPresenter,
	stateService *services.UserStateService,
	searchService *services.PlayerSearchService,
) *Handler {
	return &Handler{
		presenter:     presenter,
		keyboard:      keyboard,
		stateService:  stateService,
		searchService: searchService,
	}
}

// HandleSearch выполняет поиск
func (h *Handler) HandleSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	session.CurrentPage = 1
	return h.doSearch(ctx, bot, query, session)
}

// doSearch внутренний метод поиска
func (h *Handler) doSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, session *entities.UserSession) error {
	result, err := h.searchService.Search(ctx, session.Filters, session.CurrentPage, 5)
	if err != nil {
		return err
	}

	text := h.presenter.RenderSearchResults(result)

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboard.SearchResultsKeyboard(result.Players, result.CurrentPage, result.TotalPages)
	edit.ReplyMarkup = &markup

	_, err = bot.Send(edit)
	return err
}

// HandlePageNext следующая страница
func (h *Handler) HandlePageNext(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	session.NextPage()
	return h.doSearch(ctx, bot, query, session)
}

// HandlePagePrev предыдущая страница
func (h *Handler) HandlePagePrev(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	session.PrevPage()
	return h.doSearch(ctx, bot, query, session)
}

// HandleBackToFilters возврат к фильтрам
func (h *Handler) HandleBackToFilters(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)

	text := h.presenter.RenderFilterMenuFromSession(session)

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboard.FilterMenu(session.Filters.HasFilters())
	edit.ReplyMarkup = &markup

	_, err := bot.Send(edit)
	return err
}

// HandleBackToResults возврат к результатам
func (h *Handler) HandleBackToResults(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	return h.doSearch(ctx, bot, query, session)
}
