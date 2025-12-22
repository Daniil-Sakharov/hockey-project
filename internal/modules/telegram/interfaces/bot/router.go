package bot

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Router маршрутизирует входящие сообщения к обработчикам
type Router struct {
	commandHandlers  map[string]CommandHandler
	callbackHandlers []CallbackHandler
	messageHandler   MessageHandler
}

// NewRouter создает новый Router
func NewRouter() *Router {
	return &Router{
		commandHandlers:  make(map[string]CommandHandler),
		callbackHandlers: make([]CallbackHandler, 0),
	}
}

// RegisterCommand регистрирует обработчик команды
func (r *Router) RegisterCommand(command string, handler CommandHandler) {
	r.commandHandlers[command] = handler
}

// RegisterCallback регистрирует обработчик callback
func (r *Router) RegisterCallback(handler CallbackHandler) {
	r.callbackHandlers = append(r.callbackHandlers, handler)
}

// RegisterMessage регистрирует обработчик сообщений
func (r *Router) RegisterMessage(handler MessageHandler) {
	r.messageHandler = handler
}

// Route маршрутизирует обновление к нужному обработчику
func (r *Router) Route(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil && update.Message.IsCommand() {
		r.handleCommand(ctx, bot, update.Message)
		return
	}

	if update.CallbackQuery != nil {
		r.handleCallback(ctx, bot, update.CallbackQuery)
		return
	}

	if update.Message != nil {
		r.handleMessage(ctx, bot, update.Message)
	}
}

func (r *Router) handleCommand(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	cmd := msg.Command()
	if handler, ok := r.commandHandlers[cmd]; ok {
		_ = handler.HandleCommand(ctx, bot, msg)
	}
}

func (r *Router) handleCallback(ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	for _, handler := range r.callbackHandlers {
		if handler.CanHandle(query.Data) {
			_ = handler.HandleCallback(ctx, bot, query)
			return
		}
	}
}

func (r *Router) handleMessage(ctx context.Context, bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	if r.messageHandler != nil {
		_ = r.messageHandler.HandleMessage(ctx, bot, msg)
	}
}

// CallbackPrefix возвращает префикс из callback data
func CallbackPrefix(data string) string {
	parts := strings.Split(data, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return data
}
