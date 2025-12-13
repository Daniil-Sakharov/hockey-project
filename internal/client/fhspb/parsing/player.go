package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
)

var numberRegex = regexp.MustCompile(`№(\d+)`)

// ParsePlayer парсит профиль игрока из HTML
func ParsePlayer(html []byte, playerID string) (*dto.PlayerDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	player := &dto.PlayerDTO{ExternalID: playerID}

	// ФИО из заголовка
	doc.Find("h3 a, h2 a").Each(func(_ int, s *goquery.Selection) {
		if player.FullName == "" {
			name := strings.TrimSpace(s.Text())
			if name != "" && len(name) > 5 {
				player.FullName = name
			}
		}
	})

	if player.FullName == "" {
		doc.Find("h3, h2").Each(func(_ int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if strings.Count(text, " ") >= 1 && len(text) > 5 && len(text) < 100 {
				player.FullName = text
			}
		})
	}

	parsePlayerFieldsFromDoc(player, doc)
	parsePlayerFieldsFromText(player, string(html))

	return player, nil
}

// parsePlayerFieldsFromDoc парсит поля из таблицы через goquery
func parsePlayerFieldsFromDoc(player *dto.PlayerDTO, doc *goquery.Document) {
	doc.Find("tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}

		label := strings.TrimSpace(cells.First().Text())
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
}

// parsePlayerFieldsFromText парсит позицию, номер, роль из текста
func parsePlayerFieldsFromText(player *dto.PlayerDTO, html string) {
	// Позиция
	if player.Position == "" {
		switch {
		case strings.Contains(html, "Нападающий"):
			player.Position = "Нападающий"
		case strings.Contains(html, "Защитник"):
			player.Position = "Защитник"
		case strings.Contains(html, "Вратарь"):
			player.Position = "Вратарь"
		}
	}

	// Номер
	if player.Number == 0 {
		if m := numberRegex.FindStringSubmatch(html); len(m) >= 2 {
			player.Number, _ = strconv.Atoi(m[1])
		}
	}

	// Роль (К/А)
	if player.Role == "" {
		roleRegex := regexp.MustCompile(`№\d+\s*([КА])`)
		if m := roleRegex.FindStringSubmatch(html); len(m) >= 2 {
			player.Role = m[1]
		}
	}
}

func parseIntFromStr(s string) int {
	numRegex := regexp.MustCompile(`(\d+)`)
	if m := numRegex.FindStringSubmatch(s); len(m) >= 2 {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}
