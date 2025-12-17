package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Router интерфейс для доступа к handlers (соответствует callback.Router)
type Router interface {
	FilterHandler() interface {
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
	}
	SearchHandler() interface {
		HandleSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	}
}

// Handle обрабатывает callback фильтров (главный роутер)
func Handle(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	filterType := parts[1]

	// Создаем zap logger для обработки callback
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	switch filterType {
	case cb.FilterBack:
		// Возврат к меню фильтров
		if err := r.FilterHandler().HandleFilterMenu(ctx, bot, query); err != nil {
			logger.Error("Error handling filter menu", zap.Error(err))
		}
	case cb.FilterReset:
		if err := r.FilterHandler().HandleFilterReset(ctx, bot, query); err != nil {
			logger.Error("Error handling filter reset", zap.Error(err))
		}
	case cb.FilterApply:
		// Перенаправляем на search handler
		if err := r.SearchHandler().HandleSearch(ctx, bot, query); err != nil {
			logger.Error("Error handling search", zap.Error(err))
		}
	case cb.FilterYear:
		HandleYear(r, ctx, bot, query, parts, logger)
	case cb.FilterPosition:
		HandlePosition(r, ctx, bot, query, parts, logger)
	case cb.FilterHeight:
		HandleHeight(r, ctx, bot, query, parts, logger)
	case cb.FilterWeight:
		HandleWeight(r, ctx, bot, query, parts, logger)
	case cb.FilterRegion:
		HandleRegion(r, ctx, bot, query, parts, logger)
	case cb.FilterFio:
		HandleFio(r, ctx, bot, query, parts, logger)
	default:
		logger.Warn("Unknown filter type", zap.String("filter_type", filterType))
	}
}
