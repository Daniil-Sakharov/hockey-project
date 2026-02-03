package fhmoscow

import (
	fhmoscowrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/retry"
)

// Orchestrator оркестратор парсинга fhmoscow.com
type Orchestrator struct {
	client               *fhmoscow.Client
	tournamentRepo       *fhmoscowrepo.TournamentRepository
	teamRepo             *fhmoscowrepo.TeamRepository
	playerRepo           *fhmoscowrepo.PlayerRepository
	playerTeamRepo       *fhmoscowrepo.PlayerTeamRepository
	playerStatisticsRepo *fhmoscowrepo.PlayerStatisticsRepository
	goalieStatisticsRepo *fhmoscowrepo.GoalieStatisticsRepository
	retryManager         *retry.Manager
	config               Config
}

// New создает новый оркестратор FHMoscow
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
