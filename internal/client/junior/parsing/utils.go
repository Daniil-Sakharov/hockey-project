package parsing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ExtractTournamentID извлекает ID турнира из URL
func ExtractTournamentID(urlStr string) string {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(urlStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// ExtractGlobalSeason извлекает сезон со страницы (один для всех турниров)
func ExtractGlobalSeason(doc *goquery.Document) string {
	// 1. Пробуем извлечь из .select-current (видимый текст дропдауна)
	seasonText := strings.TrimSpace(doc.Find(".select-current").First().Text())
	if seasonText != "" {
		return seasonText // "2025/2026"
	}

	// 2. Fallback: извлекаем из <option selected>
	selectedOption := doc.Find("select.js-ajax-select option[selected]").First()
	if season, exists := selectedOption.Attr("value"); exists {
		return season // "2025/2026"
	}

	// 3. Если не найдено - возвращаем пустую строку
	return ""
}
