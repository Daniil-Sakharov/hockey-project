package helpers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// ExtractYearLinks извлекает AJAX-ссылки на ГОДЫ (из dropdown)
func ExtractYearLinks(doc *goquery.Document) []types.YearLink {
	linksMap := make(map[string]types.YearLink)

	// Селектор 1: Dropdown годов
	doc.Find(`select.select-seasons option[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		ajax, exists := s.Attr("data-ajax")
		if !exists || ajax == "" || !strings.Contains(ajax, "competitions-teams") {
			return
		}
		if year := ExtractYearFromOption(s); year > 0 {
			linksMap[ajax] = types.YearLink{Year: year, AjaxURL: ajax}
		}
	})

	// Селектор 2: Универсальный data-ajax
	doc.Find(`[data-ajax]`).Each(func(i int, s *goquery.Selection) {
		ajax, exists := s.Attr("data-ajax")
		if !exists || ajax == "" || !strings.Contains(ajax, "competitions-teams") {
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

// ExtractGroupLinks извлекает AJAX-ссылки на ГРУППЫ
func ExtractGroupLinks(doc *goquery.Document) []string {
	linksMap := make(map[string]bool)

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
