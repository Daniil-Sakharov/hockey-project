package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	fhspbParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/parser"
	fhspbOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/parser/orchestrator"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

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

	logger.Info(ctx, "🏒 Starting FHSPB parser (v3 module)...")

	// Получаем DB
	db, err := container.DB(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get database", zap.Error(err))
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
	playerTeamRepo, err := container.FHSPBPlayerTeamRepository(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to get player_team repository", zap.Error(err))
	}

	// Конфигурация парсинга
	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to load parsing config", zap.Error(err))
	}

	// Создание FHSPB клиента из НОВОГО модуля
	client := fhspb.NewClient()
	client.SetDelay(parsingConfig.FHSPB.RequestDelay)

	// Создание зависимостей
	deps := fhspbParser.Dependencies{
		DB:             db,
		Client:         client,
		TournamentRepo: tournamentRepo,
		TeamRepo:       teamRepo,
		PlayerRepo:     playerRepo,
		PlayerTeamRepo: playerTeamRepo,
	}

	// Адаптер для конфигурации
	configAdapter := &fhspbConfigAdapter{cfg: parsingConfig.FHSPB}

	// Запуск оркестратора из НОВОГО модуля
	orch := fhspbOrchestrator.New(deps, configAdapter)
	if err := orch.Run(ctx); err != nil {
		logger.Fatal(ctx, "FHSPB parser failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ FHSPB parser completed successfully")
}

// fhspbConfigAdapter адаптер для интерфейса application.Config
type fhspbConfigAdapter struct {
	cfg modules.FHSPBConfig
}

func (a *fhspbConfigAdapter) MaxBirthYear() int         { return a.cfg.MaxBirthYear }
func (a *fhspbConfigAdapter) TournamentWorkers() int    { return a.cfg.TournamentWorkers }
func (a *fhspbConfigAdapter) TeamWorkers() int          { return a.cfg.TeamWorkers }
func (a *fhspbConfigAdapter) PlayerWorkers() int        { return a.cfg.PlayerWorkers }
func (a *fhspbConfigAdapter) Mode() string              { return a.cfg.Mode }
func (a *fhspbConfigAdapter) RetryEnabled() bool        { return a.cfg.RetryEnabled }
func (a *fhspbConfigAdapter) RetryMaxAttempts() int     { return a.cfg.RetryMaxAttempts }
func (a *fhspbConfigAdapter) RetryDelay() time.Duration { return a.cfg.RetryDelay }
