package stats

import (
	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
)

type Dependencies struct {
	Client               *fhspb.Client
	TournamentRepo       *fhspbrepo.TournamentRepository
	TeamRepo             *fhspbrepo.TeamRepository
	PlayerRepo           *fhspbrepo.PlayerRepository
	PlayerStatisticsRepo *fhspbrepo.PlayerStatisticsRepository
	GoalieStatisticsRepo *fhspbrepo.GoalieStatisticsRepository
	StatisticsWorkers    int
}

type Service struct {
	deps         Dependencies
	orchestrator *Orchestrator
}

func NewService(deps Dependencies) *Service {
	return &Service{
		deps:         deps,
		orchestrator: NewOrchestrator(deps),
	}
}

func (s *Service) Orchestrator() *Orchestrator {
	return s.orchestrator
}
