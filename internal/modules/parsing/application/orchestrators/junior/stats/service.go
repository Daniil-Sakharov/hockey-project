package stats

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/interfaces"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats/detailed"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/stats/logger"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
	"go.uber.org/zap"
)

// Проверка реализации интерфейса
var _ interfaces.StatsParserService = (*service)(nil)

type service struct {
	parser         *stats.Parser
	repo           detailed.Repository
	statsLogger    logger.StatsLogger
	zapLogger      *zap.Logger
	detailedParser *detailed.Parser
}

// NewStatsParserService создает новый сервис парсинга статистики
func NewStatsParserService(
	parser *stats.Parser,
	repo detailed.Repository,
	zapLogger *zap.Logger,
) *service {
	statsLogger := logger.New(zapLogger)

	detailedParser := detailed.NewParser(
		repo,
		statsLogger,
		zapLogger,
		convertOne,
	)

	return &service{
		parser:         parser,
		repo:           repo,
		statsLogger:    statsLogger,
		zapLogger:      zapLogger,
		detailedParser: detailedParser,
	}
}

// ParseTournamentStats парсит статистику турнира и сохраняет в БД
func (s *service) ParseTournamentStats(
	ctx context.Context,
	domain string,
	tournamentURL string,
	tournamentID string,
	season string,
) (int, error) {
	return s.detailedParser.ParseTournamentStats(ctx, domain, tournamentURL, tournamentID, season)
}
