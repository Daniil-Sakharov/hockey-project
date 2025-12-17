package di

import (
	"context"
	"log"
	"os"

	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
	profileService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/profile"
	reportService "github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/report"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
	juniorService "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior"
	juniorOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior/orchestrator"
	statsService "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior/stats"
	statsOrchestrator "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior/stats/orchestrator"
	"go.uber.org/zap"
)

// Service содержит все сервисы
type Service struct {
	config *config.Config
	infra  *Infrastructure
	repo   *Repository

	juniorParserService      parser.JuniorParserService
	statsParserService       parser.StatsParserService
	orchestratorService      parser.OrchestratorService
	statsOrchestratorService parser.StatsOrchestratorService
	stateManager             bot.StateManager
	searchPlayerService      bot.SearchPlayerService
	profileService           bot.ProfileService
	reportService            *reportService.Service
}

func NewService(cfg *config.Config, infra *Infrastructure, repo *Repository) *Service {
	return &Service{config: cfg, infra: infra, repo: repo}
}

func (s *Service) JuniorParser(ctx context.Context) parser.JuniorParserService {
	if s.juniorParserService == nil {
		s.juniorParserService = juniorService.NewJuniorService(s.infra.JuniorClient())
	}
	return s.juniorParserService
}

func (s *Service) Orchestrator(ctx context.Context) parser.OrchestratorService {
	if s.orchestratorService == nil {
		s.orchestratorService = juniorOrchestrator.NewOrchestratorService(
			s.JuniorParser(ctx),
			s.repo.Player(ctx),
			s.repo.Team(ctx),
			s.repo.Tournament(ctx),
			s.repo.PlayerTeam(ctx),
			s.config.Junior,
		)
	}
	return s.orchestratorService
}

func (s *Service) StatsParser(ctx context.Context) parser.StatsParserService {
	if s.statsParserService == nil {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapLogger, err := zapConfig.Build()
		if err != nil {
			log.Fatalf("Failed to create zap logger: %v", err)
		}
		zapLogger = zapLogger.With(zap.String("service", "stats_parser"))
		s.statsParserService = statsService.NewStatsParserService(s.infra.StatsParser(), s.repo.PlayerStatistics(ctx), zapLogger)
	}
	return s.statsParserService
}

func (s *Service) StatsOrchestrator(ctx context.Context) parser.StatsOrchestratorService {
	if s.statsOrchestratorService == nil {
		logger := log.New(os.Stdout, "[STATS] ", log.LstdFlags|log.Lmsgprefix)
		s.statsOrchestratorService = statsOrchestrator.NewStatsOrchestratorService(
			s.StatsParser(ctx),
			s.repo.PlayerStatistics(ctx),
			s.repo.Tournament(ctx),
			logger,
		)
	}
	return s.statsOrchestratorService
}

func (s *Service) StateManager() bot.StateManager {
	if s.stateManager == nil {
		s.stateManager = bot.NewStateManager()
	}
	return s.stateManager
}

func (s *Service) SearchPlayer(ctx context.Context) bot.SearchPlayerService {
	if s.searchPlayerService == nil {
		s.searchPlayerService = bot.NewSearchPlayerService(s.repo.Player(ctx))
	}
	return s.searchPlayerService
}

func (s *Service) Profile(ctx context.Context) bot.ProfileService {
	if s.profileService == nil {
		s.profileService = profileService.NewService(
			s.repo.Player(ctx),
			s.repo.PlayerStatistics(ctx),
			s.repo.PlayerTeam(ctx),
			s.repo.Team(ctx),
		)
	}
	return s.profileService
}

func (s *Service) Report(ctx context.Context) *reportService.Service {
	if s.reportService == nil {
		s.reportService = reportService.NewService(
			s.repo.Player(ctx),
			s.repo.PlayerStatistics(ctx),
			s.repo.PlayerTeam(ctx),
			s.repo.Team(ctx),
			s.repo.Tournament(ctx),
		)
	}
	return s.reportService
}
