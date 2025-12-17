package parsing

import (
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

func ParseGoalieStats(doc *goquery.Document) []dto.GoalieStatsDTO {
	var results []dto.GoalieStatsDTO

	doc.Find("#StatsGridView tr").Each(func(i int, row *goquery.Selection) {
		// Пропускаем заголовки и пагинацию
		if row.Find("th").Length() > 0 || row.Find("table.pager").Length() > 0 {
			return
		}

		cells := row.Find("td")
		if cells.Length() < 14 {
			return
		}

		stats := dto.GoalieStatsDTO{}

		// Ссылка на игрока
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

		// Дата рождения
		stats.BirthDate = strings.TrimSpace(row.Find("span.description").First().Text())

		// Команда
		stats.TeamName = strings.TrimSpace(row.Find("a#TeamHyperLink").Text())

		// Статистика: И, Мин., Г, Бр., %, Ср., В, На 0, П, Шт
		stats.Games, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(4).Text()))
		stats.Minutes, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(5).Text()))
		stats.GoalsAgainst, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(6).Text()))
		stats.ShotsAgainst, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(7).Text()))
		stats.SavePercentage = parseFloat(cells.Eq(8).Text())
		stats.GoalsAgainstAvg = parseFloat(cells.Eq(9).Text())
		stats.Wins, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(10).Text()))
		stats.Shutouts, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(11).Text()))
		stats.Assists, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(12).Text()))
		stats.PenaltyMinutes, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(13).Text()))

		if stats.PlayerID != "" {
			results = append(results, stats)
		}
	})

	return results
}
