package game

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseTournamentHeader извлекает год рождения и группу из заголовка матча
// Формат: "Первенство ПФО 18/17/16/15 лет, 2008г.р, Группа А1"
func (p *Parser) parseTournamentHeader(doc *goquery.Document, d *GameDetailsDTO) {
	// Ищем заголовок в разных местах
	headerSelectors := []string{
		"th[colspan]",
		".match-header h1",
		".match-header h2",
		".tournament-title",
		".breadcrumbs a",
	}

	var headerText string
	for _, selector := range headerSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && (strings.Contains(text, "г.р") || strings.Contains(text, "Группа")) {
				headerText = text
			}
		})
		if headerText != "" {
			break
		}
	}

	if headerText == "" {
		return
	}

	// Извлекаем год рождения: "2008г.р" или "2008 г.р."
	birthYearRegex := regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
	if matches := birthYearRegex.FindStringSubmatch(headerText); len(matches) >= 2 {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			d.BirthYear = year
		}
	}

	// Извлекаем группу: "Группа А1", "Группа Б", "гр. А"
	groupRegex := regexp.MustCompile(`(?:Группа|гр\.?)\s*([А-Яа-яA-Za-z0-9]+)`)
	if matches := groupRegex.FindStringSubmatch(headerText); len(matches) >= 2 {
		d.GroupName = matches[1]
	}
}
