package stats

import (
	"log"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/interfaces"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
)

// Проверка реализации интерфейса
var _ interfaces.StatsOrchestratorService = (*service)(nil)

type service struct {
	statsService   interfaces.StatsParserService
	statsRepo      repositories.PlayerStatisticsRepository
	tournamentRepo repositories.TournamentRepository
	logger         *log.Logger
}

// NewStatsOrchestratorService создает новый orchestrator для парсинга статистики
func NewStatsOrchestratorService(
	statsService interfaces.StatsParserService,
	statsRepo repositories.PlayerStatisticsRepository,
	tournamentRepo repositories.TournamentRepository,
	logger *log.Logger,
) *service {
	return &service{
		statsService:   statsService,
		statsRepo:      statsRepo,
		tournamentRepo: tournamentRepo,
		logger:         logger,
	}
}
