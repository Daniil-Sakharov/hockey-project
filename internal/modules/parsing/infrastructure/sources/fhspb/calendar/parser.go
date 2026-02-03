package calendar

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	matchIDRegex    = regexp.MustCompile(`MatchID=(\d+)`)
	scoreRegex      = regexp.MustCompile(`(\d+)\s*:\s*(\d+)(ОТ|ПБ)?`)
	dateRegex       = regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{4})`)
	timeRegex       = regexp.MustCompile(`(\d{2}):(\d{2})`)
	teamMatchRegex  = regexp.MustCompile(`(.+?)\s*-\s*(.+)`)
	matchNumRegex   = regexp.MustCompile(`^\d+$`)
)

// Parser парсер календаря матчей FHSPB
type Parser struct{}

// NewParser создает новый парсер календаря
func NewParser() *Parser {
	return &Parser{}
}

// Parse парсит HTML страницы календаря и возвращает список матчей
func (p *Parser) Parse(html []byte, tournamentID int) ([]MatchDTO, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	return p.parseMatches(doc)
}

func (p *Parser) parseMatches(doc *goquery.Document) ([]MatchDTO, error) {
	var matches []MatchDTO
	seen := make(map[string]bool)

	// Таблица матчей: MatchGridView
	doc.Find("#MatchGridView tr").Each(func(i int, row *goquery.Selection) {
		// Пропускаем заголовок
		if i == 0 {
			return
		}

		match := p.parseRow(row)
		if match == nil {
			return
		}

		// Пропускаем дубликаты
		if seen[match.ExternalID] {
			return
		}
		seen[match.ExternalID] = true

		matches = append(matches, *match)
	})

	return matches, nil
}

func (p *Parser) parseRow(row *goquery.Selection) *MatchDTO {
	cells := row.Find("td")
	if cells.Length() < 6 {
		return nil
	}

	match := &MatchDTO{}

	// Извлекаем MatchID из любой ссылки в строке
	row.Find("a[href*='MatchID=']").Each(func(_ int, link *goquery.Selection) {
		if match.ExternalID == "" {
			href, _ := link.Attr("href")
			if m := matchIDRegex.FindStringSubmatch(href); len(m) >= 2 {
				match.ExternalID = m[1]
			}
		}
	})

	if match.ExternalID == "" {
		return nil
	}

	// Парсим ячейки
	// Структура: [турнир], №, Дата, Время, Стадион, Матч, Результат, [протокол]
	cells.Each(func(i int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())

		switch {
		case matchNumRegex.MatchString(text) && match.MatchNumber == 0:
			// Номер матча
			if num, err := strconv.Atoi(text); err == nil {
				match.MatchNumber = num
			}

		case dateRegex.MatchString(text) && match.ScheduledAt == nil:
			// Дата
			match.ScheduledAt = p.parseDate(text)

		case timeRegex.MatchString(text) && len(text) <= 5:
			// Время (обновляем дату)
			match.ScheduledAt = p.addTime(match.ScheduledAt, text)

		case len(text) <= 5 && !strings.Contains(text, ":") && match.Venue == "" && i > 3:
			// Код стадиона (3-5 букв)
			if len(text) >= 2 && len(text) <= 5 && isArenaCode(text) {
				match.Venue = text
			}

		case teamMatchRegex.MatchString(text) && match.HomeTeamName == "":
			// Названия команд "Команда А - Команда Б"
			p.parseTeams(text, match)

		case scoreRegex.MatchString(text):
			// Результат "7 : 2" или "5 : 4ОТ"
			p.parseScore(text, match)
		}
	})

	// Если нет команд, но есть ссылки на них
	if match.HomeTeamName == "" {
		p.parseTeamsFromLinks(row, match)
	}

	if match.HomeTeamName == "" || match.AwayTeamName == "" {
		return nil
	}

	return match
}

func (p *Parser) parseDate(text string) *time.Time {
	m := dateRegex.FindStringSubmatch(text)
	if len(m) < 4 {
		return nil
	}

	day, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	year, _ := strconv.Atoi(m[3])

	loc, _ := time.LoadLocation("Europe/Moscow")
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	return &t
}

func (p *Parser) addTime(date *time.Time, timeStr string) *time.Time {
	if date == nil {
		return nil
	}

	m := timeRegex.FindStringSubmatch(timeStr)
	if len(m) < 3 {
		return date
	}

	hour, _ := strconv.Atoi(m[1])
	minute, _ := strconv.Atoi(m[2])

	t := time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location())
	return &t
}

func (p *Parser) parseTeams(text string, match *MatchDTO) {
	m := teamMatchRegex.FindStringSubmatch(text)
	if len(m) < 3 {
		return
	}

	match.HomeTeamName = strings.TrimSpace(m[1])
	match.AwayTeamName = strings.TrimSpace(m[2])
}

func (p *Parser) parseTeamsFromLinks(row *goquery.Selection, match *MatchDTO) {
	teams := []string{}
	row.Find("a[href*='TeamID=']").Each(func(_ int, link *goquery.Selection) {
		name := strings.TrimSpace(link.Text())
		if name != "" && len(teams) < 2 {
			teams = append(teams, name)
		}
	})

	if len(teams) >= 2 {
		match.HomeTeamName = teams[0]
		match.AwayTeamName = teams[1]
	}
}

func (p *Parser) parseScore(text string, match *MatchDTO) {
	m := scoreRegex.FindStringSubmatch(text)
	if len(m) < 3 {
		return
	}

	home, _ := strconv.Atoi(m[1])
	away, _ := strconv.Atoi(m[2])

	match.HomeScore = &home
	match.AwayScore = &away
	match.IsFinished = true

	// Тип результата
	if len(m) >= 4 {
		switch m[3] {
		case "ОТ":
			match.ResultType = "OT"
		case "ПБ":
			match.ResultType = "SO"
		}
	}
}

func isArenaCode(s string) bool {
	for _, r := range s {
		if !((r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false
		}
	}
	return true
}
