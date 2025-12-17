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
	defer func() {
		_ = logger.Sync()
	}()

	// Загрузка конфигурации
	if err := config.Load(); err != nil {
		logger.Fatal(ctx, "Failed to load config", zap.Error(err))
	}
	cfg := config.AppConfig()

	// Создание DI контейнера для парсера
	factory := di.NewContainerFactory(cfg)
	container := factory.CreateParserContainer()

	// Запуск парсера
	logger.Info(ctx, "🏒 Starting parser...")
	orchestrator := container.Service().Orchestrator(ctx)
	if err := orchestrator.Run(ctx); err != nil {
		logger.Fatal(ctx, "Parser failed", zap.Error(err))
	}
	logger.Info(ctx, "✅ Parser completed successfully")
}
