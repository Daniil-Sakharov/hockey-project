package parsing

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseTournamentMetadata извлекает даты и флаг завершенности турнира
func ParseTournamentMetadata(s *goquery.Selection) (startDate, endDate string, isEnded bool) {
	// Ищем в родительском comp-card
	s.Parents().Each(func(j int, parent *goquery.Selection) {
		if !parent.HasClass("comp-card") {
			return
		}

		// Парсим даты из comp-period
		periodText := strings.TrimSpace(parent.Find(".comp-period").Text())
		if periodText != "" {
			// Формат: "с 01.09.2025" или "с 27.04.2025 до 06.05.2025"

			// Убираем "с " в начале
			periodText = strings.TrimPrefix(periodText, "с ")
			periodText = strings.TrimSpace(periodText)

			// Проверяем наличие " до "
			if strings.Contains(periodText, " до ") {
				// Есть дата окончания
				parts := strings.Split(periodText, " до ")
				startDate = strings.TrimSpace(parts[0]) // "01.09.2025"
				if len(parts) >= 2 {
					endDate = strings.TrimSpace(parts[1]) // "06.05.2025"
				}
			} else {
				// Только дата начала
				startDate = periodText // "01.09.2025"
			}
		}

		// Флаг завершенности (класс comp-ended)
		if parent.HasClass("comp-ended") {
			isEnded = true
		}
	})

	return
}
