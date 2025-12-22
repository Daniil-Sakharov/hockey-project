package handlers

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHandler интерфейс для start команды
type StartHandler interface {
	HandleStart(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error
	HandleMainMenuCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// FilterHandler интерфейс для фильтров
type FilterHandler interface {
	HandleFilterMenu(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleFilterReset(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleYearSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleYearValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error
	HandlePositionSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandlePositionValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error
	HandleHeightSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleHeightValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error
	HandleWeightSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleWeightValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error
	HandleRegionSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleRegionValue(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, value string) error
	HandleFioSelect(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleFioField(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error
	HandleFioClear(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, field string) error
	HandleFioApply(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleFioBack(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleTextInput(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error
}

// SearchHandler интерфейс для поиска
type SearchHandler interface {
	HandleSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandlePageNext(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandlePagePrev(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleBackToFilters(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleBackToResults(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// ProfileHandler интерфейс для профиля
type ProfileHandler interface {
	HandleProfile(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// ReportHandler интерфейс для отчётов
type ReportHandler interface {
	HandleDownloadReport(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}
