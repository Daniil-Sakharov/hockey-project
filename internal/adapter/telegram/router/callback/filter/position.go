package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandlePosition обрабатывает фильтр позиции
func HandlePosition(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string, logger *zap.Logger) {
	if len(parts) < 3 {
		return
	}

	subCommand := parts[2]

	if subCommand == cb.SubCmdSelect {
		if err := r.FilterHandler().HandlePositionSelect(ctx, bot, query); err != nil {
			logger.Error("Error handling position select", zap.Error(err))
		}
	} else {
		// Значение позиции (forward, defender, goalie, any)
		if err := r.FilterHandler().HandlePositionValue(ctx, bot, query, subCommand); err != nil {
			logger.Error("Error handling position value",
				zap.Error(err),
				zap.String("value", subCommand))
		}
	}
}
