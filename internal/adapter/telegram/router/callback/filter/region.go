package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleRegion обрабатывает фильтр региона
func HandleRegion(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string, logger *zap.Logger) {
	if len(parts) < 3 {
		return
	}

	subCommand := parts[2]

	if subCommand == cb.SubCmdSelect {
		if err := r.FilterHandler().HandleRegionSelect(ctx, bot, query); err != nil {
			logger.Error("Error handling region select", zap.Error(err))
		}
	} else {
		// Значение региона (ЦФО, СФО, УФО, any)
		if err := r.FilterHandler().HandleRegionValue(ctx, bot, query, subCommand); err != nil {
			logger.Error("Error handling region value",
				zap.Error(err),
				zap.String("value", subCommand))
		}
	}
}
