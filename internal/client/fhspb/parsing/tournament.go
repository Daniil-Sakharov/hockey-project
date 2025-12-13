package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
)

var (
	tournamentIDRegex = regexp.MustCompile(`TournamentID=(\d+)`)
	birthYearRegex    = regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
)

// ParseTournaments парсит список турниров из HTML
func ParseTournaments(html []byte) ([]dto.TournamentDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var tournaments []dto.TournamentDTO

	doc.Find("a[href*='TournamentID=']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := tournamentIDRegex.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			return
		}

		name := strings.TrimSpace(s.Text())
		if name == "" {
			return
		}

		tournaments = append(tournaments, dto.TournamentDTO{
			ID:        id,
			Name:      name,
			BirthYear: extractBirthYear(name),
		})
	})

	return deduplicateTournaments(tournaments), nil
}

func extractBirthYear(name string) int {
	matches := birthYearRegex.FindStringSubmatch(name)
	if len(matches) < 2 {
		return 0
	}
	year, _ := strconv.Atoi(matches[1])
	return year
}

func deduplicateTournaments(tournaments []dto.TournamentDTO) []dto.TournamentDTO {
	seen := make(map[int]bool)
	var result []dto.TournamentDTO
	for _, t := range tournaments {
		if !seen[t.ID] {
			seen[t.ID] = true
			result = append(result, t)
		}
	}
	return result
}

// FilterByBirthYear фильтрует турниры по году рождения
func FilterByBirthYear(tournaments []dto.TournamentDTO, maxYear int) []dto.TournamentDTO {
	var result []dto.TournamentDTO
	for _, t := range tournaments {
		if t.BirthYear > 0 && t.BirthYear <= maxYear {
			result = append(result, t)
		}
	}
	return result
}
