package filter

import (
	"context"
	"fmt"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FilterHandler struct {
	msgPresenter      *message.MessagePresenter
	keyboardPresenter *keyboard.KeyboardPresenter
	stateManager      bot.StateManager
}

func NewFilterHandler(
	msgPresenter *message.MessagePresenter,
	keyboardPresenter *keyboard.KeyboardPresenter,
	stateManager bot.StateManager,
) *FilterHandler {
	return &FilterHandler{
		msgPresenter:      msgPresenter,
		keyboardPresenter: keyboardPresenter,
		stateManager:      stateManager,
	}
}

func (h *FilterHandler) HandleFilterMenu(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	h.stateManager.SetLastMsgID(userID, query.Message.MessageID)
	h.stateManager.SetCurrentView(userID, "filter_menu")

	text, err := h.renderFilterMenu(state.Filters)
	if err != nil {
		log.Printf("Error rendering filter menu: %v", err)
		return err
	}

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.FilterMenu(state.Filters.HasFilters())
	edit.ReplyMarkup = &markup

	if _, err := botAPI.Send(edit); err != nil {
		log.Printf("Error editing message: %v", err)
		return err
	}

	return nil
}

func (h *FilterHandler) HandleFilterMenuNew(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	state := h.stateManager.GetState(userID)

	h.stateManager.SetCurrentView(userID, "filter_menu")

	text, err := h.renderFilterMenu(state.Filters)
	if err != nil {
		log.Printf("Error rendering filter menu: %v", err)
		return err
	}

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	markup := h.keyboardPresenter.FilterMenu(state.Filters.HasFilters())
	msg.ReplyMarkup = markup

	sent, err := botAPI.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	h.stateManager.SetLastMsgID(userID, sent.MessageID)
	return nil
}

func (h *FilterHandler) HandleFilterReset(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userID := query.From.ID
	h.stateManager.ResetFilters(userID)
	return h.HandleFilterMenu(ctx, botAPI, query)
}

func (h *FilterHandler) HandleFilterApply(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	return nil
}

func (h *FilterHandler) renderFilterMenu(filters domainBot.SearchFilters) (string, error) {
	data := message.FilterMenuData{
		FIO:                filters.GetFIODisplay(),
		Year:               filters.Year,
		Position:           filters.Position,
		Region:             filters.Region,
		HasFilters:         filters.HasFilters(),
		ActiveFiltersCount: filters.CountActiveFilters(),
	}

	if filters.Height != nil {
		data.Height = fmt.Sprintf("%d-%d", filters.Height.Min, filters.Height.Max)
	}

	if filters.Weight != nil {
		data.Weight = fmt.Sprintf("%d-%d", filters.Weight.Min, filters.Weight.Max)
	}

	return h.msgPresenter.RenderFilterMenu(data)
}

func (h *FilterHandler) RenderFioMenuText(fio domainBot.TempFioData) (string, error) {
	return h.msgPresenter.RenderFioMenu(fio)
}

func (h *FilterHandler) GetFioKeyboard(fio domainBot.TempFioData) tgbotapi.InlineKeyboardMarkup {
	return h.keyboardPresenter.FioMenu(fio)
}
