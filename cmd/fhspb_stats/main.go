package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	fhspbStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := logger.Init("info", false, nil); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	// Новый модульный DI контейнер
	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	// Обработка сигналов завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info(ctx, "🛑 Received shutdown signal, stopping...")
		cancel()
	}()

	logger.Info(ctx, "📊 Starting FHSPB Stats parser (v3 module)...")

	// Конфигурация парсинга
	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to load parsing config", zap.Error(err))
	}

	// Репозитории из нового DI
	tournamentRepo, err := container.FHSPBTournamentRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get tournament repository", zap.Error(err))
	}
	teamRepo, err := container.FHSPBTeamRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get team repository", zap.Error(err))
	}
	playerRepo, err := container.FHSPBPlayerRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get player repository", zap.Error(err))
	}
	playerStatsRepo, err := container.FHSPBPlayerStatisticsRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get player_statistics repository", zap.Error(err))
	}
	goalieStatsRepo, err := container.FHSPBGoalieStatisticsRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get goalie_statistics repository", zap.Error(err))
	}

	// Создание FHSPB клиента из НОВОГО модуля
	client := fhspb.NewClient()
	client.SetDelay(parsingConfig.FHSPB.RequestDelay)

	// Создание зависимостей
	deps := fhspbStats.Dependencies{
		Client:               client,
		TournamentRepo:       tournamentRepo,
		TeamRepo:             teamRepo,
		PlayerRepo:           playerRepo,
		PlayerStatisticsRepo: playerStatsRepo,
		GoalieStatisticsRepo: goalieStatsRepo,
	}

	// Запуск оркестратора из НОВОГО модуля
	orch := fhspbStats.NewOrchestrator(deps)
	if err := orch.Run(ctx); err != nil {
		logger.Fatal(ctx, "FHSPB Stats parser failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ FHSPB Stats parser completed successfully")
}
