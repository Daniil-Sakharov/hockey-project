package message

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleMessage обрабатывает обычные текстовые сообщения
func HandleMessage(ctx context.Context, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	// TODO: Обработка ввода ФИО и других текстовых данных
	// Пока просто предлагаем использовать /start
	msg := tgbotapi.NewMessage(message.Chat.ID, "Используйте /start для начала работы")
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
