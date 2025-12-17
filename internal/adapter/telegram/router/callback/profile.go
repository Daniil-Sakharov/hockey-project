package callback

import (
	"context"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleProfile обрабатывает callback для профиля игрока
func HandleProfile(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		logger.Warn(ctx, "❓ Invalid profile callback: not enough parts", zap.Strings("parts", parts))
		return
	}

	action := parts[1]

	switch action {
	case cb.PlayerProfile:
		// Обработка просмотра профиля
		if err := r.ProfileHandler().HandleProfile(ctx, bot, query); err != nil {
			logger.Error(ctx, "❌ Error handling profile",
				zap.Error(err),
				zap.String("callback_data", query.Data))
		}
	default:
		logger.Warn(ctx, "❓ Unknown player action", zap.String("action", action))
	}
}
