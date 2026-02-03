package match

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parsePenalties парсит штрафы из секции с заголовком "Удаления"
func (p *Parser) parsePenalties(doc *goquery.Document) []PenaltyDTO {
	var penalties []PenaltyDTO
	tableIndex := 0

	// Ищем таблицы с заголовками "Время", "Игрок", "Минут", "Нарушение"
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		// Проверяем заголовки
		headers := table.Find("tr th")
		if headers.Length() < 4 {
			return
		}

		headerText := ""
		headers.Each(func(_ int, th *goquery.Selection) {
			headerText += " " + strings.TrimSpace(th.Text())
		})

		// Проверяем, что это таблица штрафов
		if !strings.Contains(headerText, "Время") || !strings.Contains(headerText, "Игрок") ||
			!strings.Contains(headerText, "Минут") || !strings.Contains(headerText, "Нарушение") {
			return
		}

		// Определяем: первая таблица = домашняя команда
		isHomeTeam := (tableIndex == 0)
		tableIndex++

		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			// Пропускаем заголовок
			if row.Find("th").Length() > 0 {
				return
			}

			penalty := p.parsePenaltyRow(row, isHomeTeam)
			if penalty != nil {
				penalties = append(penalties, *penalty)
			}
		})
	})

	return penalties
}

// parsePenaltyRow парсит одну строку со штрафом
func (p *Parser) parsePenaltyRow(row *goquery.Selection, isHome bool) *PenaltyDTO {
	cells := row.Find("td")
	if cells.Length() < 4 {
		return nil
	}

	penalty := &PenaltyDTO{IsHome: isHome}

	// Колонка 0: Время
	timeCell := cells.Eq(0)
	timeSpan := timeCell.Find("span.label, span.secondary")
	if timeSpan.Length() > 0 {
		timeText := strings.TrimSpace(timeSpan.Text())
		if m := timeRegex.FindStringSubmatch(timeText); len(m) >= 3 {
			penalty.TimeMinutes, _ = strconv.Atoi(m[1])
			penalty.TimeSeconds, _ = strconv.Atoi(m[2])
		}
	} else {
		// Пробуем из текста ячейки
		timeText := strings.TrimSpace(timeCell.Text())
		if m := timeRegex.FindStringSubmatch(timeText); len(m) >= 3 {
			penalty.TimeMinutes, _ = strconv.Atoi(m[1])
			penalty.TimeSeconds, _ = strconv.Atoi(m[2])
		}
	}

	// Колонка 1: Игрок
	playerCell := cells.Eq(1)
	playerLink := playerCell.Find("a[href*='Player']")
	if playerLink.Length() > 0 {
		penalty.PlayerURL, _ = playerLink.Attr("href")
		fullName := strings.TrimSpace(playerLink.Text())
		if m := numberRegex.FindStringSubmatch(fullName); len(m) >= 2 {
			penalty.PlayerNumber, _ = strconv.Atoi(m[1])
			penalty.PlayerName = strings.TrimSpace(strings.TrimPrefix(fullName, m[0]))
		} else {
			penalty.PlayerName = fullName
		}
	} else {
		// Текст напрямую
		text := strings.TrimSpace(playerCell.Text())
		if m := numberRegex.FindStringSubmatch(text); len(m) >= 2 {
			penalty.PlayerNumber, _ = strconv.Atoi(m[1])
			penalty.PlayerName = strings.TrimSpace(strings.TrimPrefix(text, m[0]))
		} else {
			penalty.PlayerName = text
		}
	}

	// Колонка 2: Минуты штрафа
	minCell := cells.Eq(2)
	minText := strings.TrimSpace(minCell.Text())
	minRegex := regexp.MustCompile(`(\d+)`)
	if m := minRegex.FindStringSubmatch(minText); len(m) >= 2 {
		penalty.Minutes, _ = strconv.Atoi(m[1])
	}

	// Колонка 3: Нарушение
	reasonCell := cells.Eq(3)
	violationSpan := reasonCell.Find("span[title]")
	if violationSpan.Length() > 0 {
		penalty.ReasonCode = strings.TrimSpace(violationSpan.Text())
		penalty.Reason, _ = violationSpan.Attr("title")
	} else {
		penalty.Reason = strings.TrimSpace(reasonCell.Text())
		penalty.ReasonCode = penalty.Reason
	}

	if penalty.PlayerName == "" {
		return nil
	}

	return penalty
}
