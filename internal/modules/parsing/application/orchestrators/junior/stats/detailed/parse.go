package detailed

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/stats"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

// ParseTournamentStats парсит детальную статистику турнира
func (p *Parser) ParseTournamentStats(
	ctx context.Context,
	domain string,
	tournamentURL string,
	tournamentID string,
	season string,
) (int, error) {
	// Конвертируем формат сезона: "2023/2024" -> "2023-2024"
	season = strings.ReplaceAll(season, "/", "-")

	var statsURL string
	if strings.HasPrefix(tournamentURL, "http") {
		statsURL = tournamentURL + "/stats/"
	} else {
		statsURL = fmt.Sprintf("%s%sstats/", domain, tournamentURL)
	}

	p.statsLogger.LogTournamentStart(tournamentID, "Tournament", statsURL)

	resp, err := http.Get(statsURL) //nolint:gosec
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stats page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to parse HTML: %w", err)
	}

	p.zapLogger.Info("🔍 Парсинг комбинаций", zap.String("tournament_id", tournamentID), zap.String("season", season))

	combinations, err := stats.ParseCombinationsWithAjax(ctx, doc, domain, http.DefaultClient)
	if err != nil {
		p.zapLogger.Error("❌ Ошибка парсинга комбинаций", zap.Error(err))
		return 0, fmt.Errorf("failed to parse combinations: %w", err)
	}

	if len(combinations) == 0 {
		p.zapLogger.Warn("⚠️ Комбинации год+группа НЕ найдены", zap.String("tournament_id", tournamentID))
		return 0, nil
	}

	p.zapLogger.Info("✅ Найдено комбинаций", zap.Int("count", len(combinations)), zap.String("tournament_id", tournamentID))

	allDTOs, totalReceivedFromAPI := p.fetchAllCombinations(ctx, domain, tournamentID, season, combinations)

	p.zapLogger.Info("📊 Собрано записей из API", zap.Int("total", len(allDTOs)), zap.Int("from_api", totalReceivedFromAPI))

	if len(allDTOs) == 0 {
		p.statsLogger.LogTournamentSummary(0, 0)
		return 0, nil
	}

	result, conversionLosses := ConvertWithTracking(allDTOs, tournamentID, p.convertOne)

	if err := p.deleteOldStats(ctx, tournamentID); err != nil {
		return 0, err
	}

	savedCount, savingLosses := p.saveStatsEntities(ctx, result)

	p.logLosses(conversionLosses, savingLosses)

	p.statsLogger.LogTournamentSummary(totalReceivedFromAPI, savedCount)

	return savedCount, nil
}

// saveStatsEntities сохраняет статистику в БД
func (p *Parser) saveStatsEntities(
	ctx context.Context,
	result []*entities.PlayerStatistic,
) (int, []MissingPlayerInfo) {
	savedCount, err := p.repo.CreateBatch(ctx, result)
	if err != nil {
		p.zapLogger.Warn("Ошибка батча, сохраняем по одной")
		return SaveOneByOne(ctx, p.repo, result)
	}

	return savedCount, nil
}
