package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/PuerkitoBio/goquery"
)

var teamIDRegex = regexp.MustCompile(`/team/(\d+)`)

// ParseScoreboard парсит турнирную таблицу (scoreboard)
func ParseScoreboard(html []byte) ([]dto.TeamDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var teams []dto.TeamDTO

	// Ищем таблицу с классом турнирной таблицы
	doc.Find("table tr").Each(func(i int, row *goquery.Selection) {
		// Пропускаем заголовки
		if row.Find("th").Length() > 0 {
			return
		}

		cells := row.Find("td")
		if cells.Length() < 8 {
			return
		}

		team := dto.TeamDTO{}

		// Ищем ссылку на команду
		teamLink := row.Find("a[href*='/team/']")
		if href, exists := teamLink.Attr("href"); exists {
			if matches := teamIDRegex.FindStringSubmatch(href); len(matches) > 1 {
				team.ID = matches[1]
				team.ExternalURL = href
			}
		}

		team.Name = strings.TrimSpace(teamLink.Text())

		// Парсим статистику
		// Обычный порядок: место, команда, И, В, Н, П, ШЗ, ШП, О
		cellIdx := 0
		cells.Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			val, _ := strconv.Atoi(text)

			switch cellIdx {
			case 2: // И - игры
				team.Games = val
			case 3: // В - победы
				team.Wins = val
			case 4: // Н - ничьи
				team.Draws = val
			case 5: // П - поражения
				team.Losses = val
			case 6: // ШЗ - забитые
				team.GoalsFor = val
			case 7: // ШП - пропущенные
				team.GoalsAgainst = val
			case 8: // О - очки (может быть на разных позициях)
				team.Points = val
			}
			cellIdx++
		})

		team.GoalsDiff = team.GoalsFor - team.GoalsAgainst

		if team.ID != "" && team.Name != "" {
			teams = append(teams, team)
		}
	})

	return teams, nil
}
