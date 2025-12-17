package stats

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior/stats/detailed"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/junior/stats/logger"
	"go.uber.org/zap"
)

// Проверка реализации интерфейса
var _ parser.StatsParserService = (*service)(nil)

type service struct {
	parser         *stats.Parser
	repo           player_statistics.Repository
	statsLogger    logger.StatsLogger
	zapLogger      *zap.Logger
	detailedParser *detailed.Parser
}

// NewStatsParserService создает новый сервис парсинга статистики
func NewStatsParserService(
	parser *stats.Parser,
	repo player_statistics.Repository,
	zapLogger *zap.Logger,
) *service {
	statsLogger := logger.New(zapLogger)

	// Создаём детальный парсер
	detailedParser := detailed.NewParser(
		repo,
		statsLogger,
		zapLogger,
		convertOne, // функция конвертации из converter.go
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
// Возвращает количество обработанных записей
func (s *service) ParseTournamentStats(
	ctx context.Context,
	domain string,
	tournamentURL string,
	tournamentID string,
) (int, error) {
	// Используем детальный парсер
	return s.detailedParser.ParseTournamentStats(ctx, domain, tournamentURL, tournamentID)
}
