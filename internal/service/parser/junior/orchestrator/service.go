package orchestrator

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
)

type orchestratorService struct {
	juniorService  parser.JuniorParserService
	playerRepo     player.Repository
	teamRepo       team.Repository
	tournamentRepo tournament.Repository
	playerTeamRepo player_team.Repository
	config         config.JuniorConfig
}

// NewOrchestratorService создает orchestrator для парсинга
func NewOrchestratorService(
	juniorService parser.JuniorParserService,
	playerRepo player.Repository,
	teamRepo team.Repository,
	tournamentRepo tournament.Repository,
	playerTeamRepo player_team.Repository,
	cfg config.JuniorConfig,
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
