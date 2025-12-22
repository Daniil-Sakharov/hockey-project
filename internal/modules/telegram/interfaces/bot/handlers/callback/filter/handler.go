package filter

import (
	"context"
	"strconv"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/valueobjects"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter/keyboard"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает фильтры
type Handler struct {
	presenter    *presenter.Presenter
	keyboard     *keyboard.KeyboardPresenter
	stateService *services.UserStateService
}

// NewHandler создает новый Handler
func NewHandler(
	presenter *presenter.Presenter,
	keyboard *keyboard.KeyboardPresenter,
	stateService *services.UserStateService,
) *Handler {
	return &Handler{
		presenter:    presenter,
		keyboard:     keyboard,
		stateService: stateService,
	}
}

// HandleFilterMenu показывает меню фильтров
func (h *Handler) HandleFilterMenu(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)

	text, err := h.presenter.RenderFilterMenu(&session.Filters)
	if err != nil {
		return err
	}

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboard.FilterMenu(session.Filters.HasFilters())
	edit.ReplyMarkup = &markup

	_, err = bot.Send(edit)
	return err
}

// HandleFilterReset сбрасывает фильтры
func (h *Handler) HandleFilterReset(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	session.ResetFilters()
	return h.HandleFilterMenu(ctx, bot, query)
}

// HandleYearSelect показывает выбор года
func (h *Handler) HandleYearSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "🎂 **Год рождения**\n\nВыберите год рождения игрока:")
	edit.ParseMode = "Markdown"
	markup := h.keyboard.YearSelect()
	edit.ReplyMarkup = &markup
	_, err := bot.Send(edit)
	return err
}

// HandleYearValue устанавливает значение года
func (h *Handler) HandleYearValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	session := h.stateService.GetSession(query.From.ID)

	if value == "any" {
		session.Filters.Year = nil
	} else {
		year, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		session.Filters.Year = &year
	}

	return h.HandleFilterMenu(ctx, bot, query)
}

// Остальные методы (position, height, weight, region, fio) - упрощённые заглушки
func (h *Handler) HandlePositionSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "🏒 **Позиция**\n\nВыберите позицию игрока:")
	edit.ParseMode = "Markdown"
	markup := h.keyboard.PositionSelect()
	edit.ReplyMarkup = &markup
	_, err := bot.Send(edit)
	return err
}

func (h *Handler) HandlePositionValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	session := h.stateService.GetSession(query.From.ID)
	if value == "any" {
		session.Filters.Position = nil
	} else {
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
		session.Filters.Position = &position
	}
	return h.HandleFilterMenu(ctx, bot, query)
}

func (h *Handler) HandleHeightSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "📏 **Рост (см)**\n\nВыберите диапазон роста:")
	edit.ParseMode = "Markdown"
	markup := h.keyboard.HeightSelect()
	edit.ReplyMarkup = &markup
	_, err := bot.Send(edit)
	return err
}

func (h *Handler) HandleHeightValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	session := h.stateService.GetSession(query.From.ID)
	if value == "any" {
		session.Filters.Height = nil
	} else {
		session.Filters.Height = parseHeightRange(value)
	}
	return h.HandleFilterMenu(ctx, bot, query)
}

func (h *Handler) HandleWeightSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "⚖️ **Вес (кг)**\n\nВыберите диапазон веса:")
	edit.ParseMode = "Markdown"
	markup := h.keyboard.WeightSelect()
	edit.ReplyMarkup = &markup
	_, err := bot.Send(edit)
	return err
}

func (h *Handler) HandleWeightValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	session := h.stateService.GetSession(query.From.ID)
	if value == "any" {
		session.Filters.Weight = nil
	} else {
		session.Filters.Weight = parseWeightRange(value)
	}
	return h.HandleFilterMenu(ctx, bot, query)
}

func (h *Handler) HandleRegionSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "🗺️ **Регион (Федеральный округ)**\n\nВыберите федеральный округ:")
	edit.ParseMode = "Markdown"
	markup := h.keyboard.RegionSelect()
	edit.ReplyMarkup = &markup
	_, err := bot.Send(edit)
	return err
}

func (h *Handler) HandleRegionValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error {
	session := h.stateService.GetSession(query.From.ID)
	if value == "any" {
		session.Filters.Region = nil
	} else {
		session.Filters.Region = &value
	}
	return h.HandleFilterMenu(ctx, bot, query)
}

func (h *Handler) HandleFioSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)

	text, err := h.presenter.RenderFioMenu(session.TempFIO)
	if err != nil {
		return err
	}

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboard.FioMenu(session.TempFIO)
	edit.ReplyMarkup = &markup

	_, err = bot.Send(edit)
	return err
}

func (h *Handler) HandleFioField(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error {
	h.stateService.GetSession(query.From.ID).WaitingForInput = "fio_" + field

	text := h.presenter.RenderFioInputRequest(field)
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	markup := h.keyboard.FioCancelButton()
	msg.ReplyMarkup = markup

	_, err := bot.Send(msg)
	return err
}

func (h *Handler) HandleFioClear(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error {
	session := h.stateService.GetSession(query.From.ID)

	switch field {
	case "clear_last":
		session.TempFIO.LastName = ""
	case "clear_first":
		session.TempFIO.FirstName = ""
	case "clear_patr":
		session.TempFIO.Patronymic = ""
	}

	return h.HandleFioSelect(ctx, bot, query)
}

func (h *Handler) HandleFioApply(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	session.ApplyTempFIO()
	return h.HandleFilterMenu(ctx, bot, query)
}

func (h *Handler) HandleFioBack(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	session := h.stateService.GetSession(query.From.ID)
	session.ClearTempFIO()
	return h.HandleFilterMenu(ctx, bot, query)
}

func (h *Handler) HandleTextInput(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	session := h.stateService.GetSession(msg.From.ID)

	if !session.IsWaitingForInput() {
		return nil
	}

	switch session.WaitingForInput {
	case "fio_last_name":
		session.TempFIO.LastName = msg.Text
	case "fio_first_name":
		session.TempFIO.FirstName = msg.Text
	case "fio_patronymic":
		session.TempFIO.Patronymic = msg.Text
	}

	session.WaitingForInput = ""

	text, err := h.presenter.RenderFioMenu(session.TempFIO)
	if err != nil {
		return err
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	reply.ReplyMarkup = h.keyboard.FioMenu(session.TempFIO)

	_, err = bot.Send(reply)
	return err
}

// Вспомогательные функции
func parseHeightRange(value string) *valueobjects.Range {
	switch value {
	case "150-160":
		return &valueobjects.Range{Min: 150, Max: 160}
	case "160-170":
		return &valueobjects.Range{Min: 160, Max: 170}
	case "170-180":
		return &valueobjects.Range{Min: 170, Max: 180}
	case "180-190":
		return &valueobjects.Range{Min: 180, Max: 190}
	case "190-200":
		return &valueobjects.Range{Min: 190, Max: 200}
	case "200-250":
		return &valueobjects.Range{Min: 200, Max: 250}
	}
	return nil
}

func parseWeightRange(value string) *valueobjects.Range {
	switch value {
	case "40-50":
		return &valueobjects.Range{Min: 40, Max: 50}
	case "50-60":
		return &valueobjects.Range{Min: 50, Max: 60}
	case "60-70":
		return &valueobjects.Range{Min: 60, Max: 70}
	case "70-80":
		return &valueobjects.Range{Min: 70, Max: 80}
	case "80-90":
		return &valueobjects.Range{Min: 80, Max: 90}
	case "90-150":
		return &valueobjects.Range{Min: 90, Max: 150}
	}
	return nil
}
