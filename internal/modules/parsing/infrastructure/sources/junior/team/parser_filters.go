package team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/helpers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
)

// parseWithFiltersAndFallbackYear обрабатывает группы без yearLinks, назначая fallback год
func (p *Parser) parseWithFiltersAndFallbackYear(ctx context.Context, domain string, doc *goquery.Document, groups []types.GroupLink, fallbackYear *int) []types.TeamWithContext {
	var teamsWithContext []types.TeamWithContext
	totalCombinations := 0

	logger.Info(ctx, fmt.Sprintf("     📅 Годов нет, найдено %d групп", len(groups)))

	activeGroupName := helpers.ExtractActiveGroupName(doc, "competitions-teams")
	if activeGroupName != "" {
		teamsMap := make(map[string]types.TeamDTO)
		helpers.ParseTeamsFromDocWithDomain(doc, teamsMap, domain)

		for _, team := range teamsMap {
			gn := activeGroupName
			teamsWithContext = append(teamsWithContext, types.TeamWithContext{
				Team:      team,
				BirthYear: fallbackYear,
				GroupName: &gn,
			})
		}
		totalCombinations++
		if len(teamsMap) > 0 {
			logger.Info(ctx, fmt.Sprintf("        ✅ Активная группа '%s': +%d команд", activeGroupName, len(teamsMap)))
		}
	}

	for _, group := range groups {
		groupTeams := p.parseTeamsFromAjax(ctx, domain, group.AjaxURL)
		for _, team := range groupTeams {
			groupName := group.Name
			teamsWithContext = append(teamsWithContext, types.TeamWithContext{
				Team:      team,
				BirthYear: fallbackYear,
				GroupName: &groupName,
			})
		}
		totalCombinations++
	}

	logger.Info(ctx, fmt.Sprintf("     📊 Обработано комбинаций: %d", totalCombinations))
	logger.Info(ctx, fmt.Sprintf("     💾 Итого команд с контекстом: %d", len(teamsWithContext)))
	return teamsWithContext
}

func (p *Parser) parseWithFilters(ctx context.Context, domain string, doc *goquery.Document, yearLinks []types.YearLink, initialGroups []types.GroupLink) []types.TeamWithContext {
	var teamsWithContext []types.TeamWithContext
	totalCombinations := 0
	skippedYears := 0

	if len(yearLinks) == 0 {
		logger.Info(ctx, fmt.Sprintf("     📅 Годов нет, найдено %d групп", len(initialGroups)))

		// Парсим активную (дефолтную) группу из текущей страницы
		activeGroupName := helpers.ExtractActiveGroupName(doc, "competitions-teams")
		if activeGroupName != "" {
			teamsMap := make(map[string]types.TeamDTO)
			helpers.ParseTeamsFromDocWithDomain(doc, teamsMap, domain)

			for _, team := range teamsMap {
				gn := activeGroupName
				teamsWithContext = append(teamsWithContext, types.TeamWithContext{
					Team:      team,
					BirthYear: nil,
					GroupName: &gn,
				})
			}
			totalCombinations++

			if len(teamsMap) > 0 {
				logger.Info(ctx, fmt.Sprintf("        ✅ Активная группа '%s': +%d команд", activeGroupName, len(teamsMap)))
			}
		}

		// Остальные группы через AJAX
		for _, group := range initialGroups {
			groupTeams := p.parseTeamsFromAjax(ctx, domain, group.AjaxURL)
			for _, team := range groupTeams {
				groupName := group.Name
				teamsWithContext = append(teamsWithContext, types.TeamWithContext{
					Team:      team,
					BirthYear: nil,
					GroupName: &groupName,
				})
			}
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

			yearTeams := p.processYear(ctx, domain, yearLink, &totalCombinations)
			teamsWithContext = append(teamsWithContext, yearTeams...)
		}
	}

	logger.Info(ctx, fmt.Sprintf("     📊 Обработано комбинаций: %d, пропущено годов: %d", totalCombinations, skippedYears))
	logger.Info(ctx, fmt.Sprintf("     💾 Итого команд с контекстом: %d", len(teamsWithContext)))
	return teamsWithContext
}

func (p *Parser) processYear(ctx context.Context, domain string, yearLink types.YearLink, totalCombinations *int) []types.TeamWithContext {
	var teamsWithContext []types.TeamWithContext

	fullYearURL := domain + yearLink.AjaxURL
	yearResp, err := p.http.MakeRequest(fullYearURL)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("        ⚠️  Ошибка запроса года: %v", err))
		return teamsWithContext
	}

	if yearResp.StatusCode != 200 {
		_ = yearResp.Body.Close()
		return teamsWithContext
	}

	yearDoc, err := goquery.NewDocumentFromReader(yearResp.Body)
	_ = yearResp.Body.Close()
	if err != nil {
		return teamsWithContext
	}

	groupLinks := helpers.ExtractGroupLinks(yearDoc)
	logger.Info(ctx, fmt.Sprintf("        📋 ExtractGroupLinks returned %d groups", len(groupLinks)))
	for gi, gl := range groupLinks {
		logger.Info(ctx, fmt.Sprintf("           [%d] Group: '%s', AJAX: %s", gi+1, gl.Name, gl.AjaxURL[:80]))
	}

	if len(groupLinks) == 0 {
		// Год без групп - все команды принадлежат этому году
		teamsMap := make(map[string]types.TeamDTO)
		helpers.ParseTeamsFromDocWithDomain(yearDoc, teamsMap, domain)
		*totalCombinations++

		birthYear := yearLink.Year
		for _, team := range teamsMap {
			teamsWithContext = append(teamsWithContext, types.TeamWithContext{
				Team:      team,
				BirthYear: &birthYear,
				GroupName: nil,
			})
		}

		if len(teamsMap) > 0 {
			logger.Info(ctx, fmt.Sprintf("        ✅ Год без групп: +%d команд", len(teamsMap)))
		}
	} else {
		// Год с группами: сначала парсим дефолтную (активную) группу из yearDoc
		activeGroupName := helpers.ExtractActiveGroupName(yearDoc, "competitions-teams")
		if activeGroupName != "" {
			teamsMap := make(map[string]types.TeamDTO)
			helpers.ParseTeamsFromDocWithDomain(yearDoc, teamsMap, domain)
			birthYear := yearLink.Year

			for _, team := range teamsMap {
				gn := activeGroupName
				teamsWithContext = append(teamsWithContext, types.TeamWithContext{
					Team:      team,
					BirthYear: &birthYear,
					GroupName: &gn,
				})
			}
			*totalCombinations++

			if len(teamsMap) > 0 {
				logger.Info(ctx, fmt.Sprintf("        ✅ Активная группа '%s': +%d команд", activeGroupName, len(teamsMap)))
			}
		}

		// Затем обрабатываем остальные группы через AJAX
		for _, group := range groupLinks {
			logger.Info(ctx, fmt.Sprintf("        🔄 Fetching group '%s' via AJAX...", group.Name))
			groupTeams := p.parseTeamsFromAjax(ctx, domain, group.AjaxURL)
			logger.Info(ctx, fmt.Sprintf("        📊 Group '%s': got %d teams", group.Name, len(groupTeams)))
			birthYear := yearLink.Year
			groupName := group.Name

			for _, team := range groupTeams {
				teamsWithContext = append(teamsWithContext, types.TeamWithContext{
					Team:      team,
					BirthYear: &birthYear,
					GroupName: &groupName,
				})
			}
			*totalCombinations++
		}
	}

	return teamsWithContext
}

func (p *Parser) parseTeamsFromAjax(ctx context.Context, domain, ajaxURL string) []types.TeamDTO {
	fullURL := domain + ajaxURL

	ajaxResp, err := p.http.MakeRequest(fullURL)
	if err != nil {
		return nil
	}

	if ajaxResp.StatusCode != 200 {
		_ = ajaxResp.Body.Close()
		return nil
	}

	ajaxDoc, err := goquery.NewDocumentFromReader(ajaxResp.Body)
	_ = ajaxResp.Body.Close()
	if err != nil {
		return nil
	}

	teamsMap := make(map[string]types.TeamDTO)
	helpers.ParseTeamsFromDocWithDomain(ajaxDoc, teamsMap, domain)

	teams := make([]types.TeamDTO, 0, len(teamsMap))
	for _, team := range teamsMap {
		teams = append(teams, team)
	}

	if len(teams) > 0 {
		logger.Info(ctx, fmt.Sprintf("              ✅ +%d команд", len(teams)))
	}

	return teams
}
