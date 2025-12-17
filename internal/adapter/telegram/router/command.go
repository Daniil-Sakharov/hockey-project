package router

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handleCommand обрабатывает команды
func (r *Router) handleCommand(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		if err := r.startHandler.Handle(ctx, bot, message); err != nil {
			log.Printf("Error handling /start: %v", err)
		}
	default:
		// Неизвестная команда
		msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Используйте /start")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending unknown command message: %v", err)
		}
	}
}
