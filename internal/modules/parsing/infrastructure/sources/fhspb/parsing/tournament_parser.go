package parsing

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/dto"
	"github.com/PuerkitoBio/goquery"
)

func parseTournamentCard(container *goquery.Selection) *dto.TournamentDTO {
	link := container.Find("h4 a[href*='TournamentID=']").First()
	if link.Length() == 0 {
		return nil
	}

	href, _ := link.Attr("href")
	matches := tournamentIDRegex.FindStringSubmatch(href)
	if len(matches) < 2 {
		return nil
	}

	id, err := strconv.Atoi(matches[1])
	if err != nil || id == 0 {
		return nil
	}

	name := strings.TrimSpace(link.Text())
	if name == "" {
		return nil
	}

	// Извлекаем год рождения:
	// 1. Сначала пробуем из названия турнира (например "... 2008 г.р.")
	// 2. Если не нашли - ищем в родительском h5.subheader
	birthYear := extractBirthYear(name)
	if birthYear == 0 {
		// Ищем ближайший h5.subheader с годом рождения
		parent := container.Parent()
		if parent.Is("h5.subheader") {
			text := parent.Contents().First().Text()
			birthYear = extractBirthYear(text)
		}
		// Также проверяем предыдущие sibling элементы
		if birthYear == 0 {
			container.PrevAll().FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.Is("h5.subheader")
			}).First().Each(func(_ int, h5 *goquery.Selection) {
				if birthYear == 0 {
					birthYear = extractBirthYear(h5.Text())
				}
			})
		}
	}

	var startDate, endDate *time.Time
	// Bootstrap использует классы .label.label-warning, ищем по label-warning
	container.Find("span.label-warning, span[class*='warning']").Each(func(_ int, span *goquery.Selection) {
		if startDate == nil {
			startDate, endDate = extractDates(span.Text())
		}
	})

	isEnded := false
	// Bootstrap использует классы .label.label-success
	container.Find("span.label-success, span[class*='success']").Each(func(_ int, span *goquery.Selection) {
		if strings.Contains(span.Text(), "Завершен") {
			isEnded = true
		}
	})

	groupName := ""
	if m := groupNameRegex.FindStringSubmatch(name); len(m) >= 2 {
		groupName = strings.TrimSpace(m[1])
		name = strings.TrimSpace(groupNameRegex.ReplaceAllString(name, ""))
	}

	return &dto.TournamentDTO{
		ID:        id,
		Name:      name,
		GroupName: groupName,
		BirthYear: birthYear,
		Season:    determineSeason(startDate),
		StartDate: startDate,
		EndDate:   endDate,
		IsEnded:   isEnded,
	}
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
	return parseDate(matches[1]), parseDate(matches[2])
}

func parseDate(s string) *time.Time {
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return nil
	}
	return &t
}

func determineSeason(startDate *time.Time) string {
	if startDate == nil {
		return ""
	}
	year := startDate.Year()
	month := startDate.Month()

	if month >= 9 {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}
