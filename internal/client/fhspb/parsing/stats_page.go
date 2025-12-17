package parsing

import (
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

func ParseStatsPage(doc *goquery.Document) dto.StatsPageDTO {
	result := dto.StatsPageDTO{CurrentPage: 1, TotalPages: 1}

	// ViewState
	result.ViewState, _ = doc.Find(`input#__VIEWSTATE`).Attr("value")
	result.ViewStateGenerator, _ = doc.Find(`input#__VIEWSTATEGENERATOR`).Attr("value")
	result.EventValidation, _ = doc.Find(`input#__EVENTVALIDATION`).Attr("value")

	// Текущая страница
	if text := doc.Find("span.current-page").First().Text(); text != "" {
		if page, err := strconv.Atoi(text); err == nil {
			result.CurrentPage = page
		}
	}

	// Общее количество страниц - ищем максимальный номер
	result.TotalPages = result.CurrentPage
	doc.Find("table.pager a.page").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			// href="javascript:__doPostBack('...','Page$6')"
			if idx := strings.Index(href, "Page$"); idx != -1 {
				numStr := strings.TrimSuffix(href[idx+5:], "')")
				if page, err := strconv.Atoi(numStr); err == nil && page > result.TotalPages {
					result.TotalPages = page
				}
			}
		}
	})

	return result
}
