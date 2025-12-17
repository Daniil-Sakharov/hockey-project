package parsing

import (
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

func ParsePlayerStats(doc *goquery.Document) []dto.PlayerStatsDTO {
	var results []dto.PlayerStatsDTO

	doc.Find("#StatsGridView tr").Each(func(i int, row *goquery.Selection) {
		// Пропускаем заголовки и пагинацию
		if row.Find("th").Length() > 0 || row.Find("table.pager").Length() > 0 {
			return
		}

		cells := row.Find("td")
		if cells.Length() < 13 {
			return
		}

		stats := dto.PlayerStatsDTO{}

		// Ссылка на игрока - извлекаем PlayerID и TeamID
		playerLink, _ := row.Find("a#PlayerHyperLink").Attr("href")
		if matches := PlayerIDRegex.FindStringSubmatch(playerLink); len(matches) > 1 {
			stats.PlayerID = matches[1]
		}
		if matches := TeamIDRegex.FindStringSubmatch(playerLink); len(matches) > 1 {
			stats.TeamID = matches[1]
		}

		// Имя игрока
		stats.PlayerName = strings.TrimSpace(row.Find("a#PlayerHyperLink").Text())

		// Номер
		numberText := row.Find("span.label b").First().Text()
		stats.Number, _ = strconv.Atoi(numberText)

		// Роль (К/А)
		stats.Role = strings.TrimSpace(row.Find("span.warning.badge b").Text())

		// Дата рождения
		stats.BirthDate = strings.TrimSpace(row.Find("span.description").First().Text())

		// Команда
		stats.TeamName = strings.TrimSpace(row.Find("a#TeamHyperLink").Text())

		// Статистика из ячеек (начиная с 5-й)
		stats.Position = strings.TrimSpace(cells.Eq(4).Text())
		stats.Games, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(5).Text()))
		stats.Points, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(6).Text()))
		stats.PointsAvg = parseFloat(cells.Eq(7).Text())
		stats.Goals, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(8).Text()))
		stats.Assists, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(9).Text()))
		stats.PlusMinus, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(10).Text()))
		stats.PenaltyMinutes, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(11).Text()))
		stats.PenaltyAvg = parseFloat(cells.Eq(12).Text())

		if stats.PlayerID != "" {
			results = append(results, stats)
		}
	})

	return results
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	s = strings.TrimSuffix(s, "%")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
