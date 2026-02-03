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
	dateRegex   = regexp.MustCompile(`(\d{2})\.(\d{2})\.(\d{4})`)
	heightRegex = regexp.MustCompile(`(\d+)\s*см`)
	weightRegex = regexp.MustCompile(`(\d+)\s*кг`)
)

// ParsePlayerProfile парсит профиль игрока с антропометрией
func ParsePlayerProfile(html []byte, playerID string) (*dto.PlayerProfileDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}

	profile := &dto.PlayerProfileDTO{
		ID: playerID,
	}

	// Парсим текст страницы для извлечения данных
	pageText := doc.Text()

	// Имя обычно в заголовке
	// Важно: первый h1 на странице - "Статистика", нужно его пропустить
	doc.Find("h1, h2, .player-name").Each(func(_ int, s *goquery.Selection) {
		if profile.FullName == "" {
			text := strings.TrimSpace(s.Text())
			// Пропускаем служебные заголовки
			if text != "" && text != "Статистика" && text != "Федерация хоккея Москвы" {
				profile.FullName = text
			}
		}
	})

	// Ищем данные в элементах <strong>label:</strong> value
	// Это основной формат данных на страницах MIHF
	doc.Find("strong").Each(func(_ int, s *goquery.Selection) {
		label := strings.ToLower(strings.TrimSpace(s.Text()))
		// Получаем текст после </strong> до <br>
		parent := s.Parent()
		if parent == nil {
			return
		}
		fullText := parent.Text()

		// Дата рождения
		if strings.Contains(label, "рождения") {
			if matches := dateRegex.FindStringSubmatch(fullText); len(matches) == 4 {
				day, _ := strconv.Atoi(matches[1])
				month, _ := strconv.Atoi(matches[2])
				year, _ := strconv.Atoi(matches[3])
				t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
				profile.BirthDate = &t
				profile.Age = calculateAge(t)
			}
		}

		// Рост
		if strings.Contains(label, "рост") {
			if matches := heightRegex.FindStringSubmatch(fullText); len(matches) > 1 {
				profile.Height, _ = strconv.Atoi(matches[1])
			}
		}

		// Вес
		if strings.Contains(label, "вес") {
			if matches := weightRegex.FindStringSubmatch(fullText); len(matches) > 1 {
				profile.Weight, _ = strconv.Atoi(matches[1])
			}
		}

		// Амплуа
		if strings.Contains(label, "амплуа") {
			profile.Position = extractPositionFromValue(fullText, label)
		}

		// Хват
		if strings.Contains(label, "хват") {
			profile.Handedness = extractHandednessFromValue(fullText, label)
		}

		// Гражданство
		if strings.Contains(label, "гражданство") {
			profile.Citizenship = extractValueAfterLabel(fullText, label)
		}
	})

	// Дополнительно ищем данные в тексте страницы
	if profile.BirthDate == nil {
		if matches := dateRegex.FindStringSubmatch(pageText); len(matches) == 4 {
			day, _ := strconv.Atoi(matches[1])
			month, _ := strconv.Atoi(matches[2])
			year, _ := strconv.Atoi(matches[3])
			t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			profile.BirthDate = &t
			profile.Age = calculateAge(t)
		}
	}

	// Фото
	doc.Find("img.player-photo, img[src*='photo'], img[src*='player']").Each(func(_ int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists && profile.PhotoURL == "" {
			profile.PhotoURL = src
		}
	})

	return profile, nil
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
	textLower := strings.ToLower(text)
	if strings.Contains(textLower, "вратарь") || strings.Contains(text, "В") {
		return "В"
	}
	if strings.Contains(textLower, "защитник") || strings.Contains(text, "З") {
		return "З"
	}
	if strings.Contains(textLower, "нападающий") || strings.Contains(text, "Н") {
		return "Н"
	}
	return ""
}

func extractHandedness(text string) string {
	textLower := strings.ToLower(text)
	if strings.Contains(textLower, "левый") || strings.Contains(text, "Л") {
		return "Л"
	}
	if strings.Contains(textLower, "правый") || strings.Contains(text, "П") {
		return "П"
	}
	return ""
}

func extractCitizenship(text string) string {
	// Убираем "Гражданство:" и возвращаем остаток
	text = strings.TrimSpace(text)
	parts := strings.Split(text, ":")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	if strings.Contains(strings.ToLower(text), "россия") {
		return "Россия"
	}
	return ""
}

// extractValueAfterLabel извлекает значение после метки
// Пример: "гражданство: Россия" -> "Россия"
func extractValueAfterLabel(fullText, label string) string {
	// Ищем позицию метки и берём текст после неё
	idx := strings.Index(strings.ToLower(fullText), label)
	if idx == -1 {
		return ""
	}
	afterLabel := fullText[idx+len(label):]
	// Убираем двоеточие и пробелы
	afterLabel = strings.TrimLeft(afterLabel, ": \t")
	// Берём до конца строки или до следующего тега
	lines := strings.Split(afterLabel, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return strings.TrimSpace(afterLabel)
}

// extractPositionFromValue извлекает позицию из значения
func extractPositionFromValue(fullText, label string) string {
	value := extractValueAfterLabel(fullText, label)
	value = strings.TrimSpace(value)
	// Проверяем сокращения
	if strings.HasPrefix(value, "В") || strings.Contains(strings.ToLower(value), "вратарь") {
		return "В"
	}
	if strings.HasPrefix(value, "З") || strings.Contains(strings.ToLower(value), "защитник") {
		return "З"
	}
	if strings.HasPrefix(value, "Н") || strings.Contains(strings.ToLower(value), "нападающий") {
		return "Н"
	}
	return value
}

// extractHandednessFromValue извлекает хват из значения
func extractHandednessFromValue(fullText, label string) string {
	value := extractValueAfterLabel(fullText, label)
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "Л") || strings.Contains(strings.ToLower(value), "левый") {
		return "Л"
	}
	if strings.HasPrefix(value, "П") || strings.Contains(strings.ToLower(value), "правый") {
		return "П"
	}
	return value
}
