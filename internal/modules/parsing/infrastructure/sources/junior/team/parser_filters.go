package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/helpers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
)

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
				skippedYears++
				continue
			}

			yearDisplay := "неизвестный"
			if yearLink.Year > 0 {
				yearDisplay = fmt.Sprintf("%d", yearLink.Year)
			}
			logger.Info(ctx, fmt.Sprintf("     🗓️  [%d/%d] Обработка года %s...", yearIdx+1, len(yearLinks), yearDisplay))

			p.processYear(ctx, domain, yearLink, teamsMap, &totalCombinations)
		}
	}

	logger.Info(ctx, fmt.Sprintf("     📊 Обработано комбинаций: %d, пропущено годов: %d", totalCombinations, skippedYears))
	logger.Info(ctx, fmt.Sprintf("     💾 Итого уникальных команд: %d", len(teamsMap)))
}

func (p *Parser) processYear(ctx context.Context, domain string, yearLink types.YearLink, teamsMap map[string]types.TeamDTO, totalCombinations *int) {
	fullYearURL := domain + yearLink.AjaxURL
	yearResp, err := p.http.MakeRequest(fullYearURL)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("        ⚠️  Ошибка запроса года: %v", err))
		return
	}

	if yearResp.StatusCode != 200 {
		_ = yearResp.Body.Close()
		return
	}

	yearDoc, err := goquery.NewDocumentFromReader(yearResp.Body)
	_ = yearResp.Body.Close()
	if err != nil {
		return
	}

	groupLinks := helpers.ExtractGroupLinks(yearDoc)

	if len(groupLinks) == 0 {
		beforeCount := len(teamsMap)
		helpers.ParseTeamsFromDoc(yearDoc, teamsMap)
		newCount := len(teamsMap) - beforeCount
		*totalCombinations++
		if newCount > 0 {
			logger.Info(ctx, fmt.Sprintf("        ✅ Год без групп: +%d команд", newCount))
		}
	} else {
		for _, groupURL := range groupLinks {
			p.parseTeamsFromAjax(ctx, domain, groupURL, teamsMap)
			*totalCombinations++
		}
	}
}

func (p *Parser) parseTeamsFromAjax(ctx context.Context, domain, ajaxURL string, teamsMap map[string]types.TeamDTO) {
	fullURL := domain + ajaxURL

	ajaxResp, err := p.http.MakeRequest(fullURL)
	if err != nil {
		return
	}

	if ajaxResp.StatusCode != 200 {
		_ = ajaxResp.Body.Close()
		return
	}

	ajaxDoc, err := goquery.NewDocumentFromReader(ajaxResp.Body)
	_ = ajaxResp.Body.Close()
	if err != nil {
		return
	}

	beforeCount := len(teamsMap)
	helpers.ParseTeamsFromDoc(ajaxDoc, teamsMap)
	newCount := len(teamsMap) - beforeCount

	if newCount > 0 {
		logger.Info(ctx, fmt.Sprintf("              ✅ +%d команд (всего: %d)", newCount, len(teamsMap)))
	}
}
