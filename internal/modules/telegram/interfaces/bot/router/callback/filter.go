package callback

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/router/callback/filter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleFilter обрабатывает callback фильтров
func HandleFilter(r Router, ctx context.Context, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery, parts []string) {
	adapter := &filterRouterAdapter{router: r}
	filter.Handle(adapter, ctx, bot, query, parts)
}

type filterRouterAdapter struct {
	router Router
}

func (a *filterRouterAdapter) FilterHandler() filter.Handler {
	return a.router.FilterHandler()
}

func (a *filterRouterAdapter) SearchHandler() filter.SearchHandler {
	return a.router.SearchHandler()
}
