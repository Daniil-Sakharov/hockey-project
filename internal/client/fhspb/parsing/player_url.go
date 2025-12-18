package parsing

import (
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

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

		matches := PlayerIDRegex.FindStringSubmatch(href)
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
			ProfileURL:   "https://www.fhspb.ru/" + strings.TrimPrefix(href, "/"),
		})
	})

	return urls, nil
}
