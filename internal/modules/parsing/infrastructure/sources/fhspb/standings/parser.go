package standings

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	teamIDRegex = regexp.MustCompile(`TeamID=(\d+)`)
	goalsRegex  = regexp.MustCompile(`(\d+)\s*[-:]\s*(\d+)`)
)

// Parser парсер турнирной таблицы FHSPB
type Parser struct{}

// NewParser создает новый парсер
func NewParser() *Parser {
	return &Parser{}
}

// Parse парсит HTML страницы турнирной таблицы
func (p *Parser) Parse(html []byte) ([]StandingDTO, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	return p.parseStandings(doc)
}

func (p *Parser) parseStandings(doc *goquery.Document) ([]StandingDTO, error) {
	var standings []StandingDTO

	// Ищем таблицу с турнирной таблицей
	// Обычно это ctl00_ctl00_MainContent_MainContent_TeamGridView
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		headers := p.parseHeaders(table)
		if len(headers) < 3 {
			return
		}

		// Проверяем, что это турнирная таблица по наличию колонки "М" или позиции
		if _, hasPosition := headers["м"]; !hasPosition {
			if _, hasHash := headers["#"]; !hasHash {
				return
			}
		}

		groupName := p.findGroupName(table)

		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			// Пропускаем заголовки
			if row.Find("th").Length() > 0 || i == 0 {
				return
			}

			standing := p.parseRow(row, headers, groupName)
			if standing != nil {
				standings = append(standings, *standing)
			}
		})
	})

	return standings, nil
}

func (p *Parser) findGroupName(table *goquery.Selection) string {
	// Ищем заголовок группы перед таблицей
	prev := table.Prev()
	for prev.Length() > 0 {
		if prev.Is("h3, h4, h5, .group-name") {
			text := strings.TrimSpace(prev.Text())
			if text != "" && !strings.Contains(strings.ToLower(text), "таблиц") {
				return text
			}
		}
		prev = prev.Prev()
	}
	return ""
}

func (p *Parser) parseHeaders(table *goquery.Selection) map[string]int {
	headers := make(map[string]int)

	// Ищем строку с заголовками
	table.Find("tr").First().Find("th, td").Each(func(i int, cell *goquery.Selection) {
		text := strings.ToLower(strings.TrimSpace(cell.Text()))
		if text != "" {
			headers[text] = i
		}
	})

	return headers
}

func (p *Parser) parseRow(row *goquery.Selection, headers map[string]int, groupName string) *StandingDTO {
	cells := row.Find("td")
	if cells.Length() < 5 {
		return nil
	}

	standing := &StandingDTO{
		GroupName: groupName,
	}

	// Позиция
	standing.Position = p.getIntValue(cells, headers, []string{"м", "#", "место", "pos"}, 0)

	// Команда
	teamCell := cells.Eq(p.getColumnIndex(headers, []string{"команда", "к", "team", "название"}, 1))
	link := teamCell.Find("a[href*='TeamID=']")
	if link.Length() > 0 {
		href, _ := link.Attr("href")
		standing.TeamURL = href
		standing.TeamName = strings.TrimSpace(link.Text())
	} else {
		standing.TeamName = strings.TrimSpace(teamCell.Text())
	}

	if standing.TeamName == "" {
		return nil
	}

	// Чистим название
	standing.TeamName = regexp.MustCompile(`\s+`).ReplaceAllString(standing.TeamName, " ")

	// Игры
	standing.Games = p.getIntValue(cells, headers, []string{"и", "игры", "g", "games"}, 2)

	// Победы
	standing.Wins = p.getIntValue(cells, headers, []string{"в", "wins"}, 3)
	standing.WinsOT = p.getIntValue(cells, headers, []string{"во", "otw"}, -1)
	standing.WinsSO = p.getIntValue(cells, headers, []string{"вб", "sow"}, -1)

	// Поражения
	standing.LossesSO = p.getIntValue(cells, headers, []string{"пб", "sol"}, -1)
	standing.LossesOT = p.getIntValue(cells, headers, []string{"по", "otl"}, -1)
	standing.Losses = p.getIntValue(cells, headers, []string{"п", "losses"}, -1)

	// Ничьи
	standing.Draws = p.getIntValue(cells, headers, []string{"н", "draws"}, -1)

	// Голы - могут быть в отдельных колонках или объединённые
	standing.GoalsFor = p.getIntValue(cells, headers, []string{"шз", "gf"}, -1)
	standing.GoalsAgainst = p.getIntValue(cells, headers, []string{"шп", "ga"}, -1)

	// Если голы в объединённой колонке
	if standing.GoalsFor == 0 && standing.GoalsAgainst == 0 {
		goalsText := p.getTextValue(cells, headers, []string{"рш", "ш", "goals"}, -1)
		standing.GoalsFor, standing.GoalsAgainst = parseGoals(goalsText)
	}

	// Разница
	standing.GoalDifference = p.getIntValue(cells, headers, []string{"+/-", "diff"}, -1)
	if standing.GoalDifference == 0 {
		standing.GoalDifference = standing.GoalsFor - standing.GoalsAgainst
	}

	// Очки
	standing.Points = p.getIntValue(cells, headers, []string{"о", "очки", "pts", "points"}, -1)

	return standing
}

func (p *Parser) getColumnIndex(headers map[string]int, names []string, defaultIdx int) int {
	for _, name := range names {
		if idx, ok := headers[name]; ok {
			return idx
		}
	}
	return defaultIdx
}

func (p *Parser) getIntValue(cells *goquery.Selection, headers map[string]int, names []string, defaultIdx int) int {
	idx := p.getColumnIndex(headers, names, defaultIdx)
	if idx < 0 || idx >= cells.Length() {
		return 0
	}
	text := strings.TrimSpace(cells.Eq(idx).Text())
	val, _ := strconv.Atoi(text)
	return val
}

func (p *Parser) getTextValue(cells *goquery.Selection, headers map[string]int, names []string, defaultIdx int) string {
	idx := p.getColumnIndex(headers, names, defaultIdx)
	if idx < 0 || idx >= cells.Length() {
		return ""
	}
	return strings.TrimSpace(cells.Eq(idx).Text())
}

func parseGoals(text string) (goalsFor, goalsAgainst int) {
	matches := goalsRegex.FindStringSubmatch(text)
	if len(matches) >= 3 {
		goalsFor, _ = strconv.Atoi(matches[1])
		goalsAgainst, _ = strconv.Atoi(matches[2])
	}
	return
}
