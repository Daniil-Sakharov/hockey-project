package orchestrator

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	fhspbRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
	svc "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/retry"
)

// Проверка реализации интерфейса
var _ svc.Service = (*Orchestrator)(nil)

// Orchestrator оркестратор парсинга fhspb.ru
type Orchestrator struct {
	client         *fhspb.Client
	tournamentRepo *fhspbRepo.TournamentRepository
	teamRepo       *fhspbRepo.TeamRepository
	playerRepo     *fhspbRepo.PlayerRepository
	playerTeamRepo *fhspbRepo.PlayerTeamRepository
	retryManager   *retry.Manager
	config         svc.Config
}

// New создает новый оркестратор
func New(deps svc.Dependencies, config svc.Config) *Orchestrator {
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
