package junior

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
)

// MinBirthYear минимальный год рождения игроков для парсинга
// Игроки с годом рождения < MinBirthYear будут пропущены
const MinBirthYear = 2008

// ParseTeamsFromTournament парсит команды из турнира (с AJAX переключением групп/возрастов)
// Двухуровневая логика: сначала годы, потом группы для каждого года
func (c *Client) ParseTeamsFromTournament(ctx context.Context, domain, tournamentURL string) ([]TeamDTO, error) {
	// Формируем URL страницы /teams/
	teamsURL := domain + tournamentURL
	if !strings.HasSuffix(teamsURL, "/") {
		teamsURL += "/"
	}
	teamsURL += "teams/"

	logger.Info(ctx, fmt.Sprintf("  🏒 Загрузка страницы команд: %s", teamsURL))

	// Загружаем основную страницу
	resp, err := c.makeRequest(teamsURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки страницы команд: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	teamsMap := make(map[string]TeamDTO) // Дедупликация по URL

	// Шаг 1: Извлекаем все ГОДА с начальной страницы
	yearLinks := c.extractYearLinks(ctx, doc)

	// Шаг 2: Извлекаем группы с начальной страницы (для текущего года)
	initialGroups := c.extractGroupLinks(doc)

	if len(yearLinks) == 0 && len(initialGroups) == 0 {
		// Нет переключателей - парсим основную страницу
		logger.Info(ctx, "     ℹ️  Нет переключателей года/группы, парсим основную страницу")
		parseTeamsFromDoc(doc, teamsMap)
		logger.Info(ctx, fmt.Sprintf("     💾 Найдено команд: %d", len(teamsMap)))
	} else {
		// Двухуровневый парсинг: годы → группы
		totalCombinations := 0
		skippedYears := 0

		if len(yearLinks) == 0 {
			// Нет годов, но есть группы - парсим только группы
			logger.Info(ctx, fmt.Sprintf("     📅 Годов нет, найдено %d групп", len(initialGroups)))
			for _, groupURL := range initialGroups {
				c.parseTeamsFromAjax(ctx, domain, groupURL, teamsMap)
				totalCombinations++
			}
		} else {
			logger.Info(ctx, fmt.Sprintf("     📅 Найдено %d годов рождения", len(yearLinks)))

			// Для каждого года
			for yearIdx, yearLink := range yearLinks {
				// ФИЛЬТРАЦИЯ: пропускаем года рождения < MinBirthYear (2008)
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

				// Делаем AJAX запрос для года
				fullYearURL := domain + yearLink.AjaxURL
				yearResp, err := c.makeRequest(fullYearURL)
				if err != nil {
					logger.Warn(ctx, fmt.Sprintf("        ⚠️  Ошибка запроса года: %v", err))
					continue
				}

				if yearResp.StatusCode != 200 {
					yearResp.Body.Close()
					logger.Warn(ctx, fmt.Sprintf("        ⚠️  HTTP %d для года", yearResp.StatusCode))
					continue
				}

				yearDoc, err := goquery.NewDocumentFromReader(yearResp.Body)
				yearResp.Body.Close()

				if err != nil {
					logger.Warn(ctx, fmt.Sprintf("        ⚠️  Ошибка парсинга года: %v", err))
					continue
				}

				// Извлекаем группы для ЭТОГО года
				groupLinks := c.extractGroupLinks(yearDoc)

				if len(groupLinks) == 0 {
					// Нет групп для этого года - парсим команды напрямую
					beforeCount := len(teamsMap)
					parseTeamsFromDoc(yearDoc, teamsMap)
					newCount := len(teamsMap) - beforeCount
					totalCombinations++

					if newCount > 0 {
						logger.Info(ctx, fmt.Sprintf("        ✅ Год без групп: +%d команд (всего: %d)", newCount, len(teamsMap)))
					}
				} else {
					logger.Info(ctx, fmt.Sprintf("        📁 Найдено %d групп для этого года", len(groupLinks)))

					// Для каждой группы в этом году
					for groupIdx, groupURL := range groupLinks {
						logger.Info(ctx, fmt.Sprintf("           [%d/%d] Группа...", groupIdx+1, len(groupLinks)))
						c.parseTeamsFromAjax(ctx, domain, groupURL, teamsMap)
						totalCombinations++
					}
				}
			}
		}

		logger.Info(ctx, fmt.Sprintf("     📊 Обработано комбинаций: %d, пропущено годов: %d", totalCombinations, skippedYears))
		logger.Info(ctx, fmt.Sprintf("     💾 Итого уникальных команд: %d", len(teamsMap)))
	}

	// Преобразуем map в slice
	teams := make([]TeamDTO, 0, len(teamsMap))
	for _, team := range teamsMap {
		teams = append(teams, team)
	}

	return teams, nil
}

// parseTeamsFromAjax делает AJAX запрос и парсит команды
func (c *Client) parseTeamsFromAjax(ctx context.Context, domain, ajaxURL string, teamsMap map[string]TeamDTO) {
	fullURL := domain + ajaxURL

	ajaxResp, err := c.makeRequest(fullURL)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("              ⚠️  Ошибка запроса: %v", err))
		return
	}

	if ajaxResp.StatusCode != 200 {
		ajaxResp.Body.Close()
		logger.Warn(ctx, fmt.Sprintf("              ⚠️  HTTP %d", ajaxResp.StatusCode))
		return
	}

	ajaxDoc, err := goquery.NewDocumentFromReader(ajaxResp.Body)
	ajaxResp.Body.Close()

	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("              ⚠️  Ошибка парсинга: %v", err))
		return
	}

	beforeCount := len(teamsMap)
	parseTeamsFromDoc(ajaxDoc, teamsMap)
	newCount := len(teamsMap) - beforeCount

	if newCount > 0 {
		logger.Info(ctx, fmt.Sprintf("              ✅ +%d команд (всего: %d)", newCount, len(teamsMap)))
	}
}

// parseTeamsFromDoc извлекает команды из HTML-документа и добавляет в teamsMap
func parseTeamsFromDoc(doc *goquery.Document, teamsMap map[string]TeamDTO) {
	// Селектор для команд (несколько вариантов для надежности)
	doc.Find("a.team-link, li.team-item a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		// Проверяем что это ссылка на команду (должна содержать /tournaments/)
		if !strings.Contains(href, "/tournaments/") {
			return
		}

		// Дедупликация
		if _, exists := teamsMap[href]; exists {
			return
		}

		// Извлекаем название и город
		name := strings.TrimSpace(s.Find(".team-title").Text())
		city := strings.TrimSpace(s.Find(".team-city").Text())

		// Fallback если структура другая
		if name == "" {
			name = strings.TrimSpace(s.Text())
		}

		team := TeamDTO{
			URL:  href,
			Name: name,
			City: city,
		}

		teamsMap[href] = team
	})
}
