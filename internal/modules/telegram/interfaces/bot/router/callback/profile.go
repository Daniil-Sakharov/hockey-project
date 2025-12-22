package callback

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleProfile обрабатывает callback профиля
func HandleProfile(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, _ []string) {
	if err := r.ProfileHandler().HandleProfile(ctx, bot, query); err != nil {
		logger.Error(ctx, "Error handling profile", zap.Error(err))
	}
}
