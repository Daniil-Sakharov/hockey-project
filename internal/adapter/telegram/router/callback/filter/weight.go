package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleWeight обрабатывает фильтр веса
func HandleWeight(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 3 {
		return
	}

	subCommand := parts[2]

	if subCommand == cb.SubCmdSelect {
		if err := r.FilterHandler().HandleWeightSelect(ctx, bot, query); err != nil {
			logger.Error(ctx, "❌ Error handling weight select", zap.Error(err))
		}
	} else {
		// Значение веса (40-50, 90+, any)
		if err := r.FilterHandler().HandleWeightValue(ctx, bot, query, subCommand); err != nil {
			logger.Error(ctx, "❌ Error handling weight value",
				zap.Error(err),
				zap.String("value", subCommand))
		}
	}
}
