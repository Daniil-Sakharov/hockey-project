package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleYear обрабатывает фильтр года
func HandleYear(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string, logger *zap.Logger) {
	if len(parts) < 3 {
		return
	}

	subCommand := parts[2]

	if subCommand == cb.SubCmdSelect {
		if err := r.FilterHandler().HandleYearSelect(ctx, bot, query); err != nil {
			logger.Error("Error handling year select", zap.Error(err))
		}
	} else {
		// Значение года (2005, 2006, any, etc.)
		if err := r.FilterHandler().HandleYearValue(ctx, bot, query, subCommand); err != nil {
			logger.Error("Error handling year value",
				zap.Error(err),
				zap.String("value", subCommand))
		}
	}
}
