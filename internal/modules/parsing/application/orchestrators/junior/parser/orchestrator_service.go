package parser

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/interfaces"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
)

type orchestratorService struct {
	juniorService  JuniorParserService
	playerRepo     repositories.PlayerRepository
	teamRepo       repositories.TeamRepository
	tournamentRepo repositories.TournamentRepository
	playerTeamRepo repositories.PlayerTeamRepository
	config         interfaces.JuniorConfig
}

// NewOrchestratorService создает orchestrator для парсинга
func NewOrchestratorService(
	juniorService JuniorParserService,
	playerRepo repositories.PlayerRepository,
	teamRepo repositories.TeamRepository,
	tournamentRepo repositories.TournamentRepository,
	playerTeamRepo repositories.PlayerTeamRepository,
	cfg interfaces.JuniorConfig,
) *orchestratorService {
	return &orchestratorService{
		juniorService:  juniorService,
		playerRepo:     playerRepo,
		teamRepo:       teamRepo,
		tournamentRepo: tournamentRepo,
		playerTeamRepo: playerTeamRepo,
		config:         cfg,
	}
}

// Placeholder for entities usage
var _ = entities.Player{}
