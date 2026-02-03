package mihf

import (
	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/retry"
)

// Orchestrator оркестратор парсинга stats.mihf.ru
type Orchestrator struct {
	client               *mihf.Client
	tournamentRepo       *mihfrepo.TournamentRepository
	teamRepo             *mihfrepo.TeamRepository
	playerRepo           *mihfrepo.PlayerRepository
	playerTeamRepo       *mihfrepo.PlayerTeamRepository
	playerStatisticsRepo *mihfrepo.PlayerStatisticsRepository
	goalieStatisticsRepo *mihfrepo.GoalieStatisticsRepository
	retryManager         *retry.Manager
	config               Config
}

// New создает новый оркестратор MIHF
func New(deps Dependencies, config Config) *Orchestrator {
	retryManager := retry.NewManager(
		deps.DB,
		config.RetryMaxAttempts(),
		config.RetryDelay(),
	)

	return &Orchestrator{
		client:               deps.Client,
		tournamentRepo:       deps.TournamentRepo,
		teamRepo:             deps.TeamRepo,
		playerRepo:           deps.PlayerRepo,
		playerTeamRepo:       deps.PlayerTeamRepo,
		playerStatisticsRepo: deps.PlayerStatisticsRepo,
		goalieStatisticsRepo: deps.GoalieStatisticsRepo,
		retryManager:         retryManager,
		config:               config,
	}
}
