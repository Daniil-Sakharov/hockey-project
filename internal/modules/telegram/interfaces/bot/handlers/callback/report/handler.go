package report

import (
	"context"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает отчёты
type Handler struct {
	reportService *services.ReportService
}

// NewHandler создает новый Handler
func NewHandler(reportService *services.ReportService) *Handler {
	return &Handler{reportService: reportService}
}

// HandleDownloadReport генерирует и отправляет отчёт
func (h *Handler) HandleDownloadReport(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	parts := strings.Split(query.Data, ":")
	if len(parts) != 2 {
		return nil
	}

	playerID := parts[1]

	htmlData, filename, err := h.reportService.GenerateReport(ctx, playerID)
	if err != nil {
		return err
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: htmlData,
	}

	doc := tgbotapi.NewDocument(query.Message.Chat.ID, fileBytes)
	doc.Caption = "📊 Полный отчет игрока\n\nОткройте файл в браузере для просмотра графиков."

	_, err = bot.Send(doc)
	return err
}
