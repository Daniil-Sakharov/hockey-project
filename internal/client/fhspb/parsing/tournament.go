package parsing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

var (
	tournamentIDRegex = regexp.MustCompile(`TournamentID=(\d+)`)
	birthYearRegex    = regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
	dateRangeRegex    = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4})\s*-\s*(\d{2}\.\d{2}\.\d{4})`)
	groupNameRegex    = regexp.MustCompile(`\s*(Группа\s+[А-Яа-яA-Za-z0-9]+)\s*$`)
)

// ParseTournaments парсит список турниров из HTML
func ParseTournaments(html []byte) ([]dto.TournamentDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	tournaments := make(map[int]*dto.TournamentDTO)

	// Парсим каждый div.clearfix как карточку турнира
	doc.Find("div.clearfix").Each(func(_ int, container *goquery.Selection) {
		// Ищем ссылку на турнир внутри h4
		link := container.Find("h4 a[href*='TournamentID=']").First()
		if link.Length() == 0 {
			return
		}

		href, _ := link.Attr("href")
		matches := tournamentIDRegex.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil || id == 0 {
			return
		}

		if _, exists := tournaments[id]; exists {
			return
		}

		name := strings.TrimSpace(link.Text())
		if name == "" {
			return
		}

		// Birth year из родительского h5.subheader (div.clearfix находится внутри незакрытого h5)
		birthYear := 0
		parent := container.Parent()
		if parent.Is("h5.subheader") {
			// Берём только первую строку текста (до вложенных элементов)
			text := parent.Contents().First().Text()
			birthYear = extractBirthYear(text)
		}

		// Даты из span.warning.label внутри контейнера
		var startDate, endDate *time.Time
		container.Find("span.warning.label").Each(func(_ int, span *goquery.Selection) {
			if startDate == nil {
				startDate, endDate = extractDates(span.Text())
			}
		})

		// IsEnded - span.success.label с "Завершен" внутри этого контейнера
		isEnded := false
		container.Find("span.success.label").Each(func(_ int, span *goquery.Selection) {
			if strings.Contains(span.Text(), "Завершен") {
				isEnded = true
			}
		})

		// Извлекаем группу из имени (например "Группа А", "Группа Б")
		groupName := ""
		if m := groupNameRegex.FindStringSubmatch(name); len(m) >= 2 {
			groupName = strings.TrimSpace(m[1])
			name = strings.TrimSpace(groupNameRegex.ReplaceAllString(name, ""))
		}

		tournaments[id] = &dto.TournamentDTO{
			ID:        id,
			Name:      name,
			GroupName: groupName,
			BirthYear: birthYear,
			Season:    determineSeason(startDate),
			StartDate: startDate,
			EndDate:   endDate,
			IsEnded:   isEnded,
		}
	})

	result := make([]dto.TournamentDTO, 0, len(tournaments))
	for _, t := range tournaments {
		result = append(result, *t)
	}
	return result, nil
}

func extractBirthYear(name string) int {
	matches := birthYearRegex.FindStringSubmatch(name)
	if len(matches) < 2 {
		return 0
	}
	year, _ := strconv.Atoi(matches[1])
	return year
}

func extractDates(text string) (*time.Time, *time.Time) {
	matches := dateRangeRegex.FindStringSubmatch(text)
	if len(matches) < 3 {
		return nil, nil
	}

	startDate := parseDate(matches[1])
	endDate := parseDate(matches[2])
	return startDate, endDate
}

func parseDate(s string) *time.Time {
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return nil
	}
	return &t
}

// determineSeason определяет сезон по дате начала турнира
// Сезон начинается в сентябре и заканчивается в апреле-мае следующего года
func determineSeason(startDate *time.Time) string {
	if startDate == nil {
		return ""
	}
	year := startDate.Year()
	month := startDate.Month()

	// Если турнир начинается в сентябре-декабре, это сезон year-year+1
	// Если в январе-августе, это сезон year-1-year
	if month >= 9 {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}

// FilterByBirthYear фильтрует турниры по минимальному году рождения
// minYear = 2008 означает: брать турниры для 2008, 2009, 2010... (молодые игроки)
// Турниры без birth_year (Мужские команды, ЮХЛ и т.д.) пропускаются
func FilterByBirthYear(tournaments []dto.TournamentDTO, minYear int) []dto.TournamentDTO {
	var result []dto.TournamentDTO
	for _, t := range tournaments {
		if t.BirthYear >= minYear {
			result = append(result, t)
		}
	}
	return result
}
