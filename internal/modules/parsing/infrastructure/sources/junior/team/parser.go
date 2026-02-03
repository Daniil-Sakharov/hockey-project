package team

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/helpers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
)

const MinBirthYear = 2008

// Parser парсер команд
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер команд
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// ParseFromTournament парсит команды из турнира с контекстом года/группы.
// fallbackBirthYears — годы рождения со страницы списка турниров, используются
// если на странице команд нет dropdown года.
func (p *Parser) ParseFromTournament(ctx context.Context, domain, tournamentURL string, fallbackBirthYears ...int) ([]types.TeamWithContext, error) {
	teamsURL := domain + tournamentURL
	if !strings.HasSuffix(teamsURL, "/") {
		teamsURL += "/"
	}
	teamsURL += "teams/"

	logger.Info(ctx, fmt.Sprintf("  🏒 Загрузка страницы команд: %s", teamsURL))

	resp, err := p.http.MakeRequest(teamsURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки страницы команд: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	yearLinks := helpers.ExtractYearLinks(doc)
	initialGroups := helpers.ExtractGroupLinks(doc)

	var teamsWithContext []types.TeamWithContext

	if len(yearLinks) == 0 && len(initialGroups) == 0 {
		logger.Info(ctx, "     ℹ️  Нет переключателей года/группы, парсим основную страницу")
		teamsMap := make(map[string]types.TeamDTO)
		helpers.ParseTeamsFromDocWithDomain(doc, teamsMap, domain)

		// Если есть fallback год — назначаем его командам
		var birthYear *int
		if len(fallbackBirthYears) == 1 {
			birthYear = &fallbackBirthYears[0]
			logger.Info(ctx, fmt.Sprintf("     📅 Используем fallback год рождения: %d", *birthYear))
		}

		for _, team := range teamsMap {
			teamsWithContext = append(teamsWithContext, types.TeamWithContext{
				Team:      team,
				BirthYear: birthYear,
				GroupName: nil,
			})
		}
		logger.Info(ctx, fmt.Sprintf("     💾 Найдено команд: %d", len(teamsWithContext)))
	} else if len(yearLinks) == 0 && len(initialGroups) > 0 {
		// Нет годов, но есть группы — используем fallback год если есть
		var birthYear *int
		if len(fallbackBirthYears) == 1 {
			birthYear = &fallbackBirthYears[0]
			logger.Info(ctx, fmt.Sprintf("     📅 Группы без годов, используем fallback год: %d", *birthYear))
		}
		teamsWithContext = p.parseWithFiltersAndFallbackYear(ctx, domain, doc, initialGroups, birthYear)
	} else {
		teamsWithContext = p.parseWithFilters(ctx, domain, doc, yearLinks, initialGroups)
	}

	return teamsWithContext, nil
}
