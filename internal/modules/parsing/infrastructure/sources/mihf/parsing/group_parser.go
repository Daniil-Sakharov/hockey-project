package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/PuerkitoBio/goquery"
)

var (
	groupURLRegex       = regexp.MustCompile(`/championat/(\d{4})/groups/(\d+)`)
	tournamentURLRegex  = regexp.MustCompile(`/championat/(\d{4})/groups/(\d+)/tournament/(\d+)$`)
	subTournamentRegex  = regexp.MustCompile(`/championat/(\d{4})/groups/(\d+)/tournament/(\d+)/sub/(\d+)`)
	birthYearRegex      = regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
	birthYearSimple     = regexp.MustCompile(`(\d{4})`)
	groupNameRegex      = regexp.MustCompile(`[Гг]руппа\s+([А-Яа-яA-Za-z])`)
)

// ParseGroups парсит группы турниров для сезона
func ParseGroups(html []byte, seasonYear string) ([]dto.GroupDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	groups := make(map[string]dto.GroupDTO)

	doc.Find("a[href*='/groups/']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := groupURLRegex.FindStringSubmatch(href)
		if len(matches) < 3 {
			return
		}

		groupID := matches[2]
		text := strings.TrimSpace(s.Text())

		// Извлекаем год рождения
		birthYear := 0
		if byMatches := birthYearRegex.FindStringSubmatch(text); len(byMatches) > 1 {
			birthYear, _ = strconv.Atoi(byMatches[1])
		}

		if _, exists := groups[groupID]; !exists {
			groups[groupID] = dto.GroupDTO{
				ID:        groupID,
				Name:      text,
				BirthYear: birthYear,
				URL:       href,
			}
		}
	})

	result := make([]dto.GroupDTO, 0, len(groups))
	for _, g := range groups {
		result = append(result, g)
	}
	return result, nil
}

// FilterByMinBirthYear фильтрует группы по минимальному году рождения
func FilterByMinBirthYear(groups []dto.GroupDTO, minYear int) []dto.GroupDTO {
	var result []dto.GroupDTO
	for _, g := range groups {
		if g.BirthYear >= minYear {
			result = append(result, g)
		}
	}
	return result
}

// ParseTournaments парсит турниры (по году рождения) из страницы группы
// URL формат: /championat/2023/groups/76/tournament/330
func ParseTournaments(html []byte, seasonYear, groupID string) ([]dto.TournamentDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	tournaments := make(map[string]dto.TournamentDTO)

	doc.Find("a[href*='/tournament/']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Пропускаем ссылки с /sub/ - это уже подтурниры
		if strings.Contains(href, "/sub/") {
			return
		}

		matches := tournamentURLRegex.FindStringSubmatch(href)
		if len(matches) < 4 {
			return
		}

		tournamentID := matches[3]
		text := strings.TrimSpace(s.Text())

		// Извлекаем год рождения
		birthYear := 0
		if byMatches := birthYearRegex.FindStringSubmatch(text); len(byMatches) > 1 {
			birthYear, _ = strconv.Atoi(byMatches[1])
		} else if byMatches := birthYearSimple.FindStringSubmatch(text); len(byMatches) > 1 {
			year, _ := strconv.Atoi(byMatches[1])
			if year >= 2000 && year <= 2025 {
				birthYear = year
			}
		}

		if _, exists := tournaments[tournamentID]; !exists {
			tournaments[tournamentID] = dto.TournamentDTO{
				ID:        tournamentID,
				GroupID:   groupID,
				Name:      text,
				BirthYear: birthYear,
				URL:       href,
			}
		}
	})

	result := make([]dto.TournamentDTO, 0, len(tournaments))
	for _, t := range tournaments {
		result = append(result, t)
	}
	return result, nil
}

// ParseSubTournaments парсит подтурниры (Группа А, Б, В) из страницы турнира
// URL формат: /championat/2023/groups/76/tournament/330/sub/934
func ParseSubTournaments(html []byte, seasonYear, groupID, tournamentID string) ([]dto.SubTournamentDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	subTournaments := make(map[string]dto.SubTournamentDTO)

	doc.Find("a[href*='/sub/']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := subTournamentRegex.FindStringSubmatch(href)
		if len(matches) < 5 {
			return
		}

		subID := matches[4]
		text := strings.TrimSpace(s.Text())

		if _, exists := subTournaments[subID]; !exists {
			subTournaments[subID] = dto.SubTournamentDTO{
				ID:           subID,
				TournamentID: tournamentID,
				GroupID:      groupID,
				Name:         text,
				URL:          href,
			}
		}
	})

	result := make([]dto.SubTournamentDTO, 0, len(subTournaments))
	for _, st := range subTournaments {
		result = append(result, st)
	}
	return result, nil
}

// BuildTournamentPath создаёт TournamentPathDTO из компонентов
func BuildTournamentPath(seasonYear string, tournament dto.TournamentDTO, sub dto.SubTournamentDTO) dto.TournamentPathDTO {
	return dto.TournamentPathDTO{
		SeasonYear:   seasonYear,
		GroupID:      tournament.GroupID,
		TournamentID: tournament.ID,
		SubID:        sub.ID,
		BirthYear:    tournament.BirthYear,
		GroupName:    sub.Name,
	}
}
