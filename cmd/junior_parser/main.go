package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/parser"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type juniorConfigAdapter struct {
	cfg modules.JuniorConfig
}

func (a *juniorConfigAdapter) BaseURL() string    { return a.cfg.BaseURL }
func (a *juniorConfigAdapter) DomainWorkers() int { return a.cfg.DomainWorkers }
func (a *juniorConfigAdapter) MinBirthYear() int  { return a.cfg.MinBirthYear }

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := logger.Init("info", false, nil); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to load parsing config", zap.Error(err))
	}

	// Используем новые Parsing репозитории
	playerRepo, err := container.ParsingPlayerRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get player repository", zap.Error(err))
	}
	teamRepo, err := container.ParsingTeamRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get team repository", zap.Error(err))
	}
	tournamentRepo, err := container.ParsingTournamentRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get tournament repository", zap.Error(err))
	}
	playerTeamRepo, err := container.ParsingPlayerTeamRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get player_team repository", zap.Error(err))
	}

	juniorClient := junior.NewClient()
	juniorService := parser.NewJuniorService(juniorClient)
	configAdapter := &juniorConfigAdapter{cfg: parsingConfig.Junior}

	orch := parser.NewOrchestratorService(
		juniorService,
		playerRepo,
		teamRepo,
		tournamentRepo,
		playerTeamRepo,
		configAdapter,
	)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		logger.Info(ctx, "🛑 Received shutdown signal")
		cancel()
	}()

	logger.Info(ctx, "🏒 Starting Junior parser (v3 - fully modular)...",
		zap.Int("domain_workers", parsingConfig.Junior.DomainWorkers),
		zap.Int("min_birth_year", parsingConfig.Junior.MinBirthYear))

	if err := orch.Run(ctx); err != nil {
		logger.Fatal(ctx, "Parser failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ Junior parser completed successfully")
}
