package callback

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// HandleReport обрабатывает callback отчёта
func HandleReport(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	if err := r.ReportHandler().HandleDownloadReport(ctx, bot, query); err != nil {
		logger.Error(ctx, "Error handling report download", zap.Error(err))
	}
}
