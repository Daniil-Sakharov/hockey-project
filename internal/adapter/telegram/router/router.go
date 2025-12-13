package router

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/filter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/profile"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/search"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/command"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Router маршрутизирует входящие сообщения к соответствующим обработчикам
type Router struct {
	startHandler    *command.StartHandler
	filterHandler   *filter.FilterHandler
	searchHandler   *search.Handler
	profileHandler  *profile.Handler
	reportHandler   *report.Handler
	fioInputHandler *message.FioInputHandler
	stateManager    bot.StateManager
}

// NewRouter создает новый Router
func NewRouter(
	startHandler *command.StartHandler,
	filterHandler *filter.FilterHandler,
	searchHandler *search.Handler,
	profileHandler *profile.Handler,
	reportHandler *report.Handler,
	fioInputHandler *message.FioInputHandler,
	stateManager bot.StateManager,
) *Router {
	return &Router{
		startHandler:    startHandler,
		filterHandler:   filterHandler,
		searchHandler:   searchHandler,
		profileHandler:  profileHandler,
		reportHandler:   reportHandler,
		fioInputHandler: fioInputHandler,
		stateManager:    stateManager,
	}
}

// Route маршрутизирует обновление к нужному обработчику
func (r *Router) Route(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Обработка команд
	if update.Message != nil && update.Message.IsCommand() {
		r.handleCommand(ctx, bot, update.Message)
		return
	}

	// Обработка callback query (inline кнопки)
	if update.CallbackQuery != nil {
		r.handleCallback(ctx, bot, update.CallbackQuery)
		return
	}

	// Обработка обычных сообщений
	if update.Message != nil {
		r.handleMessage(ctx, bot, update.Message)
		return
	}
}

// FilterHandler возвращает filter handler (для callback роутеров)
func (r *Router) FilterHandler() *filter.FilterHandler {
	return r.filterHandler
}

// SearchHandler возвращает search handler (для callback роутеров)
func (r *Router) SearchHandler() *search.Handler {
	return r.searchHandler
}

// ProfileHandler возвращает profile handler (для callback роутеров)
func (r *Router) ProfileHandler() *profile.Handler {
	return r.profileHandler
}

// ReportHandler возвращает report handler (для callback роутеров)
func (r *Router) ReportHandler() *report.Handler {
	return r.reportHandler
}

// StartHandler возвращает start handler (для callback роутеров)
func (r *Router) StartHandler() *command.StartHandler {
	return r.startHandler
}
