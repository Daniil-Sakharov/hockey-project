package callback

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/report"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// ReportRouter интерфейс для доступа к report handler
type ReportRouter interface {
	ReportHandler() *report.Handler
}

// HandleReport обрабатывает callback для скачивания отчета
func HandleReport(r ReportRouter, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, logger *zap.Logger) {
	if err := r.ReportHandler().HandleDownloadReport(ctx, bot, query); err != nil {
		logger.Error("Error handling report download",
			zap.Error(err),
			zap.String("callback_data", query.Data))
	}
}
