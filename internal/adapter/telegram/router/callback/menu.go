package callback

import (
	"context"
	"log"

	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/filter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/profile"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/search"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Router интерфейс для доступа к handlers
type Router interface {
	FilterHandler() *filter.FilterHandler
	SearchHandler() *search.Handler
	ProfileHandler() *profile.Handler
	StartHandler() *command.StartHandler
}

// HandleMenu обрабатывает callback главного меню
func HandleMenu(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	if len(parts) < 2 {
		return
	}

	command := parts[1]

	switch command {
	case cb.MenuSearch:
		// Открываем меню фильтров
		if err := r.FilterHandler().HandleFilterMenu(ctx, bot, query); err != nil {
			log.Printf("Error handling filter menu: %v", err)
		}
	case cb.MenuMain:
		// Возврат в главное меню
		if err := r.StartHandler().HandleMainMenuCallback(ctx, bot, query); err != nil {
			log.Printf("Error handling main menu: %v", err)
		}
	case cb.MenuStats:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "📊 Поиск по статистике - будет в следующей версии 🚧")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	case cb.MenuTeam:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "🏒 Поиск команды - будет в следующей версии 🚧")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	case cb.MenuHelp:
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "❓ Помощь - в разработке")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	default:
		log.Printf("Unknown menu command: %s", command)
	}
}
