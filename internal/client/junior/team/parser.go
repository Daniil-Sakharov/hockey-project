package team

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/helpers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/types"
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

func (p *Parser) parseWithFilters(ctx context.Context, domain string, doc *goquery.Document, yearLinks []types.YearLink, initialGroups []string, teamsMap map[string]types.TeamDTO) {
	totalCombinations := 0
	skippedYears := 0

	if len(yearLinks) == 0 {
		logger.Info(ctx, fmt.Sprintf("     📅 Годов нет, найдено %d групп", len(initialGroups)))
		for _, groupURL := range initialGroups {
			p.parseTeamsFromAjax(ctx, domain, groupURL, teamsMap)
			totalCombinations++
		}
	} else {
		logger.Info(ctx, fmt.Sprintf("     📅 Найдено %d годов рождения", len(yearLinks)))

		for yearIdx, yearLink := range yearLinks {
			if yearLink.Year > 0 && yearLink.Year < MinBirthYear {
				logger.Info(ctx, fmt.Sprintf("     ⏭️  [%d/%d] Пропускаем год %d (< %d)",
					yearIdx+1, len(yearLinks), yearLink.Year, MinBirthYear))
				skippedYears++
				continue
			}

			yearDisplay := "неизвестный"
			if yearLink.Year > 0 {
				yearDisplay = fmt.Sprintf("%d", yearLink.Year)
			}
			logger.Info(ctx, fmt.Sprintf("     🗓️  [%d/%d] Обработка года %s...", yearIdx+1, len(yearLinks), yearDisplay))

			fullYearURL := domain + yearLink.AjaxURL
			yearResp, err := p.http.MakeRequest(fullYearURL)
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("        ⚠️  Ошибка запроса года: %v", err))
				continue
			}

			if yearResp.StatusCode != 200 {
				_ = yearResp.Body.Close()
				logger.Warn(ctx, fmt.Sprintf("        ⚠️  HTTP %d для года", yearResp.StatusCode))
				continue
			}

			yearDoc, err := goquery.NewDocumentFromReader(yearResp.Body)
			_ = yearResp.Body.Close()
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("        ⚠️  Ошибка парсинга года: %v", err))
				continue
			}

			groupLinks := helpers.ExtractGroupLinks(yearDoc)

			if len(groupLinks) == 0 {
				beforeCount := len(teamsMap)
				helpers.ParseTeamsFromDoc(yearDoc, teamsMap)
				newCount := len(teamsMap) - beforeCount
				totalCombinations++
				if newCount > 0 {
					logger.Info(ctx, fmt.Sprintf("        ✅ Год без групп: +%d команд (всего: %d)", newCount, len(teamsMap)))
				}
			} else {
				logger.Info(ctx, fmt.Sprintf("        📁 Найдено %d групп для этого года", len(groupLinks)))
				for groupIdx, groupURL := range groupLinks {
					logger.Info(ctx, fmt.Sprintf("           [%d/%d] Группа...", groupIdx+1, len(groupLinks)))
					p.parseTeamsFromAjax(ctx, domain, groupURL, teamsMap)
					totalCombinations++
				}
			}
		}
	}

	logger.Info(ctx, fmt.Sprintf("     📊 Обработано комбинаций: %d, пропущено годов: %d", totalCombinations, skippedYears))
	logger.Info(ctx, fmt.Sprintf("     💾 Итого уникальных команд: %d", len(teamsMap)))
}

func (p *Parser) parseTeamsFromAjax(ctx context.Context, domain, ajaxURL string, teamsMap map[string]types.TeamDTO) {
	fullURL := domain + ajaxURL

	ajaxResp, err := p.http.MakeRequest(fullURL)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("              ⚠️  Ошибка запроса: %v", err))
		return
	}

	if ajaxResp.StatusCode != 200 {
		_ = ajaxResp.Body.Close()
		logger.Warn(ctx, fmt.Sprintf("              ⚠️  HTTP %d", ajaxResp.StatusCode))
		return
	}

	ajaxDoc, err := goquery.NewDocumentFromReader(ajaxResp.Body)
	_ = ajaxResp.Body.Close()
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("              ⚠️  Ошибка парсинга: %v", err))
		return
	}

	beforeCount := len(teamsMap)
	helpers.ParseTeamsFromDoc(ajaxDoc, teamsMap)
	newCount := len(teamsMap) - beforeCount

	if newCount > 0 {
		logger.Info(ctx, fmt.Sprintf("              ✅ +%d команд (всего: %d)", newCount, len(teamsMap)))
	}
}
