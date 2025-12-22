package main

import (
	"context"
	"log"
	"os"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/migrator/pg"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	if err := logger.Init("info", false, nil); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	// Новый модульный DI контейнер
	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	logger.Info(ctx, "🔄 Connecting to database...")

	db, err := container.DB(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", zap.Error(err))
	}

	logger.Info(ctx, "✅ Connected to database")
	logger.Info(ctx, "🔄 Running migrations...")

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	migrator := pg.NewMigrator(db.DB, migrationsDir)
	if err := migrator.Up(ctx); err != nil {
		logger.Fatal(ctx, "Migration failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ Migrations applied successfully")
}
