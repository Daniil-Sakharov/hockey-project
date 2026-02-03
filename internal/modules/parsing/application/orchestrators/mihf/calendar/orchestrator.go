package calendar

import (
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	mihfClient "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/retry"
	"github.com/jmoiron/sqlx"
)

// CalendarConfig интерфейс конфигурации парсера календаря MIHF
type CalendarConfig interface {
	MinBirthYear() int
	MaxBirthYear() int
	RequestDelay() int
	GameWorkers() int
	ParseProtocol() bool
	SkipExisting() bool
	MaxSeasons() int
	TestSeason() string
	RetryEnabled() bool
	RetryMaxAttempts() int
	RetryDelay() time.Duration
}

// Source константа источника
const Source = "mihf.ru"

// Orchestrator оркестратор парсинга календаря MIHF
type Orchestrator struct {
	client *mihfClient.Client

	// Репозитории MIHF
	tournamentRepo *mihfrepo.TournamentRepository
	teamRepo       *mihfrepo.TeamRepository
	playerRepo     *mihfrepo.PlayerRepository
	playerTeamRepo *mihfrepo.PlayerTeamRepository

	// Общие репозитории
	matchRepo       repositories.MatchRepository
	matchEventRepo  repositories.MatchEventRepository
	matchLineupRepo repositories.MatchLineupRepository

	// Retry manager
	retryManager *retry.Manager

	// Конфигурация
	config CalendarConfig
}

// NewOrchestrator создает новый оркестратор
func NewOrchestrator(
	db *sqlx.DB,
	client *mihfClient.Client,
	tournamentRepo *mihfrepo.TournamentRepository,
	teamRepo *mihfrepo.TeamRepository,
	playerRepo *mihfrepo.PlayerRepository,
	playerTeamRepo *mihfrepo.PlayerTeamRepository,
	matchRepo repositories.MatchRepository,
	matchEventRepo repositories.MatchEventRepository,
	matchLineupRepo repositories.MatchLineupRepository,
	config CalendarConfig,
) *Orchestrator {
	retryManager := retry.NewManager(
		db,
		config.RetryMaxAttempts(),
		config.RetryDelay(),
	)

	return &Orchestrator{
		client:          client,
		tournamentRepo:  tournamentRepo,
		teamRepo:        teamRepo,
		playerRepo:      playerRepo,
		playerTeamRepo:  playerTeamRepo,
		matchRepo:       matchRepo,
		matchEventRepo:  matchEventRepo,
		matchLineupRepo: matchLineupRepo,
		retryManager:    retryManager,
		config:          config,
	}
}
