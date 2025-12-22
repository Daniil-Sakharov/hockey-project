package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleFio обрабатывает фильтр ФИО
func HandleFio(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	subCmd := parts[2]

	switch subCmd {
	case cb.FioSelect:
		if err := r.FilterHandler().HandleFioSelect(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling FIO select", zap.Error(err))
		}
	case cb.FioLastName, cb.FioFirstName, cb.FioPatronymic:
		if err := r.FilterHandler().HandleFioField(ctx, bot, query, subCmd); err != nil {
			logger.Error(ctx, "Error handling FIO field", zap.Error(err), zap.String("field", subCmd))
		}
	case cb.FioClearLast, cb.FioClearFirst, cb.FioClearPatr:
		if err := r.FilterHandler().HandleFioClear(ctx, bot, query, subCmd); err != nil {
			logger.Error(ctx, "Error handling FIO clear", zap.Error(err), zap.String("field", subCmd))
		}
	case cb.FioApply:
		if err := r.FilterHandler().HandleFioApply(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling FIO apply", zap.Error(err))
		}
	case cb.FioBack:
		if err := r.FilterHandler().HandleFioBack(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling FIO back", zap.Error(err))
		}
	default:
		logger.Warn(ctx, "Unknown FIO command", zap.String("command", subCmd))
	}
}
