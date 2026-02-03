package parsing

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
	"github.com/PuerkitoBio/goquery"
)

var (
	// Regex для парсинга данных
	// Ищем "Дата рождения: DD Month YYYY" или просто "DD Month YYYY" рядом с датой
	birthDateRegex = regexp.MustCompile(`[Дд]ата\s+рождения[:\s]+(\d{1,2})\s+(\S+)\s+(\d{4})`)
	dateRegexRu    = regexp.MustCompile(`(\d{1,2})\s+(\S+)\s+(\d{4})`)
	heightRegex    = regexp.MustCompile(`(\d+)\s*(?:см)?`)
	weightRegex    = regexp.MustCompile(`(\d+)\s*(?:кг)?`)
	teamLinkRegex  = regexp.MustCompile(`/team/(\d+)`)
	iceTimeRegex   = regexp.MustCompile(`(\d+):(\d+)`)

	// Карта месяцев на русском (lowercase)
	monthsRu = map[string]time.Month{
		"январь": time.January, "января": time.January,
		"февраль": time.February, "февраля": time.February,
		"март": time.March, "марта": time.March,
		"апрель": time.April, "апреля": time.April,
		"май": time.May, "мая": time.May,
		"июнь": time.June, "июня": time.June,
		"июль": time.July, "июля": time.July,
		"август": time.August, "августа": time.August,
		"сентябрь": time.September, "сентября": time.September,
		"октябрь": time.October, "октября": time.October,
		"ноябрь": time.November, "ноября": time.November,
		"декабрь": time.December, "декабря": time.December,
		// Capitalized versions (as they appear on the website)
		"Январь": time.January, "Января": time.January,
		"Февраль": time.February, "Февраля": time.February,
		"Март": time.March, "Марта": time.March,
		"Апрель": time.April, "Апреля": time.April,
		"Май": time.May, "Мая": time.May,
		"Июнь": time.June, "Июня": time.June,
		"Июль": time.July, "Июля": time.July,
		"Август": time.August, "Августа": time.August,
		"Сентябрь": time.September, "Сентября": time.September,
		"Октябрь": time.October, "Октября": time.October,
		"Ноябрь": time.November, "Ноября": time.November,
		"Декабрь": time.December, "Декабря": time.December,
	}
)

// ParsePlayerProfile парсит страницу профиля игрока /player/{id}
func ParsePlayerProfile(html []byte, playerID string) (*dto.PlayerProfileDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	profile := &dto.PlayerProfileDTO{
		ID: playerID,
	}

	// Парсим имя игрока (обычно в h1 или в блоке с классом)
	doc.Find("h1, .player-name, .name").Each(func(_ int, s *goquery.Selection) {
		if profile.FullName == "" {
			text := strings.TrimSpace(s.Text())
			if text != "" && !strings.Contains(strings.ToLower(text), "статистика") {
				profile.FullName = text
			}
		}
	})

	// Парсим данные профиля из текста страницы
	pageText := doc.Text()
	profile.BirthDate = parseBirthDateRu(pageText)
	if profile.BirthDate != nil {
		profile.Age = calculateAge(*profile.BirthDate)
	}

	// Ищем позицию по классу .position
	doc.Find(".position").Each(func(_ int, s *goquery.Selection) {
		if profile.Position == "" {
			text := strings.TrimSpace(s.Text())
			profile.Position = extractPosition(text)
		}
	})

	// Ищем данные в блоках с лейблами
	doc.Find("div, span, p").Each(func(_ int, s *goquery.Selection) {
		text := strings.ToLower(strings.TrimSpace(s.Text()))

		// Fallback для позиции, если не нашли по классу
		if profile.Position == "" && (strings.Contains(text, "позиция") || strings.Contains(text, "амплуа")) {
			profile.Position = extractPosition(text)
		}
		if strings.Contains(text, "рост") {
			profile.Height = extractNumber(text, heightRegex)
		}
		if strings.Contains(text, "вес") {
			profile.Weight = extractNumber(text, weightRegex)
		}
		if strings.Contains(text, "хват") {
			profile.Handedness = extractHandedness(text)
		}
	})

	// Парсим таблицу статистики
	doc.Find("table").Each(func(_ int, table *goquery.Selection) {
		// Проверяем что это таблица статистики
		headerText := table.Find("th, thead").Text()
		if !strings.Contains(strings.ToLower(headerText), "сезон") &&
			!strings.Contains(strings.ToLower(headerText), "игры") {
			return
		}

		table.Find("tbody tr, tr").Each(func(_ int, row *goquery.Selection) {
			stats := parseStatsRow(row)
			if stats != nil {
				profile.Stats = append(profile.Stats, *stats)
			}
		})
	})

	return profile, nil
}

func parseBirthDateRu(text string) *time.Time {
	// Сначала пробуем найти конкретно "Дата рождения: DD Month YYYY"
	matches := birthDateRegex.FindStringSubmatch(text)
	if len(matches) == 4 {
		day, _ := strconv.Atoi(matches[1])
		monthStr := strings.ToLower(matches[2])
		year, _ := strconv.Atoi(matches[3])

		if month, ok := monthsRu[monthStr]; ok {
			t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			return &t
		}
	}

	// Fallback: ищем любой паттерн "DD Month YYYY" с месяцем из map
	allMatches := dateRegexRu.FindAllStringSubmatch(text, -1)
	for _, m := range allMatches {
		if len(m) != 4 {
			continue
		}
		monthStr := strings.ToLower(m[2])
		if month, ok := monthsRu[monthStr]; ok {
			day, _ := strconv.Atoi(m[1])
			year, _ := strconv.Atoi(m[3])
			// Проверяем что год похож на год рождения (2005-2020)
			if year >= 2005 && year <= 2020 {
				t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
				return &t
			}
		}
	}

	return nil
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	return age
}

func extractPosition(text string) string {
	text = strings.ToLower(text)
	if strings.Contains(text, "вратарь") {
		return "В"
	}
	if strings.Contains(text, "защитник") {
		return "З"
	}
	if strings.Contains(text, "нападающий") {
		return "Н"
	}
	return ""
}

func extractHandedness(text string) string {
	text = strings.ToLower(text)
	if strings.Contains(text, "левый") || strings.Contains(text, "левая") {
		return "Л"
	}
	if strings.Contains(text, "правый") || strings.Contains(text, "правая") {
		return "П"
	}
	return ""
}

func extractNumber(text string, re *regexp.Regexp) int {
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		n, _ := strconv.Atoi(matches[1])
		return n
	}
	return 0
}

func parseStatsRow(row *goquery.Selection) *dto.PlayerStatsDTO {
	cells := row.Find("td")
	if cells.Length() < 7 {
		return nil
	}

	stats := &dto.PlayerStatsDTO{}

	// Порядок колонок: Команда, Сезон, Турнир, И, Г, А, О, Ш, ШП, МИН/СЕК
	cells.Each(func(i int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())

		switch i {
		case 0: // Команда
			stats.TeamName = text
			// Ищем ссылку на команду
			if link, exists := cell.Find("a").Attr("href"); exists {
				if matches := teamLinkRegex.FindStringSubmatch(link); len(matches) > 1 {
					stats.TeamID, _ = strconv.Atoi(matches[1])
				}
			}
		case 1: // Сезон
			stats.Season = text
		case 2: // Турнир
			stats.TournamentName = text
		case 3: // И - игры
			stats.Games, _ = strconv.Atoi(text)
		case 4: // Г - голы
			stats.Goals, _ = strconv.Atoi(text)
		case 5: // А - передачи
			stats.Assists, _ = strconv.Atoi(text)
		case 6: // О - очки
			stats.Points, _ = strconv.Atoi(text)
		case 7: // Ш - штрафы (количество)
			stats.PenaltyCount, _ = strconv.Atoi(text)
		case 8: // ШП - штрафные минуты
			stats.PenaltyMinutes, _ = strconv.Atoi(text)
		case 9: // МИН/СЕК - время на льду
			stats.IceTime = text
			if matches := iceTimeRegex.FindStringSubmatch(text); len(matches) == 3 {
				mins, _ := strconv.Atoi(matches[1])
				secs, _ := strconv.Atoi(matches[2])
				stats.IceTimeSeconds = mins*60 + secs
			}
		}
	})

	// Проверяем что спарсили минимально необходимые данные
	if stats.TeamName == "" {
		return nil
	}

	return stats
}
