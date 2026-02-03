package match

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseGoals парсит голы из секции с заголовком "Голы"
func (p *Parser) parseGoals(doc *goquery.Document, details *MatchDetailsDTO) []GoalDTO {
	var goals []GoalDTO
	currentPeriod := 0

	// Находим h3 с текстом "Голы"
	doc.Find("h3").Each(func(_ int, header *goquery.Selection) {
		if !strings.Contains(header.Text(), "Голы") {
			return
		}

		// Ищем div.scroll среди следующих элементов
		var scrollDiv *goquery.Selection
		header.NextAll().EachWithBreak(func(_ int, el *goquery.Selection) bool {
			if el.HasClass("scroll") {
				scrollDiv = el
				return false // прекращаем поиск
			}
			return true
		})

		if scrollDiv == nil || scrollDiv.Length() == 0 {
			return
		}

		// Проходим по всем строкам таблицы внутри div.scroll
		scrollDiv.Find("table tr").Each(func(_ int, row *goquery.Selection) {
			// Проверяем: это заголовок периода?
			periodSpan := row.Find("span[id*='PeriodLabel']")
			if periodSpan.Length() > 0 {
				periodText := periodSpan.Text()
				currentPeriod = p.detectPeriod(periodText)
				return
			}

			// Также проверяем th с классом period-title
			periodTH := row.Find("th.period-title")
			if periodTH.Length() > 0 {
				periodText := periodTH.Text()
				currentPeriod = p.detectPeriod(periodText)
				return
			}

			// Ищем таблицу score-grid внутри этой строки
			row.Find("table.score-grid tr").Each(func(_ int, goalRow *goquery.Selection) {
				goal := p.parseGoalRow(goalRow, currentPeriod, details)
				if goal != nil {
					goals = append(goals, *goal)
				}
			})
		})
	})

	return goals
}

// detectPeriod определяет номер периода из текста
func (p *Parser) detectPeriod(text string) int {
	if strings.Contains(text, "1-") || strings.Contains(text, "1й") || strings.Contains(text, "1ый") {
		return 1
	}
	if strings.Contains(text, "2-") || strings.Contains(text, "2й") || strings.Contains(text, "2ой") {
		return 2
	}
	if strings.Contains(text, "3-") || strings.Contains(text, "3й") || strings.Contains(text, "3ий") {
		return 3
	}
	if strings.Contains(text, "ОТ") || strings.Contains(text, "Овертайм") {
		return 4
	}
	return 0
}

// parseGoalRow парсит одну строку с голом
func (p *Parser) parseGoalRow(row *goquery.Selection, period int, details *MatchDetailsDTO) *GoalDTO {
	cells := row.Find("td")
	if cells.Length() < 4 {
		return nil
	}

	goal := &GoalDTO{Period: period}

	cells.Each(func(i int, cell *goquery.Selection) {
		// Время (span.warning.label)
		if timeSpan := cell.Find("span.warning"); timeSpan.Length() > 0 {
			timeText := strings.TrimSpace(timeSpan.Text())
			if m := timeRegex.FindStringSubmatch(timeText); len(m) >= 3 {
				goal.TimeMinutes, _ = strconv.Atoi(m[1])
				goal.TimeSeconds, _ = strconv.Atoi(m[2])
			}
		}

		// Счёт (h4 > b или просто b в ячейке)
		if scoreB := cell.Find("h4 b, h4 > b"); scoreB.Length() > 0 {
			scoreText := strings.TrimSpace(scoreB.Text())
			if m := scoreRegex.FindStringSubmatch(scoreText); len(m) >= 3 {
				goal.ScoreHome, _ = strconv.Atoi(m[1])
				goal.ScoreAway, _ = strconv.Atoi(m[2])
			}
		}

		// Команда (span.has-tip с title)
		if teamSpan := cell.Find("span.has-tip, span[title]"); teamSpan.Length() > 0 {
			if title, exists := teamSpan.Attr("title"); exists && title != "" {
				goal.TeamName = title
			}
			goal.TeamAbbr = strings.TrimSpace(teamSpan.Text())
		}

		// Тип гола (span.secondary.label: +1, +2, -1, -2, ПВ, ШБ, БП)
		cell.Find("span.secondary.label, span.secondary").Each(func(_ int, span *goquery.Selection) {
			text := strings.TrimSpace(span.Text())
			switch text {
			case "+1":
				goal.GoalType = "PP1" // Power play 5на4
			case "+2":
				goal.GoalType = "PP2" // Power play 5на3
			case "-1":
				goal.GoalType = "SH1" // Short-handed 4на5
			case "-2":
				goal.GoalType = "SH2" // Short-handed 3на5
			case "ПВ":
				goal.GoalType = "EN" // Empty net
			case "ШБ":
				goal.GoalType = "PS" // Penalty shot
			case "БП":
				goal.GoalType = "GWG" // Game-winning goal
			}
		})

		// Автор гола (ссылка a[href*=Player] с bold текстом)
		playerLinks := cell.Find("a[href*='Player']")
		if playerLinks.Length() > 0 {
			firstLink := playerLinks.First()
			// Проверяем, есть ли bold внутри - это автор гола
			if firstLink.Find("b").Length() > 0 || cell.Find("b a[href*='Player']").Length() > 0 {
				p.extractScorerInfo(firstLink, goal)
			} else if i >= 3 { // Колонки с 3-й и далее обычно автор
				p.extractScorerInfo(firstLink, goal)
			}

			// Остальные ссылки - ассистенты
			assistIdx := 0
			playerLinks.Each(func(j int, link *goquery.Selection) {
				href, _ := link.Attr("href")
				// Пропускаем автора гола
				if href == goal.ScorerURL {
					return
				}
				name := strings.TrimSpace(link.Text())
				if name == "" {
					return
				}

				if assistIdx == 0 {
					goal.Assist1URL = href
					goal.Assist1Name = p.cleanPlayerName(name)
				} else if assistIdx == 1 {
					goal.Assist2URL = href
					goal.Assist2Name = p.cleanPlayerName(name)
				}
				assistIdx++
			})
		}
	})

	// Определяем домашняя ли команда забила
	if goal.TeamName != "" && details != nil {
		goal.IsHome = strings.Contains(details.HomeTeamName, goal.TeamName) ||
			strings.Contains(goal.TeamName, details.HomeTeamName)
	}

	if goal.ScorerName == "" {
		return nil
	}
	return goal
}

// extractScorerInfo извлекает информацию об авторе гола из ссылки
func (p *Parser) extractScorerInfo(link *goquery.Selection, goal *GoalDTO) {
	goal.ScorerURL, _ = link.Attr("href")
	fullText := strings.TrimSpace(link.Text())
	// Извлекаем номер и имя: "14 Гагарин Тимур"
	if m := numberRegex.FindStringSubmatch(fullText); len(m) >= 2 {
		goal.ScorerNumber, _ = strconv.Atoi(m[1])
		goal.ScorerName = strings.TrimSpace(strings.TrimPrefix(fullText, m[0]))
	} else {
		goal.ScorerName = fullText
	}
}

// cleanPlayerName очищает имя игрока от номера
func (p *Parser) cleanPlayerName(name string) string {
	if m := numberRegex.FindStringSubmatch(name); len(m) >= 2 {
		return strings.TrimSpace(strings.TrimPrefix(name, m[0]))
	}
	return name
}
