package helpers

import (
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// ParseTeamsFromDoc извлекает команды из HTML-документа
func ParseTeamsFromDoc(doc *goquery.Document, teamsMap map[string]types.TeamDTO) {
	doc.Find("a.team-link, li.team-item a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || !strings.Contains(href, "/tournaments/") {
			return
		}
		if _, exists := teamsMap[href]; exists {
			return
		}

		name := strings.TrimSpace(s.Find(".team-title").Text())
		city := strings.TrimSpace(s.Find(".team-city").Text())
		if name == "" {
			name = strings.TrimSpace(s.Text())
		}

		teamsMap[href] = types.TeamDTO{
			URL:  href,
			Name: name,
			City: city,
		}
	})
}
