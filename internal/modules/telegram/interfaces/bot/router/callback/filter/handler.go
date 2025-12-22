package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Handler интерфейс filter handler
type Handler interface {
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

// SearchHandler интерфейс
type SearchHandler interface {
	HandleSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// Router интерфейс
type Router interface {
	FilterHandler() Handler
	SearchHandler() SearchHandler
}

// Handle главный роутер фильтров
func Handle(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	filterType := parts[1]

	logger.Info(ctx, "Handling filter callback",
		zap.String("filter_type", filterType),
		zap.Int64("user_id", query.From.ID))

	switch filterType {
	case cb.FilterBack:
		if err := r.FilterHandler().HandleFilterMenu(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling filter menu", zap.Error(err))
		}
	case cb.FilterReset:
		if err := r.FilterHandler().HandleFilterReset(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling filter reset", zap.Error(err))
		}
	case cb.FilterApply:
		if err := r.SearchHandler().HandleSearch(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling search", zap.Error(err))
		}
	case cb.FilterYear:
		HandleYear(r, ctx, bot, query, parts)
	case cb.FilterPosition:
		HandlePosition(r, ctx, bot, query, parts)
	case cb.FilterHeight:
		HandleHeight(r, ctx, bot, query, parts)
	case cb.FilterWeight:
		HandleWeight(r, ctx, bot, query, parts)
	case cb.FilterRegion:
		HandleRegion(r, ctx, bot, query, parts)
	case cb.FilterFio:
		HandleFio(r, ctx, bot, query, parts)
	default:
		logger.Warn(ctx, "Unknown filter type", zap.String("filter_type", filterType))
	}
}
