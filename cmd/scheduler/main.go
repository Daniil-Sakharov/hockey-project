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

	// FHMoscow parser
	fhmoscowOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhmoscow"
	// Junior calendar
	juniorCalendar "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/calendar"
	jrCalendarParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/calendar"
	jrGameParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/game"
	jrStandingsParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/standings"
	// FHSPB parser
	fhspbParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/parser"
	fhspbParserOrch "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/parser/orchestrator"
	fhspbStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/stats"
	// FHSPB calendar
	fhspbCalendarOrch "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/calendar"
	fhspbCalendarParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/calendar"
	fhspbMatchParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/match"
	fhspbStandingsParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/standings"
	// Junior parser
	juniorParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/parser"
	juniorStats "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats"
	statsOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats/orchestrator"
	// MIHF parser
	mihfOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/mihf"
	// MIHF calendar
	mihfCalendarOrch "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/mihf/calendar"
	// Repositories
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories"
	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	// Sources
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf"
	// Scheduler
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

	// Junior parser handler
	scheduler.RegisterHandler("junior_parser", func() error {
		return runJuniorParser(ctx, container, config)
	})

	// FHSPB parser handler
	scheduler.RegisterHandler("fhspb_parser", func() error {
		return runFHSPBParser(ctx, container)
	})

	// MIHF parser handler
	scheduler.RegisterHandler("mihf_parser", func() error {
		return runMIHFParser(ctx, container, config)
	})

	// FHMoscow parser handler
	scheduler.RegisterHandler("fhmoscow_parser", func() error {
		return runFHMoscowParser(ctx, container)
	})

	// Junior calendar handler
	scheduler.RegisterHandler("junior_calendar", func() error {
		return runJuniorCalendar(ctx, container, config)
	})

	// FHSPB calendar handler
	scheduler.RegisterHandler("fhspb_calendar", func() error {
		return runFHSPBCalendar(ctx, container)
	})

	// MIHF calendar handler
	scheduler.RegisterHandler("mihf_calendar", func() error {
		return runMIHFCalendar(ctx, container)
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

func runMIHFParser(
	ctx context.Context,
	container *di.Container,
	config *modules.SchedulerConfig,
) error {
	logger.Info(ctx, "Starting MIHF Parser...")

	db, err := container.DB(ctx)
	if err != nil {
		return err
	}

	tournamentRepo, err := container.MIHFTournamentRepository(ctx)
	if err != nil {
		return err
	}
	teamRepo, err := container.MIHFTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerRepo, err := container.MIHFPlayerRepository(ctx)
	if err != nil {
		return err
	}
	playerTeamRepo, err := container.MIHFPlayerTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerStatisticsRepo, err := container.MIHFPlayerStatisticsRepository(ctx)
	if err != nil {
		return err
	}
	goalieStatisticsRepo, err := container.MIHFGoalieStatisticsRepository(ctx)
	if err != nil {
		return err
	}

	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		return err
	}

	client := mihf.NewClient()
	client.SetDelay(parsingConfig.MIHF.RequestDelay)

	deps := mihfOrchestrator.Dependencies{
		DB:                   db,
		Client:               client,
		TournamentRepo:       tournamentRepo,
		TeamRepo:             teamRepo,
		PlayerRepo:           playerRepo,
		PlayerTeamRepo:       playerTeamRepo,
		PlayerStatisticsRepo: playerStatisticsRepo,
		GoalieStatisticsRepo: goalieStatisticsRepo,
	}

	configAdapter := &mihfConfigAdapter{cfg: parsingConfig.MIHF}

	orch := mihfOrchestrator.New(deps, configAdapter)
	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "MIHF Parser completed")
	return nil
}

type mihfConfigAdapter struct {
	cfg modules.MIHFConfig
}

func (a *mihfConfigAdapter) MinBirthYear() int          { return a.cfg.MinBirthYear }
func (a *mihfConfigAdapter) MaxBirthYear() int          { return a.cfg.MaxBirthYear }
func (a *mihfConfigAdapter) SeasonWorkers() int         { return a.cfg.SeasonWorkers }
func (a *mihfConfigAdapter) TournamentWorkers() int     { return a.cfg.TournamentWorkers }
func (a *mihfConfigAdapter) TeamWorkers() int           { return a.cfg.TeamWorkers }
func (a *mihfConfigAdapter) PlayerWorkers() int         { return a.cfg.PlayerWorkers }
func (a *mihfConfigAdapter) RetryEnabled() bool         { return a.cfg.RetryEnabled }
func (a *mihfConfigAdapter) RetryMaxAttempts() int      { return a.cfg.RetryMaxAttempts }
func (a *mihfConfigAdapter) RetryDelay() time.Duration  { return a.cfg.RetryDelay }
func (a *mihfConfigAdapter) MaxSeasons() int            { return a.cfg.MaxSeasons }
func (a *mihfConfigAdapter) TestSeason() string         { return a.cfg.TestSeason }

// ============================================================================
// Junior Parser
// ============================================================================

func runJuniorParser(ctx context.Context, container *di.Container, schedulerConfig *modules.SchedulerConfig) error {
	logger.Info(ctx, "🏒 Starting Junior Parser...")

	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		return err
	}

	playerRepo, err := container.ParsingPlayerRepository(ctx)
	if err != nil {
		return err
	}
	teamRepo, err := container.ParsingTeamRepository(ctx)
	if err != nil {
		return err
	}
	tournamentRepo, err := container.ParsingTournamentRepository(ctx)
	if err != nil {
		return err
	}
	playerTeamRepo, err := container.ParsingPlayerTeamRepository(ctx)
	if err != nil {
		return err
	}

	// Получаем max_tournaments из scheduler config
	maxTournaments := 0
	if jobCfg, ok := schedulerConfig.GetJob("junior_parser"); ok {
		maxTournaments = jobCfg.MaxTournaments
	}

	juniorClient := junior.NewClient()
	juniorService := juniorParser.NewJuniorService(juniorClient)
	configAdapter := &juniorConfigAdapter{cfg: parsingConfig.Junior, maxTournaments: maxTournaments}

	orch := juniorParser.NewOrchestratorService(
		juniorService,
		playerRepo,
		teamRepo,
		tournamentRepo,
		playerTeamRepo,
		configAdapter,
	)

	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "✅ Junior Parser completed")
	return nil
}

type juniorConfigAdapter struct {
	cfg            modules.JuniorConfig
	maxTournaments int
}

func (a *juniorConfigAdapter) BaseURL() string       { return a.cfg.BaseURL }
func (a *juniorConfigAdapter) DomainWorkers() int    { return a.cfg.DomainWorkers }
func (a *juniorConfigAdapter) MinBirthYear() int     { return a.cfg.MinBirthYear }
func (a *juniorConfigAdapter) MaxTournaments() int   { return a.maxTournaments }

// ============================================================================
// FHSPB Parser
// ============================================================================

func runFHSPBParser(ctx context.Context, container *di.Container) error {
	logger.Info(ctx, "🏒 Starting FHSPB Parser...")

	db, err := container.DB(ctx)
	if err != nil {
		return err
	}

	tournamentRepo, err := container.FHSPBTournamentRepository(ctx)
	if err != nil {
		return err
	}
	teamRepo, err := container.FHSPBTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerRepo, err := container.FHSPBPlayerRepository(ctx)
	if err != nil {
		return err
	}
	playerTeamRepo, err := container.FHSPBPlayerTeamRepository(ctx)
	if err != nil {
		return err
	}

	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		return err
	}

	client := fhspb.NewClient()
	client.SetDelay(parsingConfig.FHSPB.RequestDelay)

	deps := fhspbParser.Dependencies{
		DB:             db,
		Client:         client,
		TournamentRepo: tournamentRepo,
		TeamRepo:       teamRepo,
		PlayerRepo:     playerRepo,
		PlayerTeamRepo: playerTeamRepo,
	}

	configAdapter := &fhspbConfigAdapter{cfg: parsingConfig.FHSPB}

	orch := fhspbParserOrch.New(deps, configAdapter)
	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "✅ FHSPB Parser completed")
	return nil
}

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

// ============================================================================
// FHMoscow Parser
// ============================================================================

func runFHMoscowParser(ctx context.Context, container *di.Container) error {
	logger.Info(ctx, "🏒 Starting FHMoscow Parser...")

	db, err := container.DB(ctx)
	if err != nil {
		return err
	}

	tournamentRepo, err := container.FHMoscowTournamentRepository(ctx)
	if err != nil {
		return err
	}
	teamRepo, err := container.FHMoscowTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerRepo, err := container.FHMoscowPlayerRepository(ctx)
	if err != nil {
		return err
	}
	playerTeamRepo, err := container.FHMoscowPlayerTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerStatisticsRepo, err := container.FHMoscowPlayerStatisticsRepository(ctx)
	if err != nil {
		return err
	}
	goalieStatisticsRepo, err := container.FHMoscowGoalieStatisticsRepository(ctx)
	if err != nil {
		return err
	}

	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		return err
	}

	client := fhmoscow.NewClient()
	client.SetDelay(parsingConfig.FHMoscow.RequestDelay)

	deps := fhmoscowOrchestrator.Dependencies{
		DB:                   db,
		Client:               client,
		TournamentRepo:       tournamentRepo,
		TeamRepo:             teamRepo,
		PlayerRepo:           playerRepo,
		PlayerTeamRepo:       playerTeamRepo,
		PlayerStatisticsRepo: playerStatisticsRepo,
		GoalieStatisticsRepo: goalieStatisticsRepo,
	}

	configAdapter := &fhmoscowConfigAdapter{cfg: parsingConfig.FHMoscow}

	orch := fhmoscowOrchestrator.New(deps, configAdapter)
	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "✅ FHMoscow Parser completed")
	return nil
}

type fhmoscowConfigAdapter struct {
	cfg modules.FHMoscowConfig
}

func (a *fhmoscowConfigAdapter) MinBirthYear() int          { return a.cfg.MinBirthYear }
func (a *fhmoscowConfigAdapter) SeasonWorkers() int         { return a.cfg.SeasonWorkers }
func (a *fhmoscowConfigAdapter) TournamentWorkers() int     { return a.cfg.TournamentWorkers }
func (a *fhmoscowConfigAdapter) TeamWorkers() int           { return a.cfg.TeamWorkers }
func (a *fhmoscowConfigAdapter) PlayerWorkers() int         { return a.cfg.PlayerWorkers }
func (a *fhmoscowConfigAdapter) RetryEnabled() bool         { return a.cfg.RetryEnabled }
func (a *fhmoscowConfigAdapter) RetryMaxAttempts() int      { return a.cfg.RetryMaxAttempts }
func (a *fhmoscowConfigAdapter) RetryDelay() time.Duration  { return a.cfg.RetryDelay }
func (a *fhmoscowConfigAdapter) MaxSeasons() int            { return a.cfg.MaxSeasons }
func (a *fhmoscowConfigAdapter) TestSeason() string         { return a.cfg.TestSeason }
func (a *fhmoscowConfigAdapter) ScanPlayers() bool          { return a.cfg.ScanPlayers }
func (a *fhmoscowConfigAdapter) MaxPlayerID() int           { return a.cfg.MaxPlayerID }

// ============================================================================
// Junior Calendar
// ============================================================================

func runJuniorCalendar(ctx context.Context, container *di.Container, schedulerConfig *modules.SchedulerConfig) error {
	logger.Info(ctx, "🗓️ Starting Junior Calendar Parser...")

	matchRepo, err := container.MatchRepository(ctx)
	if err != nil {
		return err
	}
	matchEventRepo, err := container.MatchEventRepository(ctx)
	if err != nil {
		return err
	}
	matchLineupRepo, err := container.MatchLineupRepository(ctx)
	if err != nil {
		return err
	}
	standingRepo, err := container.StandingRepository(ctx)
	if err != nil {
		return err
	}
	tournamentRepo, err := container.ParsingTournamentRepository(ctx)
	if err != nil {
		return err
	}
	teamRepo, err := container.ParsingTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerRepo, err := container.ParsingPlayerRepository(ctx)
	if err != nil {
		return err
	}

	// Получаем max_tournaments из scheduler config
	maxTournaments := 0
	if jobCfg, ok := schedulerConfig.GetJob("junior_calendar"); ok {
		maxTournaments = jobCfg.MaxTournaments
	}

	juniorClient := junior.NewClient()
	calendarParser := jrCalendarParser.NewParser(juniorClient)
	gameParser := jrGameParser.NewParser(juniorClient)
	standingsParser := jrStandingsParser.NewParser(juniorClient)

	configAdapter := &juniorCalendarConfigAdapter{maxTournaments: maxTournaments}

	orch := juniorCalendar.NewOrchestrator(
		juniorClient, // HTTP клиент для AJAX-запросов
		calendarParser,
		gameParser,
		standingsParser,
		nil, // profileParser - будет добавлен позже
		matchRepo,
		matchEventRepo,
		matchLineupRepo,
		standingRepo,
		tournamentRepo,
		teamRepo,
		playerRepo,
		configAdapter,
	)

	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "✅ Junior Calendar Parser completed")
	return nil
}

type juniorCalendarConfigAdapter struct {
	maxTournaments int
}

func (a *juniorCalendarConfigAdapter) RequestDelay() int      { return 150 }
func (a *juniorCalendarConfigAdapter) TournamentWorkers() int { return 3 }
func (a *juniorCalendarConfigAdapter) GameWorkers() int       { return 5 }
func (a *juniorCalendarConfigAdapter) ParseProtocol() bool    { return true }
func (a *juniorCalendarConfigAdapter) ParseLineups() bool     { return true }
func (a *juniorCalendarConfigAdapter) SkipExisting() bool     { return true }
func (a *juniorCalendarConfigAdapter) MaxTournaments() int    { return a.maxTournaments }

// ============================================================================
// FHSPB Calendar
// ============================================================================

func runFHSPBCalendar(ctx context.Context, container *di.Container) error {
	logger.Info(ctx, "Starting FHSPB Calendar Parser...")

	// Репозитории
	matchRepo, err := container.MatchRepository(ctx)
	if err != nil {
		return err
	}
	matchEventRepo, err := container.MatchEventRepository(ctx)
	if err != nil {
		return err
	}
	matchLineupRepo, err := container.MatchLineupRepository(ctx)
	if err != nil {
		return err
	}
	standingRepo, err := container.StandingRepository(ctx)
	if err != nil {
		return err
	}
	matchTeamStatsRepo, err := container.MatchTeamStatsRepository(ctx)
	if err != nil {
		return err
	}
	tournamentRepo, err := container.FHSPBTournamentRepository(ctx)
	if err != nil {
		return err
	}
	teamRepo, err := container.FHSPBTeamRepository(ctx)
	if err != nil {
		return err
	}
	playerRepo, err := container.FHSPBPlayerRepository(ctx)
	if err != nil {
		return err
	}

	// Client и парсеры
	client := fhspb.NewClient()
	calendarParser := fhspbCalendarParser.NewParser()
	matchParser := fhspbMatchParser.NewParser()
	standingsParser := fhspbStandingsParser.NewParser()

	// Config adapter
	configAdapter := &fhspbCalendarConfigAdapter{}

	// Orchestrator
	orch := fhspbCalendarOrch.NewOrchestrator(
		client,
		calendarParser,
		matchParser,
		standingsParser,
		tournamentRepo,
		teamRepo,
		playerRepo,
		matchRepo,
		matchEventRepo,
		matchLineupRepo,
		standingRepo,
		matchTeamStatsRepo,
		configAdapter,
	)

	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "FHSPB Calendar Parser completed")
	return nil
}

type fhspbCalendarConfigAdapter struct{}

func (a *fhspbCalendarConfigAdapter) RequestDelay() int   { return 150 }
func (a *fhspbCalendarConfigAdapter) GameWorkers() int    { return 5 }
func (a *fhspbCalendarConfigAdapter) ParseProtocol() bool { return true }
func (a *fhspbCalendarConfigAdapter) ParseLineups() bool  { return true }
func (a *fhspbCalendarConfigAdapter) SkipExisting() bool  { return true }

// ============================================================================
// MIHF Calendar
// ============================================================================

func runMIHFCalendar(ctx context.Context, container *di.Container) error {
	logger.Info(ctx, "Starting MIHF Calendar Parser...")

	db, err := container.DB(ctx)
	if err != nil {
		return err
	}

	parsingConfig, err := container.Config().Parsing(ctx)
	if err != nil {
		return err
	}

	// Репозитории
	matchRepo, err := container.MatchRepository(ctx)
	if err != nil {
		return err
	}
	matchEventRepo, err := container.MatchEventRepository(ctx)
	if err != nil {
		return err
	}
	matchLineupRepo, err := container.MatchLineupRepository(ctx)
	if err != nil {
		return err
	}

	// MIHF репозитории
	tournamentRepo := mihfrepo.NewTournamentRepository(db)
	teamRepo := mihfrepo.NewTeamRepository(db)
	playerRepo := mihfrepo.NewPlayerRepository(db)
	playerTeamRepo := mihfrepo.NewPlayerTeamRepository(db)

	// Client
	client := mihf.NewClient()
	client.SetDelay(parsingConfig.MIHF.RequestDelay)

	// Config adapter (использует parsingConfig.MIHF для общих параметров)
	configAdapter := &mihfCalendarConfigAdapter{cfg: parsingConfig.MIHF}

	// Orchestrator
	orch := mihfCalendarOrch.NewOrchestrator(
		db,
		client,
		tournamentRepo,
		teamRepo,
		playerRepo,
		playerTeamRepo,
		matchRepo,
		matchEventRepo,
		matchLineupRepo,
		configAdapter,
	)

	if err := orch.Run(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "MIHF Calendar Parser completed")
	return nil
}

type mihfCalendarConfigAdapter struct {
	cfg modules.MIHFConfig
}

func (a *mihfCalendarConfigAdapter) MinBirthYear() int         { return a.cfg.MinBirthYear }
func (a *mihfCalendarConfigAdapter) MaxBirthYear() int         { return a.cfg.MaxBirthYear }
func (a *mihfCalendarConfigAdapter) RequestDelay() int         { return 150 }
func (a *mihfCalendarConfigAdapter) GameWorkers() int          { return 5 }
func (a *mihfCalendarConfigAdapter) ParseProtocol() bool       { return true }
func (a *mihfCalendarConfigAdapter) SkipExisting() bool        { return true }
func (a *mihfCalendarConfigAdapter) MaxSeasons() int           { return a.cfg.MaxSeasons }
func (a *mihfCalendarConfigAdapter) TestSeason() string        { return a.cfg.TestSeason }
func (a *mihfCalendarConfigAdapter) RetryEnabled() bool        { return a.cfg.RetryEnabled }
func (a *mihfCalendarConfigAdapter) RetryMaxAttempts() int     { return a.cfg.RetryMaxAttempts }
func (a *mihfCalendarConfigAdapter) RetryDelay() time.Duration { return a.cfg.RetryDelay }

// ============================================================================
// Utilities
// ============================================================================

func runAllJobsOnce(ctx context.Context, scheduler *application.SchedulerService, config *modules.SchedulerConfig) {
	handlers := scheduler.GetHandlers()

	// Используем отсортированный список jobs (по полю order)
	for _, jobWithName := range config.EnabledJobsOrdered() {
		name := jobWithName.Name
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
