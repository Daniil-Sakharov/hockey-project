package helpers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// ExtractYearLinks извлекает AJAX-ссылки на ГОДЫ (из dropdown) для teams
func ExtractYearLinks(doc *goquery.Document) []types.YearLink {
	return extractYearLinksWithFilter(doc, "competitions-teams")
}

// ExtractYearLinksForStandings извлекает AJAX-ссылки на ГОДЫ для standings
func ExtractYearLinksForStandings(doc *goquery.Document) []types.YearLink {
	return extractYearLinksWithFilter(doc, "tournament-page")
}

// ExtractYearLinksForCalendar извлекает AJAX-ссылки на ГОДЫ для календаря
func ExtractYearLinksForCalendar(doc *goquery.Document) []types.YearLink {
	return extractYearLinksWithFilter(doc, "competitions-calendar")
}

// extractYearLinksWithFilter извлекает AJAX-ссылки на ГОДЫ с указанным фильтром
func extractYearLinksWithFilter(doc *goquery.Document, componentFilter string) []types.YearLink {
	linksMap := make(map[string]types.YearLink)

	// Селектор 1: Dropdown годов
	doc.Find(`select.select-seasons option[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		ajax, exists := s.Attr("data-ajax")
		if !exists || ajax == "" || !strings.Contains(ajax, componentFilter) {
			return
		}
		if year := ExtractYearFromOption(s); year > 0 {
			linksMap[ajax] = types.YearLink{Year: year, AjaxURL: ajax}
		}
	})

	// Селектор 2: Универсальный data-ajax
	doc.Find(`[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		ajax, exists := s.Attr("data-ajax")
		if !exists || ajax == "" || !strings.Contains(ajax, componentFilter) {
			return
		}
		if !IsYearLink(ajax) {
			return
		}
		if _, exists := linksMap[ajax]; exists {
			return
		}
		if year := ExtractYearFromOption(s); year > 0 {
			linksMap[ajax] = types.YearLink{Year: year, AjaxURL: ajax}
		}
	})

	links := make([]types.YearLink, 0, len(linksMap))
	for _, link := range linksMap {
		links = append(links, link)
	}
	return links
}

// ExtractGroupLinks извлекает AJAX-ссылки на ГРУППЫ для teams
func ExtractGroupLinks(doc *goquery.Document) []types.GroupLink {
	return extractGroupLinksWithFilter(doc, "competitions-teams")
}

// ExtractGroupLinksForStandings извлекает AJAX-ссылки на ГРУППЫ для standings
func ExtractGroupLinksForStandings(doc *goquery.Document) []types.GroupLink {
	return extractGroupLinksWithFilter(doc, "tournament-page")
}

// ExtractGroupLinksForCalendar извлекает AJAX-ссылки на ГРУППЫ для календаря
func ExtractGroupLinksForCalendar(doc *goquery.Document) []types.GroupLink {
	return extractGroupLinksWithFilter(doc, "competitions-calendar")
}

// extractGroupLinksWithFilter извлекает AJAX-ссылки на ГРУППЫ с указанным фильтром
func extractGroupLinksWithFilter(doc *goquery.Document, componentFilter string) []types.GroupLink {
	linksMap := make(map[string]string) // ajax URL -> group name

	doc.Find(`div.filter-btn[data-ajax-link]`).Each(func(i int, s *goquery.Selection) {
		// Пропускаем активную кнопку — она уже обработана через ExtractActiveGroupName
		if s.HasClass("active") {
			return
		}
		if ajax, exists := s.Attr("data-ajax-link"); exists && ajax != "" {
			if strings.Contains(ajax, componentFilter) {
				groupName := strings.TrimSpace(s.Text())
				if groupName == "" {
					groupName = "unknown"
				}
				linksMap[ajax] = groupName
			}
		}
	})

	links := make([]types.GroupLink, 0, len(linksMap))
	for ajaxURL, groupName := range linksMap {
		links = append(links, types.GroupLink{
			Name:    groupName,
			AjaxURL: ajaxURL,
		})
	}
	return links
}

// ExtractActiveGroupName находит имя активной группы (без data-ajax-link)
func ExtractActiveGroupName(doc *goquery.Document, componentFilter string) string {
	var activeName string

	// Ищем активную кнопку: filter-btn без data-ajax-link (или с классом active)
	doc.Find(`div.filter-btn`).Each(func(i int, s *goquery.Selection) {
		if activeName != "" {
			return
		}

		name := strings.TrimSpace(s.Text())
		// Игнорируем кнопки с годами (4 цифры) — это не группы
		if name == "" || ParseYear(name) > 0 {
			return
		}

		_, hasAjax := s.Attr("data-ajax-link")
		if hasAjax {
			// Если есть data-ajax-link, проверяем класс active
			if !s.HasClass("active") {
				return
			}
			// Активная кнопка с ajax-link — берём имя
			ajax, _ := s.Attr("data-ajax-link")
			if strings.Contains(ajax, componentFilter) {
				activeName = name
			}
			return
		}

		// Кнопка без data-ajax-link — скорее всего активная по умолчанию
		activeName = name
	})

	return activeName
}

// ExtractYearFromOption извлекает год из <option> элемента
func ExtractYearFromOption(s *goquery.Selection) int {
	if value, exists := s.Attr("value"); exists && value != "" {
		if year := ParseYear(value); year > 0 {
			return year
		}
	}
	text := strings.TrimSpace(s.Text())
	if year := ParseYear(text); year > 0 {
		return year
	}
	return 0
}

// ParseYear парсит год из строки
func ParseYear(s string) int {
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

// IsYearLink проверяет что ссылка относится к году
func IsYearLink(ajaxURL string) bool {
	re := regexp.MustCompile(`params=([^&]+)`)
	matches := re.FindStringSubmatch(ajaxURL)
	if len(matches) < 2 {
		return true
	}
	paramsEncoded := strings.ReplaceAll(matches[1], "%3D", "=")
	return !strings.Contains(paramsEncoded, "R1JPVVBfSUQiO3")
}
