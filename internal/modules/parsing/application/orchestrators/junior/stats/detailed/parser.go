package detailed

import (
	"context"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
	"go.uber.org/zap"
)

// StatsLogger интерфейс для логирования статистики
type StatsLogger interface {
	LogTournamentStart(tournamentID, tournamentType, url string)
	LogCombinationStart(yearLabel, yearID, groupName, groupID string)
	LogCombinationError(err error)
	LogCombinationResult(receivedCount int, playerIDs []string)
	LogTournamentSummary(totalReceivedFromAPI, savedCount int)
	LogValidationSkip(playerID, teamID, reason string)
	LogFKConstraintSkip(playerID, teamID, reason string)
}

// Parser парсер детальной статистики турнира
type Parser struct {
	repo        Repository
	statsLogger StatsLogger
	zapLogger   *zap.Logger
	convertOne  func(stats.PlayerStatisticDTO, string) (*entities.PlayerStatistic, error)
}

// NewParser создаёт новый парсер
func NewParser(
	repo Repository,
	statsLogger StatsLogger,
	zapLogger *zap.Logger,
	convertOne func(stats.PlayerStatisticDTO, string) (*entities.PlayerStatistic, error),
) *Parser {
	return &Parser{
		repo:        repo,
		statsLogger: statsLogger,
		zapLogger:   zapLogger,
		convertOne:  convertOne,
	}
}

// fetchAllCombinations загружает статистику для всех комбинаций
func (p *Parser) fetchAllCombinations(
	ctx context.Context,
	domain string,
	tournamentID string,
	season string,
	combinations []stats.StatsCombination,
) ([]stats.PlayerStatisticDTO, int) {
	totalReceivedFromAPI := 0
	var allDTOs []stats.PlayerStatisticDTO

	p.zapLogger.Info("🔄 Начинаю загрузку комбинаций",
		zap.String("tournament_id", tournamentID),
		zap.String("season", season),
		zap.Int("combinations_count", len(combinations)))

	for i, combo := range combinations {
		p.zapLogger.Debug("📥 Запрос комбинации",
			zap.Int("index", i+1),
			zap.String("year", combo.YearLabel),
			zap.String("group", combo.GroupName),
			zap.String("season", season))

		p.statsLogger.LogCombinationStart(combo.YearLabel, combo.YearID, combo.GroupName, combo.GroupID)

		statsResp, err := stats.FetchStatistics(ctx, http.DefaultClient, domain, tournamentID, combo.YearID, combo.GroupID, season)
		if err != nil {
			p.zapLogger.Warn("❌ Ошибка запроса", zap.Error(err))
			p.statsLogger.LogCombinationError(err)
			continue
		}

		receivedCount := len(statsResp.Data)
		totalReceivedFromAPI += receivedCount

		p.zapLogger.Debug("✅ Получено записей", zap.Int("count", receivedCount))

		playerIDs := make([]string, 0, receivedCount)
		for _, dto := range statsResp.Data {
			if playerID := stats.ExtractPlayerID(dto.Surname); playerID != "" {
				playerIDs = append(playerIDs, playerID)
			}
		}

		p.statsLogger.LogCombinationResult(receivedCount, playerIDs)

		for i := range statsResp.Data {
			statsResp.Data[i].GroupName = combo.GroupName
			statsResp.Data[i].BirthYear = combo.YearLabel
		}

		allDTOs = append(allDTOs, statsResp.Data...)
	}

	return allDTOs, totalReceivedFromAPI
}

// deleteOldStats удаляет старую статистику турнира
func (p *Parser) deleteOldStats(ctx context.Context, tournamentID string) error {
	type deleter interface {
		DeleteByTournament(ctx context.Context, tournamentID string) error
	}

	if repo, ok := p.repo.(deleter); ok {
		return repo.DeleteByTournament(ctx, tournamentID)
	}

	return nil
}

// logLosses логирует все потери данных
func (p *Parser) logLosses(conversionLosses, savingLosses []MissingPlayerInfo) {
	for _, loss := range conversionLosses {
		p.statsLogger.LogValidationSkip(loss.PlayerID, loss.TeamID, loss.Reason)
	}

	for _, loss := range savingLosses {
		p.statsLogger.LogFKConstraintSkip(loss.PlayerID, loss.TeamID, loss.Reason)
	}
}
