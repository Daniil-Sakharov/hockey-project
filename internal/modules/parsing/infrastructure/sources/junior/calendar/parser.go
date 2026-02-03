package calendar

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер календаря матчей
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер календаря
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// Parse парсит календарь турнира (только первый год/группа по умолчанию)
func (p *Parser) Parse(tournamentURL string) ([]MatchDTO, error) {
	calendarURL := buildCalendarURL(tournamentURL)

	resp, err := p.http.MakeRequest(calendarURL)
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

	return p.parseMatches(doc, CalendarFilter{})
}

// buildCalendarURL корректно строит URL календаря, вставляя /calendar/ перед query параметрами
func buildCalendarURL(tournamentURL string) string {
	parsed, err := url.Parse(tournamentURL)
	if err != nil {
		// Fallback на старую логику
		return strings.TrimSuffix(tournamentURL, "/") + "/calendar/"
	}

	// Вставляем /calendar/ в путь
	path := strings.TrimSuffix(parsed.Path, "/") + "/calendar/"
	parsed.Path = path

	return parsed.String()
}

// ParseFromDoc парсит календарь из уже загруженного документа
func (p *Parser) ParseFromDoc(doc *goquery.Document, filter CalendarFilter) ([]MatchDTO, error) {
	return p.parseMatches(doc, filter)
}

// ParseWithFilter парсит календарь через AJAX с указанием года и группы
func (p *Parser) ParseWithFilter(baseURL, ajaxURL string, filter CalendarFilter) ([]MatchDTO, error) {
	// Формируем полный AJAX URL
	fullURL := ajaxURL
	if !strings.HasPrefix(ajaxURL, "http") {
		fullURL = strings.TrimSuffix(baseURL, "/") + ajaxURL
	}

	// Заменяем tournament-page на competitions-calendar для календаря
	calendarAjaxURL := strings.Replace(fullURL, "tournament-page", "competitions-calendar", 1)

	resp, err := p.http.MakeRequestWithHeaders(calendarAjaxURL, map[string]string{
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

	return p.parseMatches(doc, filter)
}

// ParseCalendarPage парсит страницу календаря и возвращает все матчи для всех групп
func (p *Parser) ParseCalendarPage(tournamentURL string) ([]MatchDTO, error) {
	calendarURL := buildCalendarURL(tournamentURL)

	resp, err := p.http.MakeRequest(calendarURL)
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

	return p.parseMatches(doc, CalendarFilter{})
}

func (p *Parser) parseMatches(doc *goquery.Document, filter CalendarFilter) ([]MatchDTO, error) {
	var matches []MatchDTO
	seen := make(map[string]bool)

	// Ищем все ссылки на игры
	doc.Find("a[href*='/games/']").Each(func(i int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		externalID := extractGameID(href)
		if externalID == "" {
			return
		}

		// Пропускаем дубликаты
		if seen[externalID] {
			return
		}
		seen[externalID] = true

		// Находим родительскую строку таблицы
		row := link.Closest("tr")
		if row.Length() == 0 {
			return
		}

		match := p.parseTableRow(row, externalID, href)
		if match != nil {
			// Добавляем информацию из фильтра
			match.GroupName = filter.GroupName
			match.BirthYear = filter.BirthYear
			matches = append(matches, *match)
		}
	})

	return matches, nil
}

func (p *Parser) parseTableRow(row *goquery.Selection, externalID, gameURL string) *MatchDTO {
	cells := row.Find("td")
	if cells.Length() < 5 {
		return nil
	}

	match := &MatchDTO{
		ExternalID: externalID,
		GameURL:    gameURL,
		Status:     "scheduled",
	}

	// Парсим ячейки в порядке: №, Дата, Команда А, Команда Б, Счёт, Арена
	cells.Each(func(i int, cell *goquery.Selection) {
		switch i {
		case 0:
			// Номер матча
			text := strings.TrimSpace(cell.Find(".cell").Text())
			if num, err := strconv.Atoi(text); err == nil {
				match.MatchNumber = &num
			}
		case 1:
			// Дата и время из отдельных span элементов
			dateText := strings.TrimSpace(cell.Find("span.date").Text())
			timeText := strings.TrimSpace(cell.Find("span.time").Text())
			match.ScheduledAt = parseDateTimeFromParts(dateText, timeText)
		case 2:
			// Домашняя команда
			match.HomeTeam = p.parseTeamFromCell(cell)
		case 3:
			// Гостевая команда
			match.AwayTeam = p.parseTeamFromCell(cell)
		case 4:
			// Счёт
			scoreText := strings.TrimSpace(cell.Find(".cell").Text())
			if home, away, resultType, ok := parseScore(scoreText); ok {
				match.HomeScore = &home
				match.AwayScore = &away
				match.ResultType = resultType
				match.Status = "finished"
			}
		case 5:
			// Арена
			match.Venue = strings.TrimSpace(cell.Find(".cell").Text())
		}
	})

	return match
}

func (p *Parser) parseTeamFromCell(cell *goquery.Selection) TeamInfo {
	info := TeamInfo{}

	// 1. Извлекаем ID команды из URL лого: /upload/team_logo/658725.png → 658725
	logoImg := cell.Find("img[src*='team_logo']")
	if logoImg.Length() > 0 {
		if src, exists := logoImg.Attr("src"); exists {
			info.ID = extractTeamIDFromLogo(src)
		}
	}

	// 2. Ищем ссылку на команду (fallback)
	link := cell.Find("a[href*='/teams/']")
	if link.Length() > 0 {
		if href, exists := link.Attr("href"); exists {
			info.URL = href
			// Если ID ещё не извлечён, пробуем из URL
			if info.ID == "" {
				info.ID = extractTeamIDFromURL(href)
			}
		}
	}

	// 3. Извлекаем название и город
	teamTitle := cell.Find("span.team-title")
	teamCity := cell.Find("span.team-city")

	if teamTitle.Length() > 0 {
		info.Name = strings.TrimSpace(teamTitle.Text())
	}

	if teamCity.Length() > 0 {
		city := strings.TrimSpace(teamCity.Text())
		if city != "" {
			info.Name = info.Name + " " + city
		}
	}

	// Fallback: если нет структурированных данных, берём текст ячейки
	if info.Name == "" {
		info.Name = strings.TrimSpace(cell.Text())
		info.Name = regexp.MustCompile(`\s+`).ReplaceAllString(info.Name, " ")
	}

	return info
}

// extractTeamIDFromLogo извлекает ID команды из URL лого
// /upload/team_logo/658725.png → 658725
// /upload/upload-webp/upload/team_logo/658721-70.webp → 658721
var teamLogoIDRegex = regexp.MustCompile(`team_logo/(\d+)`)

func extractTeamIDFromLogo(logoURL string) string {
	matches := teamLogoIDRegex.FindStringSubmatch(logoURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractTeamIDFromURL извлекает ID команды из URL
// /teams/lada_651237/ → 651237
var teamURLIDRegex = regexp.MustCompile(`_(\d+)/?$`)

func extractTeamIDFromURL(url string) string {
	matches := teamURLIDRegex.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

var gameIDRegex = regexp.MustCompile(`/games/(\d+)`)

func extractGameID(url string) string {
	matches := gameIDRegex.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// parseScore парсит счёт и тип результата (ОТ/Б)
// Возвращает: home, away, resultType ("regular", "OT", "SO"), ok
func parseScore(text string) (home, away int, resultType string, ok bool) {
	text = strings.TrimSpace(text)
	resultType = "regular"

	// Проверяем наличие ОТ/Б в тексте
	if regexp.MustCompile(`(?i)\(ОТ\)|ОТ\s*$`).MatchString(text) {
		resultType = "OT"
	} else if regexp.MustCompile(`(?i)\(Б\)|\(ПБ\)|Б\s*$|ПБ\s*$`).MatchString(text) {
		resultType = "SO"
	}

	// Убираем доп. информацию типа "(ОТ)", "(Б)", "ОТ", "ПБ", "Б"
	text = regexp.MustCompile(`\s*\([^)]*\)\s*`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\s+(ОТ|ПБ|Б)\s*$`).ReplaceAllString(text, "")

	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return 0, 0, "", false
	}

	h, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	a, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return 0, 0, "", false
	}

	return h, a, resultType, true
}

func parseDateTime(text string) *time.Time {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}

	// Убираем "МСК" и лишние пробелы
	text = strings.ReplaceAll(text, "МСК", "")
	text = strings.TrimSpace(text)

	formats := []string{
		"02.01.2006 15:04",
		"02.01.2006",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	for _, format := range formats {
		if t, err := time.ParseInLocation(format, text, loc); err == nil {
			return &t
		}
	}
	return nil
}

// parseDateTimeFromParts парсит дату и время из отдельных строк
func parseDateTimeFromParts(dateStr, timeStr string) *time.Time {
	dateStr = strings.TrimSpace(dateStr)
	timeStr = strings.TrimSpace(timeStr)

	if dateStr == "" {
		return nil
	}

	// Убираем "МСК" из времени
	timeStr = strings.ReplaceAll(timeStr, "МСК", "")
	timeStr = strings.TrimSpace(timeStr)

	// Собираем полную строку
	fullStr := dateStr
	if timeStr != "" {
		fullStr = dateStr + " " + timeStr
	}

	return parseDateTime(fullStr)
}
