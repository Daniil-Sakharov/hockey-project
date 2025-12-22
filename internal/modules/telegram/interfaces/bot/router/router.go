package router

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/router/callback"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// Router маршрутизирует входящие сообщения
type Router struct {
	filterHandler  FilterHandlerInterface
	searchHandler  SearchHandlerInterface
	profileHandler ProfileHandlerInterface
	reportHandler  ReportHandlerInterface
	startHandler   StartHandlerInterface
}

// NewRouter создает новый Router
func NewRouter(
	filterHandler FilterHandlerInterface,
	searchHandler SearchHandlerInterface,
	profileHandler ProfileHandlerInterface,
	reportHandler ReportHandlerInterface,
	startHandler StartHandlerInterface,
) *Router {
	return &Router{
		filterHandler:  filterHandler,
		searchHandler:  searchHandler,
		profileHandler: profileHandler,
		reportHandler:  reportHandler,
		startHandler:   startHandler,
	}
}

// Route обрабатывает входящее обновление
func (r *Router) Route(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	logger.Info(ctx, "📨 Router.Route called", zap.Int("update_id", update.UpdateID))

	switch {
	case update.Message != nil:
		logger.Debug(ctx, "📝 Routing to handleMessage")
		r.handleMessage(ctx, bot, update.Message)
	case update.CallbackQuery != nil:
		logger.Debug(ctx, "🔘 Routing to handleCallback", zap.String("data", update.CallbackQuery.Data))
		r.handleCallback(ctx, bot, update.CallbackQuery)
	default:
		logger.Debug(ctx, "❓ Unknown update type")
	}
}

// Реализация callback.Router интерфейса
func (r *Router) FilterHandler() callback.FilterHandler   { return r.filterHandler }
func (r *Router) SearchHandler() callback.SearchHandler   { return r.searchHandler }
func (r *Router) ProfileHandler() callback.ProfileHandler { return r.profileHandler }
func (r *Router) ReportHandler() callback.ReportHandler   { return r.reportHandler }
func (r *Router) StartHandler() callback.StartHandler     { return r.startHandler }
