package player

import (
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// ParseProfile парсит профиль игрока для получения дополнительных данных
func (p *Parser) ParseProfile(domain, profileURL string) (*types.PlayerProfileDTO, error) {
	fullURL := domain + profileURL

	resp, err := p.http.MakeRequest(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки профиля: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	profile := &types.PlayerProfileDTO{}

	// Парсим фото игрока
	doc.Find("img[alt='photo']").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			if strings.Contains(src, "player_photo") {
				profile.PhotoURL = buildFullPhotoURL(domain, src)
			}
		}
	})

	// Парсим данные из пар <strong>Label:</strong><span>Value</span>
	doc.Find("strong").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Text())
		// Получаем следующий элемент span
		value := strings.TrimSpace(s.NextFiltered("span").Text())
		if value == "" {
			// Иногда значение в следующем текстовом узле
			value = strings.TrimSpace(s.Parent().Text())
			value = strings.TrimPrefix(value, label)
			value = strings.TrimSpace(value)
		}

		switch {
		case strings.Contains(label, "Амплуа"):
			profile.Position = normalizePosition(value)
		case strings.Contains(label, "Рост"):
			profile.Height = extractNumber(value)
		case strings.Contains(label, "Вес"):
			profile.Weight = extractNumber(value)
		case strings.Contains(label, "Хват"):
			profile.Handedness = normalizeHandedness(value)
		case strings.Contains(label, "Гражданство"):
			profile.Citizenship = normalizeCitizenship(value)
		}
	})

	return profile, nil
}

// normalizePosition нормализует позицию
func normalizePosition(pos string) string {
	pos = strings.TrimSpace(pos)
	lower := strings.ToLower(pos)
	switch {
	case strings.Contains(lower, "защитник"):
		return "Защитник"
	case strings.Contains(lower, "нападающий"):
		return "Нападающий"
	case strings.Contains(lower, "вратарь"):
		return "Вратарь"
	}
	return pos
}

// normalizeHandedness нормализует хват
func normalizeHandedness(hand string) string {
	hand = strings.TrimSpace(hand)
	lower := strings.ToLower(hand)
	switch {
	case strings.Contains(lower, "лев"):
		return "Левый"
	case strings.Contains(lower, "прав"):
		return "Правый"
	}
	return hand
}

// normalizeCitizenship нормализует гражданство
func normalizeCitizenship(cit string) string {
	cit = strings.TrimSpace(cit)
	// Приводим к нормальному виду (первая буква заглавная)
	if cit == "" {
		return ""
	}
	lower := strings.ToLower(cit)
	if lower == "россия" || lower == "рф" || lower == "russia" {
		return "Россия"
	}
	// Capitalize first letter
	return strings.ToUpper(string(cit[0])) + strings.ToLower(cit[1:])
}

// buildFullPhotoURL формирует полный URL фото игрока.
// Добавляет домен к относительному пути и заменяет миниатюру 50x50 на 200x200.
// Возвращает пустую строку для заглушек (placeholder).
func buildFullPhotoURL(domain, src string) string {
	if !strings.Contains(src, "player_photo") {
		return ""
	}
	src = strings.Replace(src, "50_50-70.webp", "200_200-70.webp", 1)
	if strings.HasPrefix(src, "/") {
		return domain + src
	}
	return src
}

// extractNumber извлекает числовое значение из строки
func extractNumber(s string) string {
	s = strings.TrimSpace(s)
	var result strings.Builder
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result.WriteRune(c)
		}
	}
	return result.String()
}
