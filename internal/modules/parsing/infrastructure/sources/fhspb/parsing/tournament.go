package parsing

import (
	"regexp"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

var (
	tournamentIDRegex = regexp.MustCompile(`TournamentID=(\d+)`)
	birthYearRegex    = regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
	dateRangeRegex    = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4})\s*-\s*(\d{2}\.\d{2}\.\d{4})`)
	groupNameRegex    = regexp.MustCompile(`\s*(Группа\s+[А-Яа-яA-Za-z0-9]+)\s*$`)
)

// ParseTournaments парсит список турниров из HTML
func ParseTournaments(html []byte) ([]dto.TournamentDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	tournaments := make(map[int]*dto.TournamentDTO)

	doc.Find("div.clearfix").Each(func(_ int, container *goquery.Selection) {
		t := parseTournamentCard(container)
		if t != nil && t.ID != 0 {
			if _, exists := tournaments[t.ID]; !exists {
				tournaments[t.ID] = t
			}
		}
	})

	result := make([]dto.TournamentDTO, 0, len(tournaments))
	for _, t := range tournaments {
		result = append(result, *t)
	}
	return result, nil
}

// FilterByBirthYear фильтрует турниры по минимальному году рождения
func FilterByBirthYear(tournaments []dto.TournamentDTO, minYear int) []dto.TournamentDTO {
	var result []dto.TournamentDTO
	for _, t := range tournaments {
		if t.BirthYear >= minYear {
			result = append(result, t)
		}
	}
	return result
}
