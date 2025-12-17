package profile

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Handler обрабатывает отображение профиля игрока
type Handler struct {
	keyboardPresenter *keyboard.KeyboardPresenter
	profilePresenter  *presenter.ProfilePresenter
	profileService    bot.ProfileService
	logger            *zap.Logger
}

// NewHandler создает новый Handler
func NewHandler(
	keyboardPresenter *keyboard.KeyboardPresenter,
	profilePresenter *presenter.ProfilePresenter,
	profileService bot.ProfileService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		keyboardPresenter: keyboardPresenter,
		profilePresenter:  profilePresenter,
		profileService:    profileService,
		logger:            logger,
	}
}

// HandleProfile отображает профиль игрока
func (h *Handler) HandleProfile(ctx context.Context, botAPI *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	// Извлекаем player_id из callback data
	// Формат: "player:profile:player_id"
	parts := strings.Split(query.Data, ":")
	if len(parts) != 3 {
		h.logger.Error("Invalid profile callback data",
			zap.String("data", query.Data),
			zap.Int("parts_count", len(parts)))
		return fmt.Errorf("invalid callback data format")
	}

	playerID := parts[2]

	// Получаем профиль игрока
	profile, err := h.profileService.GetPlayerProfile(ctx, playerID)
	if err != nil {
		h.logger.Error("Failed to get player profile",
			zap.String("player_id", playerID),
			zap.Error(err))

		// Отправляем сообщение об ошибке
		errorMsg := tgbotapi.NewMessage(query.Message.Chat.ID,
			"❌ Не удалось загрузить профиль игрока. Попробуйте позже.")
		_, _ = botAPI.Send(errorMsg)
		return err
	}

	// Форматируем профиль через presenter
	profileText, err := h.profilePresenter.FormatPlayerProfile(profile)
	if err != nil {
		h.logger.Error("Failed to format player profile",
			zap.String("player_id", playerID),
			zap.Error(err))

		errorMsg := tgbotapi.NewMessage(query.Message.Chat.ID,
			"❌ Ошибка форматирования профиля. Попробуйте позже.")
		_, _ = botAPI.Send(errorMsg)
		return err
	}

	// Создаем клавиатуру с кнопками
	keyboard := h.createProfileKeyboard(playerID)

	// Отправляем профиль
	msg := tgbotapi.NewMessage(query.Message.Chat.ID, profileText)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "" // Обычный текст, без HTML/Markdown

	_, err = botAPI.Send(msg)
	if err != nil {
		h.logger.Error("Failed to send profile message",
			zap.String("player_id", playerID),
			zap.Error(err))
		return err
	}

	// Отвечаем на callback чтобы убрать "часики"
	callback := tgbotapi.NewCallback(query.ID, "")
	_, _ = botAPI.Request(callback)

	h.logger.Info("Profile displayed successfully",
		zap.String("player_id", playerID),
		zap.Int64("user_id", query.From.ID))

	return nil
}

// createProfileKeyboard создает клавиатуру для профиля
func (h *Handler) createProfileKeyboard(playerID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 Скачать полный отчет", fmt.Sprintf("download_report:%s", playerID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀ Назад", "search:back_to_results"),
		),
	)
}
