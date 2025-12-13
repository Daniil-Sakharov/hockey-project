package orchestrator

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	svc "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb"
)

// Проверка реализации интерфейса
var _ svc.Service = (*Orchestrator)(nil)

// Orchestrator оркестратор парсинга fhspb.ru
type Orchestrator struct {
	client         *fhspb.Client
	playerRepo     player.Repository
	teamRepo       team.Repository
	tournamentRepo tournament.Repository
	config         svc.Config
}

// New создает новый оркестратор
func New(deps svc.Dependencies, config svc.Config) *Orchestrator {
	return &Orchestrator{
		client:         deps.Client,
		playerRepo:     deps.PlayerRepo,
		teamRepo:       deps.TeamRepo,
		tournamentRepo: deps.TournamentRepo,
		config:         config,
	}
}
