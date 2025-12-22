package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/handlers/callback/filter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/handlers/callback/profile"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/handlers/callback/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/handlers/callback/search"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/handlers/command"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/router"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера с OTEL
	logConfig := getLoggerConfig()
	if err := logger.Init(getEnv("LOG_LEVEL", "info"), true, logConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info(ctx, "🚀 Starting HockeyStats Bot v3...")

	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	telegramConfig, err := container.Config().Telegram(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to load telegram config", zap.Error(err))
	}
	logger.Info(ctx, "✅ Telegram config loaded")

	// Presenter и Keyboard
	presenter, err := container.TelegramPresenter()
	if err != nil {
		logger.Fatal(ctx, "Failed to get presenter", zap.Error(err))
	}
	keyboard := container.TelegramKeyboardPresenter()
	logger.Info(ctx, "✅ Presenter and Keyboard initialized")

	// Services
	stateService := container.TelegramUserStateService()

	searchService, err := container.TelegramPlayerSearchService(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get search service", zap.Error(err))
	}

	profileService, err := container.TelegramProfileService(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get profile service", zap.Error(err))
	}

	reportService, err := container.TelegramReportService(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get report service", zap.Error(err))
	}
	logger.Info(ctx, "✅ All services initialized")

	// Handlers
	startHandler := command.NewStartHandler(presenter, keyboard)
	filterHandler := filter.NewHandler(presenter, keyboard, stateService)
	searchHandler := search.NewHandler(presenter, keyboard, stateService, searchService)
	profileHandler := profile.NewHandler(presenter, keyboard, profileService)
	reportHandler := report.NewHandler(reportService)
	logger.Info(ctx, "✅ All handlers initialized")

	// Router
	botRouter := router.NewRouter(filterHandler, searchHandler, profileHandler, reportHandler, startHandler)
	logger.Info(ctx, "✅ Router initialized")

	telegramBot, err := bot.NewBot(telegramConfig, botRouter)
	if err != nil {
		logger.Fatal(ctx, "Failed to create bot", zap.Error(err))
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info(ctx, "🛑 Received shutdown signal")
		cancel()
	}()

	logger.Info(ctx, "🤖 Bot is ready! Waiting for updates...")
	if err := telegramBot.Start(ctx); err != nil && err != context.Canceled {
		logger.Fatal(ctx, "Bot failed", zap.Error(err))
	}
}

// getLoggerConfig возвращает конфигурацию логгера для OTEL
func getLoggerConfig() *logger.LoggerConfig {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil // OTEL отключен
	}

	return &logger.LoggerConfig{
		ServiceName:  getEnv("OTEL_SERVICE_NAME", "hockey-bot"),
		Environment:  getEnv("APP_ENV", "development"),
		OTLPEndpoint: endpoint,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
