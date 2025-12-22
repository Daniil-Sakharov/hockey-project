package callback

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleSearch обрабатывает callback поиска
func HandleSearch(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	cmd := parts[1]
	sh := r.SearchHandler()

	switch cmd {
	case cb.SearchPage:
		if len(parts) < 3 {
			return
		}
		switch parts[2] {
		case cb.PageNext:
			if err := sh.HandlePageNext(ctx, bot, query); err != nil {
				logger.Error(ctx, "Error handling page next", zap.Error(err))
			}
		case cb.PagePrev:
			if err := sh.HandlePagePrev(ctx, bot, query); err != nil {
				logger.Error(ctx, "Error handling page prev", zap.Error(err))
			}
		}
	case cb.SearchBackToFilters:
		if err := sh.HandleBackToFilters(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling back to filters", zap.Error(err))
		}
	case cb.SearchBackToResults:
		if err := sh.HandleBackToResults(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling back to results", zap.Error(err))
		}
	default:
		logger.Warn(ctx, "Unknown search command", zap.String("command", cmd))
	}
}
