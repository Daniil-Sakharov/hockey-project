package junior

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
)

// extractYearLinks извлекает AJAX-ссылки на ГОДЫ (из dropdown) с информацией о годе
func (c *Client) extractYearLinks(ctx context.Context, doc *goquery.Document) []YearLink {
	linksMap := make(map[string]YearLink) // key = ajaxURL для дедупликации

	logger.Info(ctx, "     🔍 Извлечение годов...")

	// Селектор 1: Dropdown годов - <option data-ajax>
	// HTML: <option value="2015" data-ajax="...">2015</option>
	doc.Find(`select.select-seasons option[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		ajax, exists := s.Attr("data-ajax")
		if !exists || ajax == "" {
			return
		}
		if !strings.Contains(ajax, "competitions-teams") {
			return
		}

		// Извлекаем год из value или текста option
		year := extractYearFromOption(s)
		if year > 0 {
			linksMap[ajax] = YearLink{Year: year, AjaxURL: ajax}
		}
	})

	// Селектор 2: Универсальный data-ajax с competitions-teams (только годы, без GROUP_ID)
	doc.Find(`[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		ajax, exists := s.Attr("data-ajax")
		if !exists || ajax == "" {
			return
		}
		if !strings.Contains(ajax, "competitions-teams") {
			return
		}
		// Проверяем что это год (GROUP_ID = null)
		if !isYearLink(ajax) {
			return
		}
		// Пропускаем если уже добавлен
		if _, exists := linksMap[ajax]; exists {
			return
		}

		// Пытаемся извлечь год из value или текста
		year := extractYearFromOption(s)
		if year > 0 {
			linksMap[ajax] = YearLink{Year: year, AjaxURL: ajax}
		}
	})

	links := make([]YearLink, 0, len(linksMap))
	for _, link := range linksMap {
		links = append(links, link)
	}

	logger.Info(ctx, fmt.Sprintf("        Найдено годов: %d", len(links)))
	return links
}

// extractYearFromOption извлекает год из <option> элемента
func extractYearFromOption(s *goquery.Selection) int {
	// Пробуем value атрибут
	if value, exists := s.Attr("value"); exists && value != "" {
		if year := parseYear(value); year > 0 {
			return year
		}
	}

	// Пробуем текст элемента
	text := strings.TrimSpace(s.Text())
	if year := parseYear(text); year > 0 {
		return year
	}

	return 0
}

// parseYear парсит год из строки (2008, 2009, ...)
func parseYear(s string) int {
	// Ищем 4-значное число (год)
	re := regexp.MustCompile(`\b(20\d{2})\b`)
	matches := re.FindStringSubmatch(s)
	if len(matches) > 1 {
		year, err := strconv.Atoi(matches[1])
		if err == nil && year >= 2000 && year <= 2025 {
			return year
		}
	}
	return 0
}

// extractGroupLinks извлекает AJAX-ссылки на ГРУППЫ (кнопки filter-btn)
func (c *Client) extractGroupLinks(doc *goquery.Document) []string {
	linksMap := make(map[string]bool)

	// Селектор: Кнопки групп - <div class="filter-btn" data-ajax-link>
	doc.Find(`div.filter-btn[data-ajax-link]`).Each(func(i int, s *goquery.Selection) {
		if ajax, exists := s.Attr("data-ajax-link"); exists && ajax != "" {
			if strings.Contains(ajax, "competitions-teams") {
				linksMap[ajax] = true
			}
		}
	})

	links := make([]string, 0, len(linksMap))
	for link := range linksMap {
		links = append(links, link)
	}

	return links
}

// isYearLink проверяет что ссылка относится к году (GROUP_ID = null)
func isYearLink(ajaxURL string) bool {
	// В base64 закодированных params, GROUP_ID:N выглядит как "R1JPVVBfSUQiO04"
	// или в декодированном виде "GROUP_ID";N

	// Если содержит GROUP_ID с конкретным значением (не N/null) - это группа
	re := regexp.MustCompile(`params=([^&]+)`)
	matches := re.FindStringSubmatch(ajaxURL)
	if len(matches) < 2 {
		return true // Не можем определить - считаем годом
	}

	paramsEncoded := matches[1]
	paramsEncoded = strings.ReplaceAll(paramsEncoded, "%3D", "=")

	// Проверяем что GROUP_ID = NULL (в base64 это "R1JPVVBfSUQiO04" или содержит ";N;")
	// Если GROUP_ID не null - это ссылка на группу
	if strings.Contains(paramsEncoded, "R1JPVVBfSUQiO3") {
		// GROUP_ID:s: означает что есть конкретное значение (строка)
		return false
	}

	return true
}

// extractAllAjaxLinks извлекает все AJAX-ссылки для загрузки команд (группы + возраста)
// DEPRECATED: используйте extractYearLinks + extractGroupLinks для двухуровневого парсинга
func (c *Client) extractAllAjaxLinks(ctx context.Context, doc *goquery.Document) []string {
	linksMap := make(map[string]bool) // дедупликация

	logger.Info(ctx, "     🔍 Извлечение AJAX-ссылок...")

	// Счетчики для каждого селектора
	selector1Count := 0
	selector2Count := 0
	selector3Count := 0

	// Селектор 1: Dropdown годов/сезонов - <option data-ajax>
	doc.Find(`select.select-seasons option[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		if ajax, exists := s.Attr("data-ajax"); exists && ajax != "" {
			// Фильтр: только competitions-teams
			if strings.Contains(ajax, "competitions-teams") {
				if !linksMap[ajax] {
					linksMap[ajax] = true
					selector1Count++
				}
			}
		}
	})

	// Селектор 2: Кнопки групп - <div class="filter-btn" data-ajax-link>
	doc.Find(`div.filter-btn[data-ajax-link]`).Each(func(i int, s *goquery.Selection) {
		if ajax, exists := s.Attr("data-ajax-link"); exists && ajax != "" {
			if strings.Contains(ajax, "competitions-teams") {
				if !linksMap[ajax] {
					linksMap[ajax] = true
					selector2Count++
				}
			}
		}
	})

	// Селектор 3: Универсальный - любой data-ajax с competitions-teams
	// (на случай изменения структуры HTML)
	doc.Find(`[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		if ajax, exists := s.Attr("data-ajax"); exists && ajax != "" {
			if strings.Contains(ajax, "competitions-teams") {
				if !linksMap[ajax] {
					linksMap[ajax] = true
					selector3Count++
				}
			}
		}
	})

	// Логируем результаты каждого селектора
	if selector1Count > 0 {
		logger.Info(ctx, fmt.Sprintf("        Селектор 1 (option[data-ajax]): найдено %d ссылок", selector1Count))
	}
	if selector2Count > 0 {
		logger.Info(ctx, fmt.Sprintf("        Селектор 2 (div.filter-btn): найдено %d ссылок", selector2Count))
	}
	if selector3Count > 0 {
		logger.Info(ctx, fmt.Sprintf("        Селектор 3 (universal [data-ajax]): найдено %d ссылок", selector3Count))
	}

	// Преобразуем map в slice
	links := make([]string, 0, len(linksMap))
	for link := range linksMap {
		links = append(links, link)
	}

	logger.Info(ctx, fmt.Sprintf("        Итого уникальных: %d ссылок", len(links)))

	return links
}

// classifyAjaxLink определяет тип AJAX-ссылки (год или группа) для логирования
func classifyAjaxLink(ajaxURL string) string {
	// Извлекаем params из URL
	re := regexp.MustCompile(`params=([^&]+)`)
	matches := re.FindStringSubmatch(ajaxURL)
	if len(matches) < 2 {
		return "unknown"
	}

	// Декодируем base64 (может быть с padding или без)
	paramsEncoded := matches[1]
	paramsEncoded = strings.ReplaceAll(paramsEncoded, "%3D", "=") // URL decode =

	// Пробуем декодировать
	// PHP сериализует так: a:4:{s:13:"TOURNAMENT_ID";s:8:"16743807";s:8:"GROUP_ID";N;...}
	// Если есть GROUP_ID не NULL - это группа
	if strings.Contains(paramsEncoded, "R1JPVVBfSUQ") || strings.Contains(paramsEncoded, "GROUP_ID") {
		// Проверяем что GROUP_ID не NULL
		if !strings.Contains(paramsEncoded, "R1JPVVBfSUQiO04") && !strings.Contains(paramsEncoded, `"GROUP_ID";N`) {
			return "group"
		}
	}

	return "year"
}
