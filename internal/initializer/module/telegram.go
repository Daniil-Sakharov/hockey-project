package module

import (
	"context"
	"fmt"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/filter"
	profileHandler "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/profile"
	reportHandler "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/search"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/command"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router"
	routerMessage "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/template"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/closer"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Telegram содержит зависимости Telegram бота
type Telegram struct {
	config  *config.Config
	service *Service

	templateEngine    *template.Engine
	msgPresenter      *message.MessagePresenter
	keyboardPresenter *keyboard.KeyboardPresenter
	startHandler      *command.StartHandler
	filterHandler     *filter.FilterHandler
	searchHandler     *search.Handler
	profileHandler    *profileHandler.Handler
	reportHandler     *reportHandler.Handler
	fioInputHandler   *routerMessage.FioInputHandler
	telegramRouter    *router.Router
	telegramBot       *telegram.Bot
}

func NewTelegram(cfg *config.Config, service *Service) *Telegram {
	return &Telegram{config: cfg, service: service}
}

func (t *Telegram) TemplateEngine() *template.Engine {
	if t.templateEngine == nil {
		engine, err := template.NewEngine()
		if err != nil {
			panic(fmt.Sprintf("failed to create template engine: %s", err.Error()))
		}
		t.templateEngine = engine
	}
	return t.templateEngine
}

func (t *Telegram) MessagePresenter() *message.MessagePresenter {
	if t.msgPresenter == nil {
		t.msgPresenter = message.NewMessagePresenter(t.TemplateEngine())
	}
	return t.msgPresenter
}

func (t *Telegram) KeyboardPresenter() *keyboard.KeyboardPresenter {
	if t.keyboardPresenter == nil {
		t.keyboardPresenter = keyboard.NewKeyboardPresenter()
	}
	return t.keyboardPresenter
}

func (t *Telegram) StartHandler(ctx context.Context) *command.StartHandler {
	if t.startHandler == nil {
		t.startHandler = command.NewStartHandler(t.MessagePresenter(), t.KeyboardPresenter())
	}
	return t.startHandler
}

func (t *Telegram) FilterHandler(ctx context.Context) *filter.FilterHandler {
	if t.filterHandler == nil {
		t.filterHandler = filter.NewFilterHandler(t.MessagePresenter(), t.KeyboardPresenter(), t.service.StateManager())
	}
	return t.filterHandler
}

func (t *Telegram) SearchHandler(ctx context.Context) *search.Handler {
	if t.searchHandler == nil {
		t.searchHandler = search.NewHandler(t.MessagePresenter(), t.KeyboardPresenter(), t.service.StateManager(), t.service.SearchPlayer(ctx))
	}
	return t.searchHandler
}

func (t *Telegram) ProfileHandler(ctx context.Context) *profileHandler.Handler {
	if t.profileHandler == nil {
		profilePresenter := presenter.NewProfilePresenter(t.TemplateEngine())
		zapLogger := t.createLogger("profile")
		t.profileHandler = profileHandler.NewHandler(t.KeyboardPresenter(), profilePresenter, t.service.Profile(ctx), zapLogger)
	}
	return t.profileHandler
}

func (t *Telegram) ReportHandler(ctx context.Context) *reportHandler.Handler {
	if t.reportHandler == nil {
		zapLogger := t.createLogger("report")
		t.reportHandler = reportHandler.NewHandler(t.service.Report(ctx), zapLogger)
	}
	return t.reportHandler
}

func (t *Telegram) FioInputHandler(ctx context.Context) *routerMessage.FioInputHandler {
	if t.fioInputHandler == nil {
		t.fioInputHandler = routerMessage.NewFioInputHandler(t.service.StateManager(), t.FilterHandler(ctx))
	}
	return t.fioInputHandler
}

func (t *Telegram) Router(ctx context.Context) *router.Router {
	if t.telegramRouter == nil {
		t.telegramRouter = router.NewRouter(
			t.StartHandler(ctx),
			t.FilterHandler(ctx),
			t.SearchHandler(ctx),
			t.ProfileHandler(ctx),
			t.ReportHandler(ctx),
			t.FioInputHandler(ctx),
			t.service.StateManager(),
		)
	}
	return t.telegramRouter
}

func (t *Telegram) Bot(ctx context.Context) *telegram.Bot {
	if t.telegramBot == nil {
		bot, err := telegram.NewBot(t.config.Telegram, t.Router(ctx))
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s", err.Error()))
		}
		closer.AddNamed("TelegramBot", func(ctx context.Context) error {
			bot.Stop()
			return nil
		})
		logger.Info(ctx, "✅ Telegram Bot initialized")
		t.telegramBot = bot
	}
	return t.telegramBot
}

func (t *Telegram) createLogger(name string) *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Failed to create zap logger for %s: %v", name, err)
	}
	return zapLogger.With(zap.String("handler", name))
}
