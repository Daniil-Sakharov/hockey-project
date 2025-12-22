package parser

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/fhspb/parser"
	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/retry"
)

// Проверка реализации интерфейса
var _ parser.Service = (*Orchestrator)(nil)

// Orchestrator оркестратор парсинга fhspb.ru
type Orchestrator struct {
	client         *fhspb.Client
	tournamentRepo *fhspbrepo.TournamentRepository
	teamRepo       *fhspbrepo.TeamRepository
	playerRepo     *fhspbrepo.PlayerRepository
	playerTeamRepo *fhspbrepo.PlayerTeamRepository
	retryManager   *retry.Manager
	config         parser.Config
}

// New создает новый оркестратор
func New(deps parser.Dependencies, config parser.Config) *Orchestrator {
	retryManager := retry.NewManager(
		deps.DB,
		config.RetryMaxAttempts(),
		config.RetryDelay(),
	)

	return &Orchestrator{
		client:         deps.Client,
		tournamentRepo: deps.TournamentRepo,
		teamRepo:       deps.TeamRepo,
		playerRepo:     deps.PlayerRepo,
		playerTeamRepo: deps.PlayerTeamRepo,
		retryManager:   retryManager,
		config:         config,
	}
}
