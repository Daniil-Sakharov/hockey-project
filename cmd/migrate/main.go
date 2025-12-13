package main

import (
	"context"
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/migrator/pg"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx := context.Background()

	// Загружаем конфигурацию
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.AppConfig()

	log.Println("🔄 Connecting to database...")

	// Подключаемся к БД
	db, err := sqlx.Connect("pgx", cfg.Postgres.URI())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to database")
	log.Println("🔄 Running migrations...")

	// Создаем migrator
	migrator := pg.NewMigrator(db.DB, cfg.Postgres.MigrationsDir())

	// Запускаем миграции
	if err := migrator.Up(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Migrations applied successfully")
}
