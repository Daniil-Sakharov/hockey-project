package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/PuerkitoBio/goquery"
)

var (
	playerLinkRegex  = regexp.MustCompile(`/players/info/(\d+)`)
	timeRegex        = regexp.MustCompile(`'(\d{1,2}):(\d{2})`)
	periodScoreRegex = regexp.MustCompile(`\((\d+:\d+(?:,\d+:\d+)*)\)`)
)

// ParseMatchProtocol парсит протокол матча
func ParseMatchProtocol(html []byte) (*dto.MatchProtocolDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	proto := &dto.MatchProtocolDTO{}

	parseTeamLogos(doc, proto)
	parseScoreByPeriods(doc, proto)
	parseGoalsTable(doc, proto)
	parsePenaltiesTable(doc, proto)
	proto.HomeLineup, proto.AwayLineup = parseLineupsTable(doc)

	return proto, nil
}

// parseTeamLogos извлекает URL логотипов команд
func parseTeamLogos(doc *goquery.Document, proto *dto.MatchProtocolDTO) {
	logos := doc.Find("img[src*='/contents/clubs/']")
	if logos.Length() >= 2 {
		proto.HomeLogoURL, _ = logos.Eq(0).Attr("src")
		proto.AwayLogoURL, _ = logos.Eq(1).Attr("src")
	}
}

// parseScoreByPeriods извлекает счет по периодам
func parseScoreByPeriods(doc *goquery.Document, proto *dto.MatchProtocolDTO) {
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		if m := periodScoreRegex.FindStringSubmatch(text); len(m) > 1 {
			periods := strings.Split(m[1], ",")
			for i, p := range periods {
				if i >= 4 {
					break
				}
				scores := strings.Split(strings.TrimSpace(p), ":")
				if len(scores) == 2 {
					proto.ScoreByPeriod[i][0], _ = strconv.Atoi(scores[0])
					proto.ScoreByPeriod[i][1], _ = strconv.Atoi(scores[1])
				}
			}
		}
	})
}

// parseGoalsTable парсит голы из таблицы с 4 колонками
// Структура: время_дом | игрок_дом | время_гость | игрок_гость
func parseGoalsTable(doc *goquery.Document, proto *dto.MatchProtocolDTO) {
	doc.Find("table.table-hover").Each(func(_ int, table *goquery.Selection) {
		headerRow := table.Find("tr.title").First()
		cells := headerRow.Find("td")

		// Таблица голов: 4 колонки (время, команда, время, команда)
		if cells.Length() != 4 {
			return
		}

		// Проверяем что это не таблица нарушений (у нее 6 колонок)
		// и не таблица составов (у нее 12 колонок)
		table.Find("tr.out_blue, tr.out_white").Each(func(_ int, row *goquery.Selection) {
			rowCells := row.Find("td")
			if rowCells.Length() < 4 {
				return
			}

			// Домашняя команда (колонки 0-1)
			homeGoal := parseGoalFromCells(rowCells.Eq(0), rowCells.Eq(1), true)
			if homeGoal != nil {
				proto.Goals = append(proto.Goals, *homeGoal)
			}

			// Гостевая команда (колонки 2-3)
			awayGoal := parseGoalFromCells(rowCells.Eq(2), rowCells.Eq(3), false)
			if awayGoal != nil {
				proto.Goals = append(proto.Goals, *awayGoal)
			}
		})
	})
}

// parseGoalFromCells парсит гол из ячеек времени и игрока
func parseGoalFromCells(timeCell, playerCell *goquery.Selection, isHome bool) *dto.GoalEventDTO {
	timeText := strings.TrimSpace(timeCell.Text())
	timeMatch := timeRegex.FindStringSubmatch(timeText)
	if len(timeMatch) != 3 {
		return nil
	}

	minutes, _ := strconv.Atoi(timeMatch[1])
	seconds, _ := strconv.Atoi(timeMatch[2])

	// Ищем автора гола
	scorerLink := playerCell.Find("a[href*='/players/info/']").First()
	if scorerLink.Length() == 0 {
		return nil
	}

	href, _ := scorerLink.Attr("href")
	m := playerLinkRegex.FindStringSubmatch(href)
	if len(m) < 2 || m[1] == "" {
		return nil
	}

	goal := &dto.GoalEventDTO{
		Period:      calculatePeriod(minutes),
		TimeMinutes: minutes,
		TimeSeconds: seconds,
		ScorerID:    m[1],
		ScorerName:  strings.TrimSpace(scorerLink.Text()),
		GoalType:    detectGoalType(strings.ToLower(playerCell.Text())),
		IsHome:      isHome,
	}

	// Парсим ассистентов (все ссылки кроме первой)
	assistCount := 0
	playerCell.Find("a[href*='/players/info/']").Each(func(i int, link *goquery.Selection) {
		if i == 0 {
			return // Пропускаем автора гола
		}
		href, _ := link.Attr("href")
		am := playerLinkRegex.FindStringSubmatch(href)
		if len(am) < 2 || am[1] == "" {
			return
		}
		assistCount++
		if assistCount == 1 {
			goal.Assist1ID = am[1]
			goal.Assist1Name = strings.TrimSpace(link.Text())
		} else if assistCount == 2 {
			goal.Assist2ID = am[1]
			goal.Assist2Name = strings.TrimSpace(link.Text())
		}
	})

	return goal
}

// parsePenaltiesTable парсит нарушения из таблицы с 6 колонками
// Структура: время_дом | минуты_дом | игрок_дом | время_гость | минуты_гость | игрок_гость
func parsePenaltiesTable(doc *goquery.Document, proto *dto.MatchProtocolDTO) {
	// Ищем заголовок "Нарушения"
	doc.Find("h1").Each(func(_ int, h1 *goquery.Selection) {
		if !strings.Contains(h1.Text(), "Нарушения") {
			return
		}

		// Ищем таблицу в следующем sibling div
		table := h1.Next().Find("table.table-hover").First()
		if table.Length() == 0 {
			return
		}

		table.Find("tr.out_blue, tr.out_white").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() < 6 {
				return
			}

			// Домашняя команда (колонки 0-2)
			homePenalty := parsePenaltyFromCells(cells.Eq(0), cells.Eq(1), cells.Eq(2), true)
			if homePenalty != nil {
				proto.Penalties = append(proto.Penalties, *homePenalty)
			}

			// Гостевая команда (колонки 3-5)
			awayPenalty := parsePenaltyFromCells(cells.Eq(3), cells.Eq(4), cells.Eq(5), false)
			if awayPenalty != nil {
				proto.Penalties = append(proto.Penalties, *awayPenalty)
			}
		})
	})
}

// parsePenaltyFromCells парсит нарушение из 3 ячеек
func parsePenaltyFromCells(timeCell, minCell, playerCell *goquery.Selection, isHome bool) *dto.PenaltyEventDTO {
	timeText := strings.TrimSpace(timeCell.Text())
	timeMatch := timeRegex.FindStringSubmatch(timeText)
	if len(timeMatch) != 3 {
		return nil
	}

	minutes, _ := strconv.Atoi(timeMatch[1])
	seconds, _ := strconv.Atoi(timeMatch[2])

	// Минуты штрафа
	penaltyCode, _ := strconv.Atoi(strings.TrimSpace(minCell.Text()))
	penaltyMinutes := convertPenaltyCode(penaltyCode)
	if penaltyMinutes == 0 {
		return nil
	}

	// Игрок
	playerLink := playerCell.Find("a[href*='/players/info/']").First()
	if playerLink.Length() == 0 {
		return nil
	}

	href, _ := playerLink.Attr("href")
	m := playerLinkRegex.FindStringSubmatch(href)
	if len(m) < 2 || m[1] == "" {
		return nil
	}

	return &dto.PenaltyEventDTO{
		Period:      calculatePeriod(minutes),
		TimeMinutes: minutes,
		TimeSeconds: seconds,
		PlayerID:    m[1],
		PlayerName:  strings.TrimSpace(playerLink.Text()),
		Minutes:     penaltyMinutes,
		IsHome:      isHome,
	}
}

// convertPenaltyCode преобразует код штрафа в реальные минуты
func convertPenaltyCode(code int) int {
	switch code {
	case 1:
		return 2
	case 2:
		return 4
	default:
		if code >= 5 && code <= 25 {
			return code
		}
		return 0
	}
}

// parseLineupsTable парсит составы из таблицы с 12 колонками
// Структура: №|Поз|?|флаг|игрок|да × 2 команды
func parseLineupsTable(doc *goquery.Document) (home, away []dto.LineupPlayerDTO) {
	// Ищем заголовок "Составы"
	doc.Find("h1").Each(func(_ int, h1 *goquery.Selection) {
		if !strings.Contains(h1.Text(), "Составы") {
			return
		}

		// Ищем таблицу в следующем sibling div
		table := h1.Next().Find("table.table-hover").First()
		if table.Length() == 0 {
			return
		}

		table.Find("tr.out_blue, tr.out_white").Each(func(_ int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() < 12 {
				return
			}

			// Домашняя команда (колонки 0-5)
			homePlayer := parseLineupPlayer(cells, 0)
			if homePlayer != nil {
				home = append(home, *homePlayer)
			}

			// Гостевая команда (колонки 6-11)
			awayPlayer := parseLineupPlayer(cells, 6)
			if awayPlayer != nil {
				away = append(away, *awayPlayer)
			}
		})
	})

	return home, away
}

// parseLineupPlayer парсит игрока из 6 ячеек начиная с offset
// Структура: номер | позиция | ? | флаг | игрок | да
func parseLineupPlayer(cells *goquery.Selection, offset int) *dto.LineupPlayerDTO {
	if offset+5 >= cells.Length() {
		return nil
	}

	number, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(offset).Text()))
	position := strings.TrimSpace(cells.Eq(offset + 1).Text())

	// Игрок в 5-й колонке (offset + 4)
	playerCell := cells.Eq(offset + 4)
	playerLink := playerCell.Find("a[href*='/players/info/']").First()
	if playerLink.Length() == 0 {
		return nil
	}

	href, _ := playerLink.Attr("href")
	m := playerLinkRegex.FindStringSubmatch(href)
	if len(m) < 2 || m[1] == "" {
		return nil
	}

	return &dto.LineupPlayerDTO{
		PlayerID:   m[1],
		PlayerName: strings.TrimSpace(playerLink.Text()),
		Number:     number,
		Position:   position,
	}
}

// calculatePeriod вычисляет период по времени
func calculatePeriod(minutes int) int {
	period := (minutes / 20) + 1
	if period > 4 {
		period = 4
	}
	return period
}

// detectGoalType определяет тип гола
func detectGoalType(text string) string {
	switch {
	case strings.Contains(text, "больш"):
		return "pp"
	case strings.Contains(text, "меньш"):
		return "sh"
	default:
		return "even"
	}
}
