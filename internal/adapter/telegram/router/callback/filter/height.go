package filter

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleHeight обрабатывает фильтр роста
func HandleHeight(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string, logger *zap.Logger) {
	if len(parts) < 3 {
		return
	}

	subCommand := parts[2]

	if subCommand == cb.SubCmdSelect {
		if err := r.FilterHandler().HandleHeightSelect(ctx, bot, query); err != nil {
			logger.Error("Error handling height select", zap.Error(err))
		}
	} else {
		// Значение роста (150-160, 200+, any)
		if err := r.FilterHandler().HandleHeightValue(ctx, bot, query, subCommand); err != nil {
			logger.Error("Error handling height value",
				zap.Error(err),
				zap.String("value", subCommand))
		}
	}
}
