package parsing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
)

var teamIDRegex = regexp.MustCompile(`TeamID=([a-f0-9-]{36})`)

// ParseTeams парсит список команд турнира из HTML
func ParseTeams(html []byte, tournamentID int) ([]dto.TeamDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var teams []dto.TeamDTO

	doc.Find("a[href*='TeamID=']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := teamIDRegex.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}

		name := strings.TrimSpace(s.Text())
		if name == "" {
			return
		}

		teams = append(teams, dto.TeamDTO{
			ID:           matches[1],
			TournamentID: tournamentID,
			Name:         name,
		})
	})

	return deduplicateTeams(teams), nil
}

func deduplicateTeams(teams []dto.TeamDTO) []dto.TeamDTO {
	seen := make(map[string]bool)
	var result []dto.TeamDTO
	for _, t := range teams {
		if !seen[t.ID] {
			seen[t.ID] = true
			result = append(result, t)
		}
	}
	return result
}
