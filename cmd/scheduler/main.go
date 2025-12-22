package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fhspbStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/stats"
	juniorStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats"
	statsOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats/orchestrator"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories"
	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/scheduler/application"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/scheduler/infrastructure"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var runOnce = flag.Bool("run-once", false, "Run all jobs once and exit")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация логгера с OTEL
	logConfig := getLoggerConfig()
	if err := logger.Init(getEnv("LOG_LEVEL", "info"), true, logConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	config, err := modules.LoadSchedulerConfig("config/scheduler.yaml")
	if err != nil {
		logger.Fatal(ctx, "Failed to load scheduler config", zap.Error(err))
	}

	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	db, err := container.DB(ctx)
	if err != nil {
		logger.Fatal(ctx, "Failed to connect to database", zap.Error(err))
	}

	lockRepo := infrastructure.NewLockRepository(db.DB)

	// Инициализируем метрики
	metrics, err := infrastructure.NewSchedulerMetrics()
	if err != nil {
		logger.Warn(ctx, "Failed to initialize metrics: "+err.Error())
	}

	scheduler := application.NewSchedulerService(config, lockRepo, metrics)

	// Создаём репозитории
	tournamentRepo := repositories.NewTournamentPostgres(db)
	failedJobRepo := repositories.NewFailedJobRepository(db)

	// Регистрируем handlers
	registerHandlers(ctx, scheduler, container, config, tournamentRepo, failedJobRepo)

	// Обработка сигналов
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		logger.Info(ctx, "🛑 Received shutdown signal...")
		cancel()
	}()

	// Режим run-once или run_immediately
	if *runOnce || config.RunImmediately {
		logger.Info(ctx, "🚀 Running all jobs once...")
		runAllJobsOnce(ctx, scheduler, config)
		logger.Info(ctx, "✅ All jobs completed")
		return
	}

	// Bootstrap mode
	if config.BootstrapMode {
		logger.Info(ctx, "🚀 Running in BOOTSTRAP mode")
		if err := runBootstrap(ctx, container, tournamentRepo); err != nil {
			logger.Fatal(ctx, "Bootstrap failed", zap.Error(err))
		}
		logger.Info(ctx, "✅ Bootstrap completed")
		return
	}

	// Обычный режим - запускаем scheduler
	logger.Info(ctx, "🚀 Starting scheduler...")
	if err := scheduler.Start(ctx); err != nil {
		logger.Fatal(ctx, "Failed to start scheduler", zap.Error(err))
	}

	<-ctx.Done()

	logger.Info(ctx, "🛑 Stopping scheduler...")
	if err := scheduler.Stop(context.Background()); err != nil {
		logger.Error(ctx, "Failed to stop scheduler: "+err.Error())
	}

	logger.Info(ctx, "✅ Scheduler stopped")
}

func registerHandlers(
	ctx context.Context,
	scheduler *application.SchedulerService,
	container *di.Container,
	config *modules.SchedulerConfig,
	tournamentRepo *repositories.TournamentPostgres,
	failedJobRepo *repositories.FailedJobRepository,
) {
	metrics := scheduler.GetMetrics()

	// Junior Stats handler
	scheduler.RegisterHandler("junior_stats", func() error {
		return runJuniorStats(ctx, container, config, tournamentRepo, metrics)
	})

	// FHSPB Stats handler
	scheduler.RegisterHandler("fhspb_stats", func() error {
		return runFHSPBStats(ctx, container, config, tournamentRepo, metrics)
	})

	// Retry worker
	retryWorker := application.NewRetryWorker(failedJobRepo)
	scheduler.RegisterHandler("retry_worker", retryWorker.Run)

	// TODO: junior_parser, fhspb_parser handlers
	scheduler.RegisterHandler("junior_parser", func() error {
		logger.Info(ctx, "junior_parser: not implemented yet")
		return nil
	})
	scheduler.RegisterHandler("fhspb_parser", func() error {
		logger.Info(ctx, "fhspb_parser: not implemented yet")
		return nil
	})

	logger.Info(ctx, "📋 Handlers registered")
}

func runJuniorStats(
	ctx context.Context,
	container *di.Container,
	config *modules.SchedulerConfig,
	tournamentRepo *repositories.TournamentPostgres,
	metrics *infrastructure.SchedulerMetrics,
) error {
	logger.Info(ctx, "📊 Starting Junior Stats...")

	// Получаем турниры по приоритету
	tournaments, err := tournamentRepo.GetAllForParsing(ctx, true, "junior")
	if err != nil {
		return err
	}

	// Применяем лимит
	jobCfg, _ := config.GetJob("junior_stats")
	if jobCfg.MaxTournaments > 0 && len(tournaments) > jobCfg.MaxTournaments {
		tournaments = tournaments[:jobCfg.MaxTournaments]
	}

	if len(tournaments) == 0 {
		logger.Info(ctx, "✅ No tournaments need parsing")
		return nil
	}

	// Создаём orchestrator
	playerStatsRepo, err := container.ParsingPlayerStatisticsRepository(ctx)
	if err != nil {
		return err
	}
	parsingTournamentRepo, err := container.ParsingTournamentRepository(ctx)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	statsParser := stats.NewParser(httpClient)
	zapLogger := zap.NewNop()
	statsService := juniorStats.NewStatsParserService(statsParser, playerStatsRepo, zapLogger)

	stdLogger := log.New(os.Stdout, "[JUNIOR_STATS] ", log.LstdFlags|log.Lmsgprefix)
	orch := statsOrchestrator.NewStatsOrchestratorService(
		statsService,
		playerStatsRepo,
		parsingTournamentRepo,
		stdLogger,
	)

	// Парсим
	if err := orch.RunForTournaments(ctx, tournaments); err != nil {
		return err
	}

	// Записываем метрики
	if metrics != nil {
		metrics.RecordTournamentsParsed(ctx, "junior", int64(len(tournaments)))
	}

	// Обновляем last_stats_parsed_at
	for _, t := range tournaments {
		if err := tournamentRepo.UpdateLastStatsParsed(ctx, t.ID); err != nil {
			logger.Error(ctx, "Failed to update last_stats_parsed_at: "+err.Error())
		}
	}

	logger.Info(ctx, "✅ Junior Stats completed")
	return nil
}

func runFHSPBStats(
	ctx context.Context,
	container *di.Container,
	config *modules.SchedulerConfig,
	tournamentRepo *repositories.TournamentPostgres,
	metrics *infrastructure.SchedulerMetrics,
) error {
	logger.Info(ctx, "📊 Starting FHSPB Stats...")

	// Получаем турниры по приоритету
	tournaments, err := tournamentRepo.GetAllForParsing(ctx, true, "fhspb")
	if err != nil {
		return err
	}

	// Применяем лимит
	jobCfg, _ := config.GetJob("fhspb_stats")
	if jobCfg.MaxTournaments > 0 && len(tournaments) > jobCfg.MaxTournaments {
		tournaments = tournaments[:jobCfg.MaxTournaments]
	}

	if len(tournaments) == 0 {
		logger.Info(ctx, "✅ No FHSPB tournaments need parsing")
		return nil
	}

	db, err := container.DB(ctx)
	if err != nil {
		return err
	}

	// Конвертируем в FHSPB Tournament тип
	fhspbTournaments := make([]*fhspbrepo.Tournament, len(tournaments))
	for i, t := range tournaments {
		externalID := ""
		if t.ExternalID != nil {
			externalID = *t.ExternalID
		}
		fhspbTournaments[i] = &fhspbrepo.Tournament{
			ID:         t.ID,
			ExternalID: externalID,
			Name:       t.Name,
		}
	}

	// Создаём orchestrator
	fhspbClient := fhspb.NewClient()

	deps := fhspbStats.Dependencies{
		Client:               fhspbClient,
		TournamentRepo:       fhspbrepo.NewTournamentRepository(db),
		TeamRepo:             fhspbrepo.NewTeamRepository(db),
		PlayerRepo:           fhspbrepo.NewPlayerRepository(db),
		PlayerStatisticsRepo: fhspbrepo.NewPlayerStatisticsRepository(db),
		GoalieStatisticsRepo: fhspbrepo.NewGoalieStatisticsRepository(db),
	}

	orch := fhspbStats.NewOrchestrator(deps)

	// Парсим
	if err := orch.RunForTournaments(ctx, fhspbTournaments); err != nil {
		return err
	}

	// Записываем метрики
	if metrics != nil {
		metrics.RecordTournamentsParsed(ctx, "fhspb", int64(len(tournaments)))
	}

	// Обновляем last_stats_parsed_at
	for _, t := range tournaments {
		if err := tournamentRepo.UpdateLastStatsParsed(ctx, t.ID); err != nil {
			logger.Error(ctx, "Failed to update last_stats_parsed_at: "+err.Error())
		}
	}

	logger.Info(ctx, "✅ FHSPB Stats completed")
	return nil
}

func runAllJobsOnce(ctx context.Context, scheduler *application.SchedulerService, config *modules.SchedulerConfig) {
	handlers := scheduler.GetHandlers()

	for name := range config.EnabledJobs() {
		logger.Info(ctx, "▶️ Running job: "+name)

		handler, ok := handlers[name]
		if !ok {
			logger.Warn(ctx, "No handler for job: "+name)
			continue
		}

		if err := handler(); err != nil {
			logger.Error(ctx, "Job failed: "+name+": "+err.Error())
		} else {
			logger.Info(ctx, "✅ Job completed: "+name)
		}
	}
}

func runBootstrap(ctx context.Context, container *di.Container, tournamentRepo *repositories.TournamentPostgres) error {
	logger.Info(ctx, "📦 Bootstrap: checking if data exists...")

	hasData, err := tournamentRepo.HasData(ctx)
	if err != nil {
		return err
	}

	if hasData {
		logger.Info(ctx, "✅ Data already exists, skipping bootstrap")
		return nil
	}

	logger.Info(ctx, "📦 Bootstrap: Step 1/4 - Junior Parser (TODO)")
	logger.Info(ctx, "📦 Bootstrap: Step 2/4 - FHSPB Parser (TODO)")
	logger.Info(ctx, "📦 Bootstrap: Step 3/4 - Junior Stats (TODO)")
	logger.Info(ctx, "📦 Bootstrap: Step 4/4 - FHSPB Stats (TODO)")

	return nil
}

// getLoggerConfig возвращает конфигурацию логгера для OTEL
func getLoggerConfig() *logger.LoggerConfig {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil // OTEL отключен
	}

	return &logger.LoggerConfig{
		ServiceName:  getEnv("OTEL_SERVICE_NAME", "hockey-scheduler"),
		Environment:  getEnv("APP_ENV", "development"),
		OTLPEndpoint: endpoint,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
