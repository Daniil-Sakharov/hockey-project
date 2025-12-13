package parsing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
)

var (
	tournamentIDRegex = regexp.MustCompile(`TournamentID=(\d+)`)
	birthYearRegex    = regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
	seasonRegex       = regexp.MustCompile(`Сезон\s+(\d{4}-\d{4})`)
	dateRangeRegex    = regexp.MustCompile(`(\d{2}\.\d{2}\.\d{4})\s*-\s*(\d{2}\.\d{2}\.\d{4})`)
)

// ParseTournaments парсит список турниров из HTML
func ParseTournaments(html []byte) ([]dto.TournamentDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	tournaments := make(map[int]*dto.TournamentDTO)

	// Парсим турниры из ссылок с TournamentID
	doc.Find("a[href*='TournamentID=']").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		matches := tournamentIDRegex.FindStringSubmatch(href)
		if len(matches) < 2 {
			return
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			return
		}

		// Пропускаем если уже обработали этот турнир
		if _, exists := tournaments[id]; exists {
			return
		}

		name := strings.TrimSpace(s.Text())
		if name == "" {
			return
		}

		// Извлекаем даты из title атрибута
		title, _ := s.Attr("title")
		startDate, endDate := extractDates(title)

		// Если в title нет дат, ищем в родительском контейнере
		if startDate == nil {
			parent := s.Parent()
			for i := 0; i < 5 && parent.Length() > 0; i++ {
				parentText := parent.Text()
				startDate, endDate = extractDates(parentText)
				if startDate != nil {
					break
				}
				parent = parent.Parent()
			}
		}

		// Определяем isEnded - ищем в ближайшем контейнере .clearfix
		isEnded := false
		// Поднимаемся до div.clearfix
		container := s.Closest("div.clearfix")
		if container.Length() > 0 {
			// Ищем span.success.label с текстом "Завершен"
			container.Find("span.success.label").Each(func(_ int, span *goquery.Selection) {
				if strings.Contains(span.Text(), "Завершен") {
					isEnded = true
				}
			})
		}

		// Определяем сезон по датам турнира
		season := determineSeason(startDate)

		tournaments[id] = &dto.TournamentDTO{
			ID:        id,
			Name:      name,
			BirthYear: extractBirthYear(name),
			Season:    season,
			StartDate: startDate,
			EndDate:   endDate,
			IsEnded:   isEnded,
		}
	})

	// Конвертируем map в slice
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
func FilterByBirthYear(tournaments []dto.TournamentDTO, minYear int) []dto.TournamentDTO {
	var result []dto.TournamentDTO
	for _, t := range tournaments {
		if t.BirthYear == 0 || t.BirthYear >= minYear {
			result = append(result, t)
		}
	}
	return result
}
