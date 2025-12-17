package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/migrator/pg"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	// Создание DI контейнера для миграций
	factory := di.NewContainerFactory(cfg)
	container := factory.CreateMigrateContainer()

	logger.Info(ctx, "🔄 Connecting to database...")

	db := container.Infrastructure().PostgresDB(ctx)

	logger.Info(ctx, "✅ Connected to database")
	logger.Info(ctx, "🔄 Running migrations...")

	migrator := pg.NewMigrator(db.DB, cfg.Postgres.MigrationsDir())
	if err := migrator.Up(ctx); err != nil {
		logger.Fatal(ctx, "Migration failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ Migrations applied successfully")
}
