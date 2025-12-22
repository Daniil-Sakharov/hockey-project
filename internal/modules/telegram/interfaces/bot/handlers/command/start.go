package command

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// StartHandler обрабатывает команду /start
type StartHandler struct {
	presenter *presenter.Presenter
	keyboard  *keyboard.KeyboardPresenter
}

// NewStartHandler создает новый StartHandler
func NewStartHandler(presenter *presenter.Presenter, keyboard *keyboard.KeyboardPresenter) *StartHandler {
	return &StartHandler{
		presenter: presenter,
		keyboard:  keyboard,
	}
}

// HandleStart обрабатывает команду /start
func (h *StartHandler) HandleStart(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error {
	logger.Info(ctx, "🚀 HandleStart called",
		zap.Int64("user_id", msg.From.ID),
		zap.String("username", msg.From.UserName))

	userName := ""
	if msg.From.FirstName != "" {
		userName = msg.From.FirstName
	}

	text, err := h.presenter.RenderWelcome(userName)
	if err != nil {
		logger.Error(ctx, "❌ Failed to render welcome", zap.Error(err))
		return err
	}
	logger.Debug(ctx, "✅ Welcome text rendered", zap.Int("text_len", len(text)))

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	reply.ReplyMarkup = h.keyboard.MainMenu()

	_, err = bot.Send(reply)
	if err != nil {
		logger.Error(ctx, "❌ Failed to send message", zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ Start message sent successfully", zap.Int64("user_id", msg.From.ID))
	return nil
}

// HandleMainMenuCallback обрабатывает возврат в главное меню
func (h *StartHandler) HandleMainMenuCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	userName := ""
	if query.From.FirstName != "" {
		userName = query.From.FirstName
	}

	text, err := h.presenter.RenderWelcome(userName)
	if err != nil {
		return err
	}

	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboard.MainMenu()
	edit.ReplyMarkup = &markup

	_, err = bot.Send(edit)
	return err
}
