package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
	"github.com/PuerkitoBio/goquery"
)

var playerLinkRegex = regexp.MustCompile(`/player/(\d+)`)

// ParseTeamRoster парсит состав команды со страницы /team/{id}
// Возвращает список членов команды с их ID
func ParseTeamRoster(html []byte) ([]dto.TeamMemberDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var members []dto.TeamMemberDTO
	seen := make(map[string]bool)

	// Ищем все ссылки на игроков
	doc.Find("a[href*='/player/']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := playerLinkRegex.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}

		playerID := matches[1]
		if seen[playerID] {
			return
		}
		seen[playerID] = true

		member := dto.TeamMemberDTO{
			PlayerID: playerID,
			Name:     strings.TrimSpace(s.Text()),
		}

		// Пробуем найти номер в соседних элементах
		parent := s.Parent()
		if parent != nil {
			siblingText := parent.Text()
			if num := extractPlayerNumber(siblingText); num > 0 {
				member.Number = num
			}
		}

		members = append(members, member)
	})

	return members, nil
}

// ParseTeamRosterFromTable парсит таблицу состава команды
func ParseTeamRosterFromTable(html []byte) ([]dto.TeamMemberDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var members []dto.TeamMemberDTO
	seen := make(map[string]bool)

	// Ищем таблицу с составом
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		table.Find("tr").Each(func(_ int, row *goquery.Selection) {
			var member dto.TeamMemberDTO

			// Ищем ссылку на игрока в строке
			row.Find("a[href*='/player/']").Each(func(_ int, link *goquery.Selection) {
				href, _ := link.Attr("href")
				if matches := playerLinkRegex.FindStringSubmatch(href); len(matches) > 1 {
					member.PlayerID = matches[1]
					member.Name = strings.TrimSpace(link.Text())
				}
			})

			if member.PlayerID == "" || seen[member.PlayerID] {
				return
			}
			seen[member.PlayerID] = true

			// Парсим остальные ячейки
			row.Find("td").Each(func(i int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())

				// Первая ячейка часто номер
				if i == 0 {
					if num, err := strconv.Atoi(text); err == nil && num > 0 && num < 100 {
						member.Number = num
					}
				}

				// Ищем позицию
				textLower := strings.ToLower(text)
				if strings.Contains(textLower, "вратарь") || text == "В" {
					member.Position = "В"
				} else if strings.Contains(textLower, "защитник") || text == "З" {
					member.Position = "З"
				} else if strings.Contains(textLower, "нападающий") || text == "Н" {
					member.Position = "Н"
				}
			})

			members = append(members, member)
		})
	})

	return members, nil
}

func extractPlayerNumber(text string) int {
	// Ищем номер в формате #17 или № 17 или просто 17
	re := regexp.MustCompile(`[#№]?\s*(\d{1,2})`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		num, _ := strconv.Atoi(matches[1])
		if num > 0 && num < 100 {
			return num
		}
	}
	return 0
}
