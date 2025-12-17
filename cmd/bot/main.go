package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Инициализация логгера
	if err := logger.Init("info", false, nil); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	if err := config.Load(); err != nil {
		logger.Fatal(ctx, "Failed to load config", zap.Error(err))
	}
	cfg := config.AppConfig()

	// Создание DI контейнера для бота
	factory := di.NewContainerFactory(cfg)
	container := factory.CreateBotContainer()

	// Запуск бота
	logger.Info(ctx, "🤖 Starting Telegram bot...")
	if err := container.Telegram().Bot(ctx).Start(ctx); err != nil {
		logger.Fatal(ctx, "Bot failed", zap.Error(err))
	}
}
