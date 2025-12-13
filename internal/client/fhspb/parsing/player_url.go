package parsing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
)

var playerIDRegex = regexp.MustCompile(`PlayerID=([a-f0-9-]{36})`)

// ParsePlayerURLs извлекает URL игроков из страницы команды
func ParsePlayerURLs(html []byte, tournamentID int, teamID string) ([]dto.PlayerURLDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var urls []dto.PlayerURLDTO
	seen := make(map[string]bool)

	doc.Find("a[href*='PlayerID=']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := playerIDRegex.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}

		playerID := matches[1]
		if seen[playerID] {
			return
		}
		seen[playerID] = true

		urls = append(urls, dto.PlayerURLDTO{
			URL:          href,
			PlayerID:     playerID,
			TeamID:       teamID,
			TournamentID: tournamentID,
		})
	})

	return urls, nil
}
