package command

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// StartHandler обрабатывает команду /start
type StartHandler struct {
	msgPresenter      *message.MessagePresenter
	keyboardPresenter *keyboard.KeyboardPresenter
}

// NewStartHandler создает новый StartHandler
func NewStartHandler(
	msgPresenter *message.MessagePresenter,
	keyboardPresenter *keyboard.KeyboardPresenter,
) *StartHandler {
	return &StartHandler{
		msgPresenter:      msgPresenter,
		keyboardPresenter: keyboardPresenter,
	}
}

// Handle обрабатывает команду /start
func (h *StartHandler) Handle(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	logger.Info(ctx, "🏁 Handling /start command",
		zap.Int64("user_id", msg.From.ID),
		zap.String("username", msg.From.UserName))

	// Получаем имя пользователя
	userName := msg.From.FirstName
	if userName == "" {
		userName = msg.From.UserName
	}

	// Подготавливаем данные для шаблона
	data := message.WelcomeData{
		UserName:         userName,
		PlayersCount:     14545, // TODO: получать из БД
		TeamsCount:       613,   // TODO: получать из БД
		TournamentsCount: 36,    // TODO: получать из БД
	}

	// Рендерим сообщение
	text, err := h.msgPresenter.RenderWelcome(data)
	if err != nil {
		logger.Error(ctx, "❌ Failed to render welcome message", zap.Error(err))
		text = "Добро пожаловать в Hockey Scout Bot!"
	}

	// Создаем сообщение с inline клавиатурой
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	reply.ReplyMarkup = h.keyboardPresenter.MainMenu()

	// Отправляем
	_, err = bot.Send(reply)
	if err != nil {
		logger.Error(ctx, "❌ Failed to send welcome message", zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ Welcome message sent successfully")
	return nil
}

// HandleMainMenuCallback обрабатывает возврат в главное меню из callback
func (h *StartHandler) HandleMainMenuCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	logger.Info(ctx, "🏠 Handling main menu callback",
		zap.Int64("user_id", query.From.ID),
		zap.String("username", query.From.UserName))

	// Получаем имя пользователя
	userName := query.From.FirstName
	if userName == "" {
		userName = query.From.UserName
	}

	// Подготавливаем данные для шаблона
	data := message.WelcomeData{
		UserName:         userName,
		PlayersCount:     14545, // TODO: получать из БД
		TeamsCount:       613,   // TODO: получать из БД
		TournamentsCount: 36,    // TODO: получать из БД
	}

	// Рендерим сообщение
	text, err := h.msgPresenter.RenderWelcome(data)
	if err != nil {
		logger.Error(ctx, "❌ Failed to render welcome message", zap.Error(err))
		text = "Добро пожаловать в Hockey Scout Bot!"
	}

	// Редактируем текущее сообщение
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.MainMenu()
	edit.ReplyMarkup = &markup

	_, err = bot.Send(edit)
	if err != nil {
		logger.Error(ctx, "❌ Error editing message to main menu", zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ Main menu displayed successfully")
	return nil
}
