package command

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		log.Printf("Failed to render welcome message: %v", err)
		text = "Добро пожаловать в Hockey Scout Bot!"
	}

	// Создаем сообщение с inline клавиатурой
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	reply.ReplyMarkup = h.keyboardPresenter.MainMenu()

	// Отправляем
	_, err = bot.Send(reply)
	if err != nil {
		return err
	}

	return nil
}

// HandleMainMenuCallback обрабатывает возврат в главное меню из callback
func (h *StartHandler) HandleMainMenuCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
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
		log.Printf("Failed to render welcome message: %v", err)
		text = "Добро пожаловать в Hockey Scout Bot!"
	}

	// Редактируем текущее сообщение
	edit := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, text)
	edit.ParseMode = "Markdown"
	markup := h.keyboardPresenter.MainMenu()
	edit.ReplyMarkup = &markup

	_, err = bot.Send(edit)
	if err != nil {
		log.Printf("Error editing message to main menu: %v", err)
		return err
	}

	return nil
}
