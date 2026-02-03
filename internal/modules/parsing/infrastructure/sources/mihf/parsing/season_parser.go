package parsing

import (
	"regexp"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/PuerkitoBio/goquery"
)

var seasonRegex = regexp.MustCompile(`/championat/(\d{4})-(\d{4})`)

// ParseSeasons парсит список сезонов с главной страницы
func ParseSeasons(html []byte) ([]dto.SeasonDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	seasons := make(map[string]dto.SeasonDTO)

	// Ищем ссылки на сезоны в меню "Статистика"
	doc.Find("a[href*='/championat/']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := seasonRegex.FindStringSubmatch(href)
		if len(matches) < 3 {
			return
		}

		year := matches[1]
		fullName := matches[1] + "-" + matches[2]

		if _, exists := seasons[year]; !exists {
			seasons[year] = dto.SeasonDTO{
				Year:     year,
				FullName: fullName,
				URL:      "/championat/" + fullName,
			}
		}
	})

	result := make([]dto.SeasonDTO, 0, len(seasons))
	for _, s := range seasons {
		result = append(result, s)
	}
	return result, nil
}
