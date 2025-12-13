package stats_orchestrator

import (
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
)

// Проверка реализации интерфейса
var _ parser.StatsOrchestratorService = (*service)(nil)

type service struct {
	statsService   parser.StatsParserService
	statsRepo      player_statistics.Repository
	tournamentRepo tournament.Repository
	logger         *log.Logger
}

// NewStatsOrchestratorService создает новый orchestrator для парсинга статистики
func NewStatsOrchestratorService(
	statsService parser.StatsParserService,
	statsRepo player_statistics.Repository,
	tournamentRepo tournament.Repository,
	logger *log.Logger,
) *service {
	return &service{
		statsService:   statsService,
		statsRepo:      statsRepo,
		tournamentRepo: tournamentRepo,
		logger:         logger,
	}
}
