package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	juniorStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats"
	statsOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats/orchestrator"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
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

	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info(ctx, "🛑 Received shutdown signal, stopping...")
		cancel()
	}()

	logger.Info(ctx, "📊 Starting Junior Stats parser (v3 module)...")

	// Используем новые Parsing репозитории
	playerStatsRepo, err := container.ParsingPlayerStatisticsRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get player_statistics repository", zap.Error(err))
	}
	tournamentRepo, err := container.ParsingTournamentRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get tournament repository", zap.Error(err))
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	statsParser := stats.NewParser(httpClient)
	zapLogger := zap.NewNop()
	statsService := juniorStats.NewStatsParserService(statsParser, playerStatsRepo, zapLogger)

	stdLogger := log.New(os.Stdout, "[STATS] ", log.LstdFlags|log.Lmsgprefix)
	orch := statsOrchestrator.NewStatsOrchestratorService(
		statsService,
		playerStatsRepo,
		tournamentRepo,
		stdLogger,
	)

	if err := orch.Run(ctx); err != nil {
		logger.Fatal(ctx, "Junior Stats parser failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ Junior Stats parser completed successfully")
}
