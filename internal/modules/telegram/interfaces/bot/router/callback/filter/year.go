package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleYear обрабатывает фильтр года
func HandleYear(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	subCmd := parts[2]

	if subCmd == cb.SubCmdSelect {
		if err := r.FilterHandler().HandleYearSelect(ctx, bot, query); err != nil {
			logger.Error(ctx, "Error handling year select", zap.Error(err))
		}
	} else {
		if err := r.FilterHandler().HandleYearValue(ctx, bot, query, subCmd); err != nil {
			logger.Error(ctx, "Error handling year value", zap.Error(err), zap.String("value", subCmd))
		}
	}
}
