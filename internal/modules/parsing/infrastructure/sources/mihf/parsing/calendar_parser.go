package parsing

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/PuerkitoBio/goquery"
)

var (
	calendarTeamIDRegex = regexp.MustCompile(`/team/(\d+)`)
	calendarProtoRegex  = regexp.MustCompile(`/proto/(\d+)`)
	scoreRegex          = regexp.MustCompile(`^(\d+):(\d+)$`)
)

// ParseCalendar парсит календарь матчей турнира
func ParseCalendar(html []byte) ([]dto.MatchDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	var matches []dto.MatchDTO
	currentRound := 0

	doc.Find("table.table-hover tr").Each(func(_ int, row *goquery.Selection) {
		// Проверяем строку-заголовок тура
		if dataRound, exists := row.Attr("data-round"); exists {
			currentRound, _ = strconv.Atoi(dataRound)
			return
		}

		// Пропускаем заголовки таблицы
		if row.Find("th").Length() > 0 {
			return
		}

		cells := row.Find("td")
		if cells.Length() < 8 {
			return
		}

		match := parseMatchRow(cells, currentRound)
		if match != nil {
			matches = append(matches, *match)
		}
	})

	return matches, nil
}

// parseMatchRow парсит одну строку матча
func parseMatchRow(cells *goquery.Selection, round int) *dto.MatchDTO {
	match := &dto.MatchDTO{Round: round}

	// Колонка 0: номер матча
	match.MatchNumber, _ = strconv.Atoi(strings.TrimSpace(cells.Eq(0).Text()))

	// Колонка 2: команда A (домашняя)
	teamALink := cells.Eq(2).Find("a[href*='/team/']")
	if href, exists := teamALink.Attr("href"); exists {
		if m := calendarTeamIDRegex.FindStringSubmatch(href); len(m) > 1 {
			match.HomeTeamID = m[1]
		}
	}
	match.HomeTeamName = strings.TrimSpace(teamALink.Text())

	// Колонка 3: команда B (гостевая)
	teamBLink := cells.Eq(3).Find("a[href*='/team/']")
	if href, exists := teamBLink.Attr("href"); exists {
		if m := calendarTeamIDRegex.FindStringSubmatch(href); len(m) > 1 {
			match.AwayTeamID = m[1]
		}
	}
	match.AwayTeamName = strings.TrimSpace(teamBLink.Text())

	// Колонка 4: счет и ссылка на протокол
	scoreCell := cells.Eq(4)
	scoreLink := scoreCell.Find("a")
	if href, exists := scoreLink.Attr("href"); exists {
		match.ProtoURL = href
		if m := calendarProtoRegex.FindStringSubmatch(href); len(m) > 1 {
			match.ExternalID = m[1]
		}
	}
	parseScore(strings.TrimSpace(scoreLink.Text()), match)

	// Колонка 5: дата (DD.MM.YYYY)
	dateStr := strings.TrimSpace(cells.Eq(5).Text())
	// Колонка 6: время (HH:MM)
	timeStr := strings.TrimSpace(cells.Eq(6).Text())
	match.ScheduledAt = parseDateTime(dateStr, timeStr)

	// Колонка 7: стадион + город
	match.Venue = strings.TrimSpace(cells.Eq(7).Text())
	match.VenueCity = extractCityFromVenue(match.Venue)

	// Валидация
	if match.HomeTeamID == "" || match.AwayTeamID == "" {
		return nil
	}

	return match
}

// parseScore парсит счет матча
func parseScore(scoreStr string, match *dto.MatchDTO) {
	if m := scoreRegex.FindStringSubmatch(scoreStr); len(m) == 3 {
		match.HomeScore, _ = strconv.Atoi(m[1])
		match.AwayScore, _ = strconv.Atoi(m[2])
	}
}

// parseDateTime парсит дату и время матча
func parseDateTime(dateStr, timeStr string) time.Time {
	datetime := dateStr
	if timeStr != "" {
		datetime += " " + timeStr
	}
	t, _ := time.Parse("02.01.2006 15:04", datetime)
	return t
}

// extractCityFromVenue извлекает город из названия стадиона
// "СК Локомотив, Ярославль" → "Ярославль"
func extractCityFromVenue(venue string) string {
	parts := strings.Split(venue, ",")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return ""
}

// FindTournamentDates определяет даты начала и окончания турнира
func FindTournamentDates(matches []dto.MatchDTO) (start, end time.Time) {
	for _, m := range matches {
		if m.ScheduledAt.IsZero() {
			continue
		}
		if start.IsZero() || m.ScheduledAt.Before(start) {
			start = m.ScheduledAt
		}
		if m.ScheduledAt.After(end) {
			end = m.ScheduledAt
		}
	}
	return
}
