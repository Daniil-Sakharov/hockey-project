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

	// Загрузка конфигурации для парсера (без Telegram)
	if err := config.LoadForParser(); err != nil {
		logger.Fatal(ctx, "Failed to load config", zap.Error(err))
	}
	cfg := config.AppConfig()

	// Создание DI контейнера для парсера
	factory := di.NewContainerFactory(cfg)
	container := factory.CreateParserContainer()

	logger.Info(ctx, "📊 Starting Stats parser...")

	// Запуск Stats оркестратора
	orchestrator := container.Service().StatsOrchestrator(ctx)
	if err := orchestrator.Run(ctx); err != nil {
		logger.Fatal(ctx, "Stats parser failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ Stats parser completed successfully")
}
