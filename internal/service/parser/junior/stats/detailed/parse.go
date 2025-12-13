package detailed

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/stats"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	"github.com/PuerkitoBio/goquery"
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
	convertOne  func(stats.PlayerStatisticDTO, string) (*player_statistics.PlayerStatistic, error)
}

// NewParser создаёт новый парсер
func NewParser(
	repo Repository,
	statsLogger StatsLogger,
	zapLogger *zap.Logger,
	convertOne func(stats.PlayerStatisticDTO, string) (*player_statistics.PlayerStatistic, error),
) *Parser {
	return &Parser{
		repo:        repo,
		statsLogger: statsLogger,
		zapLogger:   zapLogger,
		convertOne:  convertOne,
	}
}

// ParseTournamentStats парсит детальную статистику турнира с подробным логированием
// Использует двухуровневый парсинг: сначала года, потом группы для каждого года
func (p *Parser) ParseTournamentStats(
	ctx context.Context,
	domain string,
	tournamentURL string,
	tournamentID string,
) (int, error) {
	statsURL := fmt.Sprintf("%s%sstats/", domain, tournamentURL)

	// Логируем начало парсинга турнира
	p.statsLogger.LogTournamentStart(tournamentID, "Tournament", statsURL)

	// 1. Загружаем HTML страницу
	resp, err := http.Get(statsURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch stats page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 2. Парсим HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// 3. Извлекаем комбинации год+группа с двухуровневой логикой
	// Для каждого года делаем AJAX запрос чтобы получить ВСЕ группы
	combinations, err := stats.ParseCombinationsWithAjax(ctx, doc, domain, http.DefaultClient)
	if err != nil {
		return 0, fmt.Errorf("failed to parse combinations: %w", err)
	}

	if len(combinations) == 0 {
		p.zapLogger.Warn("Комбинации год+группа НЕ найдены")
		return 0, nil
	}

	p.zapLogger.Info("Найдено комбинаций (двухуровневый парсинг)",
		zap.Int("count", len(combinations)))

	// 4. Парсим статистику для каждой комбинации с детальным логированием
	allDTOs, totalReceivedFromAPI := p.fetchAllCombinations(ctx, domain, tournamentID, combinations)

	p.zapLogger.Info("Собрано записей из API",
		zap.Int("total", len(allDTOs)))

	if len(allDTOs) == 0 {
		p.statsLogger.LogTournamentSummary(0, 0)
		return 0, nil
	}

	// 5. Конвертируем DTO → Domain entities с отслеживанием потерь
	entities, conversionLosses := ConvertWithTracking(allDTOs, tournamentID, p.convertOne)

	// 6. Удаляем старую статистику турнира
	if err := p.deleteOldStats(ctx, tournamentID); err != nil {
		return 0, err
	}

	// 7. Сохраняем в БД
	savedCount, savingLosses := p.saveStats(ctx, entities)

	// 8. Логируем потери
	p.logLosses(conversionLosses, savingLosses)

	// 9. Финальный итог по турниру
	p.statsLogger.LogTournamentSummary(totalReceivedFromAPI, savedCount)

	return savedCount, nil
}

// fetchAllCombinations загружает статистику для всех комбинаций год+группа
func (p *Parser) fetchAllCombinations(
	ctx context.Context,
	domain string,
	tournamentID string,
	combinations []stats.StatsCombination,
) ([]stats.PlayerStatisticDTO, int) {
	totalReceivedFromAPI := 0
	var allDTOs []stats.PlayerStatisticDTO

	for _, combo := range combinations {
		p.statsLogger.LogCombinationStart(
			combo.YearLabel,
			combo.YearID,
			combo.GroupName,
			combo.GroupID,
		)

		// Делаем запрос к API
		statsResp, err := stats.FetchStatistics(
			ctx,
			http.DefaultClient,
			domain,
			tournamentID,
			combo.YearID,
			combo.GroupID,
		)
		if err != nil {
			p.statsLogger.LogCombinationError(err)
			continue
		}

		receivedCount := len(statsResp.Data)
		totalReceivedFromAPI += receivedCount

		// Извлекаем player_id для логирования
		playerIDs := make([]string, 0, receivedCount)
		for _, dto := range statsResp.Data {
			playerID := stats.ExtractPlayerID(dto.Surname)
			if playerID != "" {
				playerIDs = append(playerIDs, playerID)
			}
		}

		// Логируем результат
		p.statsLogger.LogCombinationResult(receivedCount, playerIDs)

		// Добавляем контекст к DTO
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
		if err := repo.DeleteByTournament(ctx, tournamentID); err != nil {
			return fmt.Errorf("failed to delete old statistics: %w", err)
		}
	}

	return nil
}

// saveStats сохраняет статистику в БД с обработкой ошибок
func (p *Parser) saveStats(
	ctx context.Context,
	entities []*player_statistics.PlayerStatistic,
) (int, []MissingPlayerInfo) {
	savedCount, err := p.repo.CreateBatch(ctx, entities)

	if err != nil {
		// Если ошибка FK - пробуем сохранить по одной
		p.zapLogger.Warn("Ошибка батча, сохраняем по одной")
		return SaveOneByOne(ctx, p.repo, entities)
	}

	return savedCount, nil
}

// logLosses логирует все потери данных
func (p *Parser) logLosses(
	conversionLosses []MissingPlayerInfo,
	savingLosses []MissingPlayerInfo,
) {
	// Логируем потери при конвертации
	for _, loss := range conversionLosses {
		p.statsLogger.LogValidationSkip(loss.PlayerID, loss.TeamID, loss.Reason)
	}

	// Логируем потери при сохранении
	for _, loss := range savingLosses {
		p.statsLogger.LogFKConstraintSkip(loss.PlayerID, loss.TeamID, loss.Reason)
	}
}
