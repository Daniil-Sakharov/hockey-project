package callback

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// FilterHandler интерфейс
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
}

// SearchHandler интерфейс
type SearchHandler interface {
	HandleSearch(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandlePageNext(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandlePagePrev(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleBackToFilters(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	HandleBackToResults(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// ProfileHandler интерфейс
type ProfileHandler interface {
	HandleProfile(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// StartHandler интерфейс
type StartHandler interface {
	HandleMainMenuCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// ReportHandler интерфейс
type ReportHandler interface {
	HandleDownloadReport(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
}

// Router интерфейс для доступа к handlers
type Router interface {
	FilterHandler() FilterHandler
	SearchHandler() SearchHandler
	ProfileHandler() ProfileHandler
	StartHandler() StartHandler
	ReportHandler() ReportHandler
}

// HandleMenu обрабатывает callback главного меню
func HandleMenu(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	cmd := parts[1]

	switch cmd {
	case cb.MenuSearch:
		if err := r.FilterHandler().HandleFilterMenu(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling filter menu", zap.Error(err))
		}
	case cb.MenuMain:
		if err := r.StartHandler().HandleMainMenuCallback(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling main menu", zap.Error(err))
		}
	case cb.MenuStats:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "📊 Поиск по статистике - будет в следующей версии 🚧")
		_, _ = bot.Send(msg)
	case cb.MenuTeam:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "🏒 Поиск команды - будет в следующей версии 🚧")
		_, _ = bot.Send(msg)
	case cb.MenuHelp:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "❓ Помощь - в разработке")
		_, _ = bot.Send(msg)
	default:
		logger.Warn(ctx, "Unknown menu command", zap.String("command", cmd))
	}
}
