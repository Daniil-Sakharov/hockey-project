package helpers

import (
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// ParseTeamsFromDoc извлекает команды из HTML-документа
func ParseTeamsFromDoc(doc *goquery.Document, teamsMap map[string]types.TeamDTO) {
	ParseTeamsFromDocWithDomain(doc, teamsMap, "")
}

// ParseTeamsFromDocWithDomain извлекает команды из HTML-документа с привязкой домена для лого
func ParseTeamsFromDocWithDomain(doc *goquery.Document, teamsMap map[string]types.TeamDTO, domain string) {
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

		var logoURL string
		if img := s.Find("img"); img.Length() > 0 {
			if src, ok := img.Attr("src"); ok && strings.Contains(src, "team_logo") {
				if domain != "" && strings.HasPrefix(src, "/") {
					logoURL = domain + src
				} else {
					logoURL = src
				}
			}
		}

		teamsMap[href] = types.TeamDTO{
			URL:     href,
			Name:    name,
			City:    city,
			LogoURL: logoURL,
		}
	})
}
