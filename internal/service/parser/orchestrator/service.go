package orchestrator

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
)

// Константы конфигурации Worker Pools
const (
	domainWorkers = 5 // Worker Pool для доменов (было 3)
)

type orchestratorService struct {
	juniorService  parser.JuniorParserService
	playerRepo     player.Repository
	teamRepo       team.Repository
	tournamentRepo tournament.Repository
	playerTeamRepo player_team.Repository
}

// NewOrchestratorService создает orchestrator для парсинга
func NewOrchestratorService(
	juniorService parser.JuniorParserService,
	playerRepo player.Repository,
	teamRepo team.Repository,
	tournamentRepo tournament.Repository,
	playerTeamRepo player_team.Repository,
) *orchestratorService {
	return &orchestratorService{
		juniorService:  juniorService,
		playerRepo:     playerRepo,
		teamRepo:       teamRepo,
		tournamentRepo: tournamentRepo,
		playerTeamRepo: playerTeamRepo,
	}
}
