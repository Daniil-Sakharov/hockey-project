package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

// ParsePlayer парсит профиль игрока из HTML
func ParsePlayer(html []byte, playerID string) (*dto.PlayerDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	player := &dto.PlayerDTO{ExternalID: playerID}

	// ФИО из заголовка h3 > a
	doc.Find("h3 a[href*='PlayerID']").First().Each(func(_ int, s *goquery.Selection) {
		player.FullName = strings.TrimSpace(s.Text())
	})

	// Позиция из h5.subheader (Вратарь/Защитник/Нападающий)
	doc.Find("h5.subheader").Each(func(_ int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "Вратарь" || text == "Защитник" || text == "Нападающий" {
			player.Position = text
		}
	})

	// Номер из span.label рядом с позицией: <span class="label">№<b>31</b></span>
	doc.Find("div.medium-8 span.label").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "№") {
			numStr := strings.TrimPrefix(text, "№")
			if n, err := strconv.Atoi(numStr); err == nil && n > 0 && n < 100 {
				player.Number = n
				// Роль (К/А) - следующий sibling span.warning.label
				s.NextFiltered("span.warning.label").Each(func(_ int, role *goquery.Selection) {
					text := strings.TrimSpace(role.Text())
					if text == "К" || text == "А" {
						player.Role = text
					}
				})
			}
		}
	})

	// Данные из таблицы
	doc.Find("table.panel tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}

		// Заменяем &nbsp; на пробел
		label := strings.ReplaceAll(cells.First().Text(), "\u00a0", " ")
		label = strings.TrimSpace(label)
		value := strings.TrimSpace(cells.Last().Text())

		switch label {
		case "Дата рождения":
			player.BirthDate = value
		case "Место рождения":
			player.BirthPlace = value
		case "Гражданство":
			player.Citizenship = value
		case "Рост":
			player.Height = parseIntFromStr(value)
		case "Вес":
			player.Weight = parseIntFromStr(value)
		case "Хват":
			player.Stick = value
		case "Воспитанник":
			player.School = value
		}
	})

	return player, nil
}

func parseIntFromStr(s string) int {
	numRegex := regexp.MustCompile(`(\d+)`)
	if m := numRegex.FindStringSubmatch(s); len(m) >= 2 {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}
