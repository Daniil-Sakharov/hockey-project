package match

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parsePeriodScores парсит счёт по периодам из таблицы #ScoreGridView
func (p *Parser) parsePeriodScores(doc *goquery.Document, details *MatchDetailsDTO) {
	doc.Find("#ScoreGridView tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return // заголовок
		}

		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}

		teamName := strings.TrimSpace(cells.Eq(0).Text())
		p1, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(1).Text()))
		p2, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(2).Text()))
		p3, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(3).Text()))

		// Проверяем есть ли колонка ОТ (6-я колонка = индекс 5)
		// Структура: Команда | 1 | 2 | 3 | ОТ? | Всего
		hasOT := cells.Length() >= 6
		otVal := 0
		if hasOT {
			otText := strings.TrimSpace(cells.Eq(4).Text())
			otVal, _ = strconv.Atoi(otText)
		}

		if i == 1 {
			details.HomeTeamName = teamName
			details.HomeScoreP1 = p1
			details.HomeScoreP2 = p2
			details.HomeScoreP3 = p3
			if hasOT {
				details.HomeScoreOT = otVal
			}
		} else if i == 2 {
			details.AwayTeamName = teamName
			details.AwayScoreP1 = p1
			details.AwayScoreP2 = p2
			details.AwayScoreP3 = p3
			if hasOT {
				details.AwayScoreOT = otVal
			}
		}
	})
}

// parseShots парсит броски в створ ворот из таблицы #ShotGridView
func (p *Parser) parseShots(doc *goquery.Document, details *MatchDetailsDTO) {
	doc.Find("#ShotGridView tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return // заголовок
		}

		cells := row.Find("td")
		if cells.Length() < 5 {
			return
		}

		p1, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(1).Text()))
		p2, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(2).Text()))
		p3, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(3).Text()))

		// Итого может быть в bold
		totalText := strings.TrimSpace(cells.Eq(4).Text())
		if b := cells.Eq(4).Find("b"); b.Length() > 0 {
			totalText = strings.TrimSpace(b.Text())
		}
		total, _ := strconv.Atoi(totalText)

		// Проверяем есть ли колонка ОТ
		ot := 0
		if cells.Length() >= 6 {
			otText := strings.TrimSpace(cells.Eq(4).Text())
			ot, _ = strconv.Atoi(otText)
			// Тогда итого в 5-й колонке
			totalText = strings.TrimSpace(cells.Eq(5).Text())
			if b := cells.Eq(5).Find("b"); b.Length() > 0 {
				totalText = strings.TrimSpace(b.Text())
			}
			total, _ = strconv.Atoi(totalText)
		}

		shots := ShotsDTO{
			P1:    p1,
			P2:    p2,
			P3:    p3,
			OT:    ot,
			Total: total,
		}

		if i == 1 {
			details.HomeShots = shots
		} else if i == 2 {
			details.AwayShots = shots
		}
	})
}
