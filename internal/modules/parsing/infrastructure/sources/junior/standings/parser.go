package standings

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер турнирных таблиц
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// Parse парсит турнирную таблицу (только первый год/группа по умолчанию)
func (p *Parser) Parse(tournamentURL string) ([]StandingDTO, error) {
	resp, err := p.http.MakeRequest(tournamentURL)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	return p.parseStandings(doc, StandingsFilter{})
}

// ParseFromDoc парсит standings из уже загруженного документа
func (p *Parser) ParseFromDoc(doc *goquery.Document, filter StandingsFilter) ([]StandingDTO, error) {
	return p.parseStandings(doc, filter)
}

// ParseWithFilter парсит турнирную таблицу через AJAX с указанием года и группы
func (p *Parser) ParseWithFilter(baseURL, ajaxURL string, filter StandingsFilter) ([]StandingDTO, error) {
	// Формируем полный AJAX URL
	fullURL := ajaxURL
	if !strings.HasPrefix(ajaxURL, "http") {
		fullURL = strings.TrimSuffix(baseURL, "/") + ajaxURL
	}

	resp, err := p.http.MakeRequestWithHeaders(fullURL, map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	})
	if err != nil {
		return nil, fmt.Errorf("ajax request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	return p.parseStandings(doc, filter)
}

func (p *Parser) parseStandings(doc *goquery.Document, filter StandingsFilter) ([]StandingDTO, error) {
	var standings []StandingDTO

	// Ищем таблицу с турнирной таблицей
	// Находим таблицу, содержащую заголовок "М" (место) и ссылки на команды
	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		// Проверяем, что это таблица standings по заголовкам
		headers := p.parseHeaders(table)
		if _, hasPosition := headers["м"]; !hasPosition {
			return
		}

		// Используем группу из фильтра или пытаемся найти на странице
		groupName := filter.GroupName
		if groupName == "" {
			groupName = p.findGroupName(table)
		}

		table.Find("tbody tr, tr").Each(func(j int, row *goquery.Selection) {
			// Пропускаем строки заголовков
			if row.Find("th").Length() > 0 {
				return
			}
			if standing := p.parseRow(row, headers, groupName); standing != nil {
				// Добавляем год рождения из фильтра
				standing.BirthYear = filter.BirthYear
				standings = append(standings, *standing)
			}
		})
	})

	return standings, nil
}

func (p *Parser) findGroupName(table *goquery.Selection) string {
	// Проверяем предыдущие элементы
	prev := table.PrevFiltered("h3, h4, .group-name, .table-title")
	if prev.Length() > 0 {
		return strings.TrimSpace(prev.Text())
	}

	// Проверяем родительский контейнер
	parent := table.Parent()
	title := parent.Find(".group-title, .section-title").First()
	if title.Length() > 0 {
		return strings.TrimSpace(title.Text())
	}

	return ""
}

func (p *Parser) parseHeaders(table *goquery.Selection) map[string]int {
	headers := make(map[string]int)

	// Ищем первую строку с th элементами
	var headerRow *goquery.Selection

	// Сначала проверяем thead
	theadTr := table.Find("thead tr").First()
	if theadTr.Length() > 0 && theadTr.Find("th").Length() > 0 {
		headerRow = theadTr
	} else {
		// Иначе ищем первую tr с th
		table.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
			if tr.Find("th").Length() > 0 {
				headerRow = tr
				return false // break
			}
			return true
		})
	}

	if headerRow == nil {
		return headers
	}

	headerRow.Find("th").Each(func(i int, th *goquery.Selection) {
		text := strings.ToLower(strings.TrimSpace(th.Text()))
		headers[text] = i
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
	standing.Position = p.getIntValue(cells, headers, []string{"м", "место", "#", "pos"}, 0)

	// Команда - "к", "команда", "команды"
	teamCell := cells.Eq(p.getColumnIndex(headers, []string{"к", "команда", "команды", "team", "название"}, 1))
	link := teamCell.Find("a")
	if link.Length() > 0 {
		standing.TeamURL, _ = link.Attr("href")
		standing.TeamName = strings.TrimSpace(link.Text())
	} else {
		standing.TeamName = strings.TrimSpace(teamCell.Text())
	}

	if standing.TeamName == "" {
		return nil
	}

	// Убираем лишние пробелы и переносы в названии
	standing.TeamName = regexp.MustCompile(`\s+`).ReplaceAllString(standing.TeamName, " ")

	// Игры
	standing.Games = p.getIntValue(cells, headers, []string{"и", "игры", "g", "games"}, 2)

	// Победы
	standing.Wins = p.getIntValue(cells, headers, []string{"в", "wins"}, 3)
	standing.WinsOT = p.getIntValue(cells, headers, []string{"во", "wins_ot", "otw"}, 4)
	standing.WinsSO = p.getIntValue(cells, headers, []string{"вб", "wins_so", "sow"}, 5)

	// Поражения
	standing.LossesSO = p.getIntValue(cells, headers, []string{"пб", "losses_so", "sol"}, 6)
	standing.LossesOT = p.getIntValue(cells, headers, []string{"по", "losses_ot", "otl"}, 7)
	standing.Losses = p.getIntValue(cells, headers, []string{"п", "losses"}, 8)

	// Ничьи
	standing.Draws = p.getIntValue(cells, headers, []string{"н", "draws"}, -1)

	// Голы - отдельные колонки ШЗ и ШП
	standing.GoalsFor = p.getIntValue(cells, headers, []string{"шз", "gf", "goals_for"}, 9)
	standing.GoalsAgainst = p.getIntValue(cells, headers, []string{"шп", "ga", "goals_against"}, 10)

	// Если голы в объединённой колонке (РШ формат "119:74")
	if standing.GoalsFor == 0 && standing.GoalsAgainst == 0 {
		goalsText := p.getTextValue(cells, headers, []string{"рш", "ш", "goals", "gf-ga"}, -1)
		standing.GoalsFor, standing.GoalsAgainst = parseGoals(goalsText)
	}

	// Разница
	standing.GoalDifference = p.getIntValue(cells, headers, []string{"+/-", "diff", "разница"}, 11)
	if standing.GoalDifference == 0 {
		standing.GoalDifference = standing.GoalsFor - standing.GoalsAgainst
	}

	// Очки - обычно последняя колонка
	standing.Points = p.getIntValue(cells, headers, []string{"о", "очки", "pts", "points"}, 12)

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

var goalsRegex = regexp.MustCompile(`(\d+)\s*[-:]\s*(\d+)`)

func parseGoals(text string) (goalsFor, goalsAgainst int) {
	matches := goalsRegex.FindStringSubmatch(text)
	if len(matches) >= 3 {
		goalsFor, _ = strconv.Atoi(matches[1])
		goalsAgainst, _ = strconv.Atoi(matches[2])
	}
	return
}
