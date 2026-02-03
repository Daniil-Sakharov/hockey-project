package match

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseLineups парсит составы полевых игроков из секций h5:contains("Полевые игроки")
func (p *Parser) parseLineups(doc *goquery.Document, details *MatchDetailsDTO) {
	sectionIndex := 0

	doc.Find("h5").Each(func(_ int, header *goquery.Selection) {
		headerText := strings.TrimSpace(header.Text())
		if headerText != "Полевые игроки" {
			return
		}

		sectionIndex++
		// Первая секция = домашняя команда, вторая = гостевая
		isHomeTeam := (sectionIndex == 1)

		// Ищем div.scroll после заголовка
		scrollDiv := header.Next()
		if !scrollDiv.HasClass("scroll") {
			scrollDiv = header.NextAll().Filter(".scroll").First()
		}

		if scrollDiv.Length() == 0 {
			return
		}

		// Парсим таблицу
		scrollDiv.Find("table tr").Each(func(i int, row *goquery.Selection) {
			// Пропускаем заголовок
			if row.Find("th").Length() > 0 {
				return
			}

			player := p.parsePlayerLineupRow(row)
			if player == nil {
				return
			}

			if isHomeTeam {
				details.HomeLineup = append(details.HomeLineup, *player)
			} else {
				details.AwayLineup = append(details.AwayLineup, *player)
			}
		})
	})
}

// parsePlayerLineupRow парсит строку с полевым игроком
func (p *Parser) parsePlayerLineupRow(row *goquery.Selection) *PlayerLineupDTO {
	cells := row.Find("td")
	if cells.Length() < 7 {
		return nil
	}

	player := &PlayerLineupDTO{}

	// Колонка 0: Игрок + роль (А/К)
	playerCell := cells.Eq(0)
	playerLink := playerCell.Find("a[href*='Player']")
	if playerLink.Length() > 0 {
		player.PlayerURL, _ = playerLink.Attr("href")
		fullName := strings.TrimSpace(playerLink.Text())
		if m := numberRegex.FindStringSubmatch(fullName); len(m) >= 2 {
			player.Number, _ = strconv.Atoi(m[1])
			player.PlayerName = strings.TrimSpace(strings.TrimPrefix(fullName, m[0]))
		} else {
			player.PlayerName = fullName
		}
	}

	// Роль капитана/ассистента (span.label b или span.warning b)
	roleSpan := playerCell.Find("span.label b, span.warning b")
	if roleSpan.Length() > 0 {
		role := strings.TrimSpace(roleSpan.Text())
		if role == "К" {
			player.CaptainRole = "C" // Captain
		} else if role == "А" {
			player.CaptainRole = "A" // Assistant
		}
	}

	// Колонка 1: Амплуа
	posText := strings.TrimSpace(cells.Eq(1).Text())
	switch posText {
	case "Зщ":
		player.Position = "defenseman"
	case "Нп":
		player.Position = "forward"
	default:
		player.Position = posText
	}

	// Колонка 2: Играл
	playedText := strings.TrimSpace(cells.Eq(2).Text())
	player.Played = (playedText == "Да" || playedText == "да")

	// Колонка 3: Очки
	pointsText := strings.TrimSpace(cells.Eq(3).Text())
	if b := cells.Eq(3).Find("b"); b.Length() > 0 {
		pointsText = strings.TrimSpace(b.Text())
	}
	player.Points, _ = strconv.Atoi(pointsText)

	// Колонка 4: Голы
	player.Goals, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(4).Text()))

	// Колонка 5: Передачи
	player.Assists, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(5).Text()))

	// Колонка 6: Штрафные минуты
	player.PenaltyMinutes, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(6).Text()))

	// Колонка 7: +/- (если есть)
	if cells.Length() > 7 {
		player.PlusMinus, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(7).Text()))
	}

	if player.PlayerName == "" {
		return nil
	}

	return player
}

// parseGoalieStats парсит статистику вратарей из секций h5:contains("Вратари")
func (p *Parser) parseGoalieStats(doc *goquery.Document, details *MatchDetailsDTO) {
	sectionIndex := 0

	doc.Find("h5").Each(func(_ int, header *goquery.Selection) {
		headerText := strings.TrimSpace(header.Text())
		if headerText != "Вратари" {
			return
		}

		sectionIndex++
		// Первая секция = домашняя команда, вторая = гостевая
		isHomeTeam := (sectionIndex == 1)

		// Ищем div.scroll после заголовка
		scrollDiv := header.Next()
		if !scrollDiv.HasClass("scroll") {
			scrollDiv = header.NextAll().Filter(".scroll").First()
		}

		if scrollDiv.Length() == 0 {
			return
		}

		// Парсим таблицу
		scrollDiv.Find("table tr").Each(func(i int, row *goquery.Selection) {
			// Пропускаем заголовок
			if row.Find("th").Length() > 0 {
				return
			}

			goalie := p.parseGoalieRow(row)
			if goalie == nil {
				return
			}

			if isHomeTeam {
				details.HomeGoalies = append(details.HomeGoalies, *goalie)
			} else {
				details.AwayGoalies = append(details.AwayGoalies, *goalie)
			}
		})
	})
}

// parseGoalieRow парсит строку со статистикой вратаря
func (p *Parser) parseGoalieRow(row *goquery.Selection) *GoalieStatsDTO {
	cells := row.Find("td")
	if cells.Length() < 6 {
		return nil
	}

	goalie := &GoalieStatsDTO{}

	// Колонка 0: Вратарь
	goalieCell := cells.Eq(0)
	goalieLink := goalieCell.Find("a[href*='Player']")
	if goalieLink.Length() > 0 {
		goalie.PlayerURL, _ = goalieLink.Attr("href")
		fullName := strings.TrimSpace(goalieLink.Text())
		if m := numberRegex.FindStringSubmatch(fullName); len(m) >= 2 {
			goalie.Number, _ = strconv.Atoi(m[1])
			goalie.PlayerName = strings.TrimSpace(strings.TrimPrefix(fullName, m[0]))
		} else {
			goalie.PlayerName = fullName
		}
	}

	// Колонка 1: Играл
	playedText := strings.TrimSpace(cells.Eq(1).Text())
	goalie.Played = (playedText == "Да" || playedText == "да")

	// Колонка 2: Время на льду
	timeCell := cells.Eq(2)
	timeSpan := timeCell.Find("span.label, span.secondary")
	timeText := ""
	if timeSpan.Length() > 0 {
		timeText = strings.TrimSpace(timeSpan.Text())
	} else {
		timeText = strings.TrimSpace(timeCell.Text())
	}
	if m := timeRegex.FindStringSubmatch(timeText); len(m) >= 3 {
		min, _ := strconv.Atoi(m[1])
		sec, _ := strconv.Atoi(m[2])
		goalie.TimeOnIce = min*60 + sec
	}

	// Колонка 3: Пропущено голов
	goalsText := strings.TrimSpace(cells.Eq(3).Text())
	if b := cells.Eq(3).Find("b"); b.Length() > 0 {
		goalsText = strings.TrimSpace(b.Text())
	}
	goalie.GoalsAgainst, _ = strconv.Atoi(goalsText)

	// Колонка 4: Бросков
	goalie.ShotsAgainst, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(4).Text()))

	// Колонка 5: Процент отражённых
	pctText := strings.TrimSpace(cells.Eq(5).Text())
	if b := cells.Eq(5).Find("b"); b.Length() > 0 {
		pctText = strings.TrimSpace(b.Text())
	}
	pctText = strings.ReplaceAll(pctText, ",", ".")
	goalie.SavePercentage, _ = strconv.ParseFloat(pctText, 64)

	// Колонка 6: Штрафы (если есть)
	if cells.Length() > 6 {
		goalie.PenaltyMinutes, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(6).Text()))
	}

	if goalie.PlayerName == "" {
		return nil
	}

	return goalie
}
