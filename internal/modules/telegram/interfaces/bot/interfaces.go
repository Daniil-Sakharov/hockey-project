package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler интерфейс для обработчиков
type Handler interface {
	Handle(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error
}

// CommandHandler обработчик команд
type CommandHandler interface {
	HandleCommand(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error
}

// CallbackHandler обработчик callback query
type CallbackHandler interface {
	HandleCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error
	CanHandle(data string) bool
}

// MessageHandler обработчик сообщений
type MessageHandler interface {
	HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) error
}
