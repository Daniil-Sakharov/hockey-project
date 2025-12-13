package initializer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/filter"
	profileHandler "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/profile"
	reportHandler "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/callback/search"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/handler/command"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/presenter/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router"
	routerMessage "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/router/message"
	"github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/template"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	playerRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/player"
	playerStatisticsRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/player_statistics"
	playerTeamRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/player_team"
	teamRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/team"
	tournamentRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	profileService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/profile"
	reportService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
	juniorService "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/orchestrator"
	statsService "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/stats_orchestrator"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/closer"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// diContainer DI контейнер для всех зависимостей
type diContainer struct {
	config *config.Config

	// Infrastructure
	db *sqlx.DB

	// Repositories
	playerRepository           player.Repository
	teamRepository             team.Repository
	tournamentRepository       tournament.Repository
	playerTeamRepository       player_team.Repository
	playerStatisticsRepository player_statistics.Repository

	// Clients
	juniorClient *junior.Client
	statsParser  *stats.Parser

	// Services
	juniorParserService      parser.JuniorParserService
	statsParserService       parser.StatsParserService
	orchestratorService      parser.OrchestratorService
	statsOrchestratorService parser.StatsOrchestratorService

	// Bot Services
	stateManager        bot.StateManager
	searchPlayerService bot.SearchPlayerService
	profileService      bot.ProfileService
	reportService       *reportService.Service

	// Telegram Bot dependencies
	templateEngine    *template.Engine
	msgPresenter      *message.MessagePresenter
	keyboardPresenter *keyboard.KeyboardPresenter
	startHandler      *command.StartHandler
	filterHandler     *filter.FilterHandler
	searchHandler     *search.Handler
	profileHandler    *profileHandler.Handler
	reportHandler     *reportHandler.Handler
	fioInputHandler   *routerMessage.FioInputHandler
	telegramRouter    *router.Router
	telegramBot       *telegram.Bot
}

// NewDiContainer создает новый DI контейнер
func NewDiContainer(cfg *config.Config) *diContainer {
	return &diContainer{
		config: cfg,
	}
}

// PostgresDB возвращает подключение к PostgreSQL (lazy initialization)
func (d *diContainer) PostgresDB(ctx context.Context) *sqlx.DB {
	if d.db == nil {
		db, err := sqlx.Connect("pgx", d.config.Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to PostgreSQL: %s", err.Error()))
		}

		// Настройки connection pool
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		db.SetConnMaxIdleTime(1 * time.Minute)

		if err := db.Ping(); err != nil {
			panic(fmt.Sprintf("failed to ping PostgreSQL: %s", err.Error()))
		}

		closer.AddNamed("PostgreSQL", func(ctx context.Context) error {
			return db.Close()
		})

		logger.Info(ctx, "✅ PostgreSQL connected")

		d.db = db
	}
	return d.db
}

// PlayerRepository возвращает репозиторий игроков
func (d *diContainer) PlayerRepository(ctx context.Context) player.Repository {
	if d.playerRepository == nil {
		db := d.PostgresDB(ctx)
		d.playerRepository = playerRepo.NewRepository(db)
	}
	return d.playerRepository
}

// TeamRepository возвращает репозиторий команд
func (d *diContainer) TeamRepository(ctx context.Context) team.Repository {
	if d.teamRepository == nil {
		db := d.PostgresDB(ctx)
		d.teamRepository = teamRepo.NewRepository(db)
	}
	return d.teamRepository
}

// TournamentRepository возвращает репозиторий турниров
func (d *diContainer) TournamentRepository(ctx context.Context) tournament.Repository {
	if d.tournamentRepository == nil {
		db := d.PostgresDB(ctx)
		d.tournamentRepository = tournamentRepo.NewRepository(db)
	}
	return d.tournamentRepository
}

// PlayerTeamRepository возвращает репозиторий связей player_teams
func (d *diContainer) PlayerTeamRepository(ctx context.Context) player_team.Repository {
	if d.playerTeamRepository == nil {
		db := d.PostgresDB(ctx)
		d.playerTeamRepository = playerTeamRepo.NewRepository(db)
	}
	return d.playerTeamRepository
}

// PlayerStatisticsRepository возвращает репозиторий статистики игроков
func (d *diContainer) PlayerStatisticsRepository(ctx context.Context) player_statistics.Repository {
	if d.playerStatisticsRepository == nil {
		db := d.PostgresDB(ctx)
		d.playerStatisticsRepository = playerStatisticsRepo.NewRepository(db)
	}
	return d.playerStatisticsRepository
}

// JuniorClient возвращает клиент для junior.fhr.ru
func (d *diContainer) JuniorClient() *junior.Client {
	if d.juniorClient == nil {
		d.juniorClient = junior.NewClient()
	}
	return d.juniorClient
}

// StatsParser возвращает парсер статистики
func (d *diContainer) StatsParser() *stats.Parser {
	if d.statsParser == nil {
		httpClient := &http.Client{
			Timeout: 60 * time.Second, // Увеличили timeout для медленных доменов
		}
		d.statsParser = stats.NewParser(httpClient)
	}
	return d.statsParser
}

// JuniorParserService возвращает сервис парсинга junior.fhr.ru
func (d *diContainer) JuniorParserService(ctx context.Context) parser.JuniorParserService {
	if d.juniorParserService == nil {
		client := d.JuniorClient()
		d.juniorParserService = juniorService.NewJuniorService(client)
	}
	return d.juniorParserService
}

// OrchestratorService возвращает orchestrator для парсинга
func (d *diContainer) OrchestratorService(ctx context.Context) parser.OrchestratorService {
	if d.orchestratorService == nil {
		juniorSvc := d.JuniorParserService(ctx)
		playerRepo := d.PlayerRepository(ctx)
		teamRepo := d.TeamRepository(ctx)
		tournamentRepo := d.TournamentRepository(ctx)
		playerTeamRepo := d.PlayerTeamRepository(ctx)

		d.orchestratorService = orchestrator.NewOrchestratorService(
			juniorSvc,
			playerRepo,
			teamRepo,
			tournamentRepo,
			playerTeamRepo,
		)
	}
	return d.orchestratorService
}

// StatsParserService возвращает сервис парсинга статистики
func (d *diContainer) StatsParserService(ctx context.Context) parser.StatsParserService {
	if d.statsParserService == nil {
		parser := d.StatsParser()
		repo := d.PlayerStatisticsRepository(ctx)

		// Создаем zap logger для stats parser
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapLogger, err := zapConfig.Build()
		if err != nil {
			log.Fatalf("Failed to create zap logger: %v", err)
		}

		// Добавляем поля контекста
		zapLogger = zapLogger.With(zap.String("service", "stats_parser"))

		d.statsParserService = statsService.NewStatsParserService(parser, repo, zapLogger)
	}
	return d.statsParserService
}

// StatsOrchestratorService возвращает orchestrator для парсинга статистики
func (d *diContainer) StatsOrchestratorService(ctx context.Context) parser.StatsOrchestratorService {
	if d.statsOrchestratorService == nil {
		statsService := d.StatsParserService(ctx)
		statsRepo := d.PlayerStatisticsRepository(ctx)
		tournamentRepo := d.TournamentRepository(ctx)
		logger := log.New(os.Stdout, "[STATS] ", log.LstdFlags|log.Lmsgprefix)

		d.statsOrchestratorService = stats_orchestrator.NewStatsOrchestratorService(
			statsService,
			statsRepo,
			tournamentRepo,
			logger,
		)
	}
	return d.statsOrchestratorService
}

// TemplateEngine возвращает template engine для рендеринга сообщений
func (d *diContainer) TemplateEngine() *template.Engine {
	if d.templateEngine == nil {
		engine, err := template.NewEngine()
		if err != nil {
			panic(fmt.Sprintf("failed to create template engine: %s", err.Error()))
		}
		d.templateEngine = engine
	}
	return d.templateEngine
}

// MessagePresenter возвращает presenter для форматирования сообщений
func (d *diContainer) MessagePresenter() *message.MessagePresenter {
	if d.msgPresenter == nil {
		engine := d.TemplateEngine()
		d.msgPresenter = message.NewMessagePresenter(engine)
	}
	return d.msgPresenter
}

// KeyboardPresenter возвращает presenter для создания клавиатур
func (d *diContainer) KeyboardPresenter() *keyboard.KeyboardPresenter {
	if d.keyboardPresenter == nil {
		d.keyboardPresenter = keyboard.NewKeyboardPresenter()
	}
	return d.keyboardPresenter
}

// StateManager возвращает менеджер состояния пользователей
func (d *diContainer) StateManager() bot.StateManager {
	if d.stateManager == nil {
		d.stateManager = bot.NewStateManager()
	}
	return d.stateManager
}

// SearchPlayerService возвращает сервис поиска игроков
func (d *diContainer) SearchPlayerService(ctx context.Context) bot.SearchPlayerService {
	if d.searchPlayerService == nil {
		playerRepo := d.PlayerRepository(ctx)
		d.searchPlayerService = bot.NewSearchPlayerService(playerRepo)
	}
	return d.searchPlayerService
}

// ProfileService возвращает сервис профилей игроков
func (d *diContainer) ProfileService(ctx context.Context) bot.ProfileService {
	if d.profileService == nil {
		playerRepo := d.PlayerRepository(ctx)
		statsRepo := d.PlayerStatisticsRepository(ctx)
		playerTeamRepo := d.PlayerTeamRepository(ctx)
		teamRepo := d.TeamRepository(ctx)
		d.profileService = profileService.NewService(playerRepo, statsRepo, playerTeamRepo, teamRepo)
	}
	return d.profileService
}

// ReportService возвращает сервис генерации отчетов
func (d *diContainer) ReportService(ctx context.Context) *reportService.Service {
	if d.reportService == nil {
		playerRepo := d.PlayerRepository(ctx)
		statsRepo := d.PlayerStatisticsRepository(ctx)
		playerTeamRepo := d.PlayerTeamRepository(ctx)
		teamRepo := d.TeamRepository(ctx)
		tournamentRepo := d.TournamentRepository(ctx)
		d.reportService = reportService.NewService(playerRepo, statsRepo, playerTeamRepo, teamRepo, tournamentRepo)
	}
	return d.reportService
}

// StartHandler возвращает обработчик команды /start
func (d *diContainer) StartHandler(ctx context.Context) *command.StartHandler {
	if d.startHandler == nil {
		msg := d.MessagePresenter()
		kbd := d.KeyboardPresenter()
		d.startHandler = command.NewStartHandler(msg, kbd)
	}
	return d.startHandler
}

// FilterHandler возвращает обработчик фильтров
func (d *diContainer) FilterHandler(ctx context.Context) *filter.FilterHandler {
	if d.filterHandler == nil {
		msg := d.MessagePresenter()
		kbd := d.KeyboardPresenter()
		state := d.StateManager()
		d.filterHandler = filter.NewFilterHandler(msg, kbd, state)
	}
	return d.filterHandler
}

// SearchHandler возвращает обработчик поиска
func (d *diContainer) SearchHandler(ctx context.Context) *search.Handler {
	if d.searchHandler == nil {
		msg := d.MessagePresenter()
		kbd := d.KeyboardPresenter()
		state := d.StateManager()
		searchSvc := d.SearchPlayerService(ctx)
		d.searchHandler = search.NewHandler(msg, kbd, state, searchSvc)
	}
	return d.searchHandler
}

// ProfileHandler возвращает обработчик профилей
func (d *diContainer) ProfileHandler(ctx context.Context) *profileHandler.Handler {
	if d.profileHandler == nil {
		kbd := d.KeyboardPresenter()
		templateEngine := d.TemplateEngine()
		profilePresenter := presenter.NewProfilePresenter(templateEngine)
		profileSvc := d.ProfileService(ctx)

		// Создаем zap logger для profile handler
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapLogger, err := zapConfig.Build()
		if err != nil {
			log.Fatalf("Failed to create zap logger for profile handler: %v", err)
		}
		zapLogger = zapLogger.With(zap.String("handler", "profile"))

		d.profileHandler = profileHandler.NewHandler(kbd, profilePresenter, profileSvc, zapLogger)
	}
	return d.profileHandler
}

// ReportHandler возвращает обработчик отчетов
func (d *diContainer) ReportHandler(ctx context.Context) *reportHandler.Handler {
	if d.reportHandler == nil {
		reportSvc := d.ReportService(ctx)

		// Создаем zap logger для report handler
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapLogger, err := zapConfig.Build()
		if err != nil {
			log.Fatalf("Failed to create zap logger for report handler: %v", err)
		}
		zapLogger = zapLogger.With(zap.String("handler", "report"))

		d.reportHandler = reportHandler.NewHandler(reportSvc, zapLogger)
	}
	return d.reportHandler
}

// FioInputHandler возвращает обработчик ввода ФИО
func (d *diContainer) FioInputHandler(ctx context.Context) *routerMessage.FioInputHandler {
	if d.fioInputHandler == nil {
		state := d.StateManager()
		filterHandler := d.FilterHandler(ctx)
		d.fioInputHandler = routerMessage.NewFioInputHandler(state, filterHandler)
	}
	return d.fioInputHandler
}

// TelegramRouter возвращает router для маршрутизации команд и callback
func (d *diContainer) TelegramRouter(ctx context.Context) *router.Router {
	if d.telegramRouter == nil {
		start := d.StartHandler(ctx)
		filterHandler := d.FilterHandler(ctx)
		searchHandler := d.SearchHandler(ctx)
		profileHandler := d.ProfileHandler(ctx)
		reportHandler := d.ReportHandler(ctx)
		fioInputHandler := d.FioInputHandler(ctx)
		stateManager := d.StateManager()
		d.telegramRouter = router.NewRouter(start, filterHandler, searchHandler, profileHandler, reportHandler, fioInputHandler, stateManager)
	}
	return d.telegramRouter
}

// TelegramBot возвращает экземпляр Telegram бота
func (d *diContainer) TelegramBot(ctx context.Context) *telegram.Bot {
	if d.telegramBot == nil {
		router := d.TelegramRouter(ctx)
		bot, err := telegram.NewBot(d.config.Telegram, router)
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s", err.Error()))
		}

		closer.AddNamed("TelegramBot", func(ctx context.Context) error {
			bot.Stop()
			return nil
		})

		logger.Info(ctx, "✅ Telegram Bot initialized")

		d.telegramBot = bot
	}
	return d.telegramBot
}
