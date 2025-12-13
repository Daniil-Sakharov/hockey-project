package callback

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router/callback/filter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// routerAdapter адаптер для преобразования callback.Router в filter.Router
type routerAdapter struct {
	router Router
}

func (a *routerAdapter) FilterHandler() interface {
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
} {
	return a.router.FilterHandler()
}

func (a *routerAdapter) SearchHandler() interface {
	HandleSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
} {
	return a.router.SearchHandler()
}

// HandleFilter обрабатывает callback фильтров
func HandleFilter(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	// Делегируем обработку в пакет filter через адаптер
	adapter := &routerAdapter{router: r}
	filter.Handle(adapter, ctx, bot, query, parts)
}
