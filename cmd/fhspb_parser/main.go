package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/di"
	svc "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb/orchestrator"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Обработка сигналов завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info(ctx, "🛑 Received shutdown signal, stopping...")
		cancel()
	}()

	logger.Info(ctx, "🏒 Starting FHSPB parser...")

	// Настройка клиента
	client := fhspb.NewClient()
	fhspbCfg := cfg.FHSPB
	client.SetDelay(fhspbCfg.RequestDelay())

	// Создание зависимостей
	deps := svc.Dependencies{
		Client:         client,
		TournamentRepo: container.Repository().FHSPBTournament(ctx),
		TeamRepo:       container.Repository().FHSPBTeam(ctx),
		PlayerRepo:     container.Repository().FHSPBPlayer(ctx),
		PlayerTeamRepo: container.Repository().FHSPBPlayerTeam(ctx),
	}

	// Создание конфигурации
	parserCfg := svc.Config{
		MaxBirthYear:      fhspbCfg.MaxBirthYear(),
		TournamentWorkers: fhspbCfg.TournamentWorkers(),
		TeamWorkers:       fhspbCfg.TeamWorkers(),
		PlayerWorkers:     fhspbCfg.PlayerWorkers(),
	}

	// Запуск оркестратора
	orch := orchestrator.New(deps, parserCfg)
	if err := orch.Run(ctx); err != nil {
		logger.Fatal(ctx, "FHSPB parser failed", zap.Error(err))
	}

	logger.Info(ctx, "✅ FHSPB parser completed successfully")
}
