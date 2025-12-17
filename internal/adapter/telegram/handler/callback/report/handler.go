package report

import (
	"context"
	"strings"

	reportService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/report"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Handler обрабатывает скачивание отчета
type Handler struct {
	reportService *reportService.Service
	logger        *zap.Logger
}

// NewHandler создает новый Handler
func NewHandler(
	reportService *reportService.Service,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		reportService: reportService,
		logger:        logger,
	}
}

// HandleDownloadReport генерирует и отправляет HTML отчет игрока
func (h *Handler) HandleDownloadReport(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	// Извлекаем player_id из callback data
	// Формат: "download_report:player_id"
	parts := strings.Split(query.Data, ":")
	if len(parts) != 2 {
		h.logger.Error("Invalid report callback data",
			zap.String("data", query.Data),
			zap.Int("parts_count", len(parts)))

		callback := tgbotapi.NewCallback(query.ID, "❌ Ошибка формата данных")
		_, _ = botAPI.Request(callback)
		return nil
	}

	playerID := parts[1]

	// Отправляем сообщение о начале генерации
	waitMsg := tgbotapi.NewMessage(query.Message.Chat.ID, "⏳ Генерирую полный отчет игрока, это может занять несколько секунд...")
	sentMsg, err := botAPI.Send(waitMsg)
	if err != nil {
		h.logger.Error("Failed to send wait message", zap.Error(err))
	}

	// Генерируем отчет
	htmlData, filename, err := h.reportService.GeneratePlayerReport(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to generate player report",
			zap.String("player_id", playerID),
			zap.Error(err))

		errorMsg := tgbotapi.NewMessage(query.Message.Chat.ID,
			"❌ Не удалось сгенерировать отчет. Попробуйте позже.")
		_, _ = botAPI.Send(errorMsg)

		// Удаляем сообщение ожидания
		if sentMsg.MessageID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, sentMsg.MessageID)
			_, _ = botAPI.Request(deleteMsg)
		}

		return err
	}

	// Удаляем сообщение ожидания
	if sentMsg.MessageID != 0 {
		deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, sentMsg.MessageID)
		_, _ = botAPI.Request(deleteMsg)
	}

	// Отправляем файл
	fileBytes := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: htmlData,
	}

	doc := tgbotapi.NewDocument(query.Message.Chat.ID, fileBytes)
	doc.Caption = "📊 Полный отчет игрока\n\nОткройте файл в браузере для просмотра графиков."

	_, err = botAPI.Send(doc)
	if err != nil {
		h.logger.Error("Failed to send report file",
			zap.String("player_id", playerID),
			zap.Int("file_size", len(htmlData)),
			zap.Error(err))

		// Если не удалось отправить файл, попробуем отправить как текст
		if len(htmlData) < 4096 {
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, string(htmlData))
			_, _ = botAPI.Send(msg)
		} else {
			errorMsg := tgbotapi.NewMessage(query.Message.Chat.ID,
				"❌ Не удалось отправить файл отчета. Попробуйте позже.")
			_, _ = botAPI.Send(errorMsg)
		}
		return err
	}

	h.logger.Info("Report sent successfully",
		zap.String("player_id", playerID),
		zap.String("filename", filename),
		zap.Int("file_size", len(htmlData)),
		zap.Int64("user_id", query.From.ID))

	return nil
}
