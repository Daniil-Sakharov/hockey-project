package stats

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Parser клиент для парсинга статистики турниров
type Parser struct {
	httpClient *http.Client
}

// NewParser создает новый парсер статистики
func NewParser(httpClient *http.Client) *Parser {
	return &Parser{
		httpClient: httpClient,
	}
}

// ParseTournamentStats парсит всю статистику турнира
func (p *Parser) ParseTournamentStats(
	ctx context.Context,
	domain string,
	tournamentURL string,
	tournamentID string,
) ([]PlayerStatisticDTO, error) {
	statsURL := fmt.Sprintf("%s%sstats/", domain, tournamentURL)
	resp, err := p.httpClient.Get(statsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 3. Парсим HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// 4. Извлекаем комбинации год+группа с AJAX запросами для каждого года
	combinations, err := ParseCombinationsWithAjax(ctx, doc, domain, p.httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to parse combinations: %w", err)
	}

	if len(combinations) == 0 {
		// Логирование происходит в service layer
		return []PlayerStatisticDTO{}, nil
	}

	// Логирование происходит в service layer через statsLogger

	// 5. Парсим статистику для каждой комбинации
	var allStats []PlayerStatisticDTO

	for _, combo := range combinations {
		// Логирование происходит в service layer через statsLogger
		stats, err := p.fetchCombinationStats(ctx, domain, tournamentID, combo)
		if err != nil {
			// Логирование происходит в service layer
			continue
		}

		allStats = append(allStats, stats...)
	}

	// Логирование происходит в service layer через statsLogger
	return allStats, nil
}

// fetchCombinationStats парсит статистику для конкретной комбинации год+группа
func (p *Parser) fetchCombinationStats(
	ctx context.Context,
	domain string,
	tournamentID string,
	combo StatsCombination,
) ([]PlayerStatisticDTO, error) {
	// Делаем запрос к JSON API (season пустой - используется fallback)
	statsResp, err := FetchStatistics(
		ctx,
		p.httpClient,
		domain,
		tournamentID,
		combo.YearID,
		combo.GroupID,
		"", // season - будет использован fallback
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch statistics: %w", err)
	}

	// Добавляем контекст (group_name, birth_year) в каждый DTO
	for i := range statsResp.Data {
		statsResp.Data[i].GroupName = combo.GroupName
		statsResp.Data[i].BirthYear = combo.YearLabel
	}

	return statsResp.Data, nil
}
