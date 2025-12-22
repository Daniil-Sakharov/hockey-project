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

// ParseFromTournament парсит команды из турнира
func (p *Parser) ParseFromTournament(ctx context.Context, domain, tournamentURL string) ([]types.TeamDTO, error) {
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

	teamsMap := make(map[string]types.TeamDTO)

	yearLinks := helpers.ExtractYearLinks(doc)
	initialGroups := helpers.ExtractGroupLinks(doc)

	if len(yearLinks) == 0 && len(initialGroups) == 0 {
		logger.Info(ctx, "     ℹ️  Нет переключателей года/группы, парсим основную страницу")
		helpers.ParseTeamsFromDoc(doc, teamsMap)
		logger.Info(ctx, fmt.Sprintf("     💾 Найдено команд: %d", len(teamsMap)))
	} else {
		p.parseWithFilters(ctx, domain, doc, yearLinks, initialGroups, teamsMap)
	}

	teams := make([]types.TeamDTO, 0, len(teamsMap))
	for _, team := range teamsMap {
		teams = append(teams, team)
	}
	return teams, nil
}
