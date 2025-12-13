package player

import (
	"regexp"
	"strings"
	"time"
)

// Player представляет игрока
type Player struct {
	ID         string    `db:"id"`          // Уникальный ID из URL (924040)
	ProfileURL string    `db:"profile_url"` // /player/...-924040/
	Name       string    `db:"name"`        // ФИО
	BirthDate  time.Time `db:"birth_date"`  // Дата рождения
	Position   string    `db:"position"`    // Защитник/Нападающий/Вратарь
	Height     *int      `db:"height"`      // Рост в см (nullable)
	Weight     *int      `db:"weight"`      // Вес в кг (nullable)
	Handedness *string   `db:"handedness"`  // Левый/Правый (nullable)

	// Данные из registrynew.fhr.ru (будут добавлены позже)
	RegistryID *string `db:"registry_id"` // ID из registrynew.fhr.ru (nullable)
	School     *string `db:"school"`      // Школа подготовки (nullable)
	Rank       *string `db:"rank"`        // Разряд (nullable)

	// Сезон из которого взяты актуальные данные (рост, вес, хват)
	DataSeason *string `db:"data_season"` // Формат: 2024/2025 (nullable)

	// Поля для fhspb.ru
	ExternalID  *string `db:"external_id"`  // Внешний ID из источника (PlayerID UUID для fhspb.ru)
	Citizenship *string `db:"citizenship"`  // Гражданство игрока
	Role        *string `db:"role"`         // Роль в команде: К - капитан, А - ассистент
	BirthPlace  *string `db:"birth_place"`  // Место рождения

	Source    string    `db:"source"`     // junior.fhr.ru, fhspb.ru или both
	CreatedAt time.Time `db:"created_at"` // Дата создания
	UpdatedAt time.Time `db:"updated_at"` // Дата обновления
}

// ExtractIDFromURL извлекает ID из URL профиля
// Пример: /player/abdrashitov-daniyar-rustamovich-2008-05-13-924040/ -> "924040"
func ExtractIDFromURL(profileURL string) string {
	re := regexp.MustCompile(`-(\d+)/?$`)
	matches := re.FindStringSubmatch(profileURL)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// Position constants
const (
	PositionDefender   = "Защитник"
	PositionForward    = "Нападающий"
	PositionGoalkeeper = "Вратарь"
)

// Handedness constants
const (
	HandednessLeft  = "Левый"
	HandednessRight = "Правый"
)

// Source constants
const (
	SourceJunior = "junior.fhr.ru"
	SourceFHSPB  = "fhspb.ru"
	SourceBoth   = "both"
)

// IsValidPosition проверяет валидность позиции
func IsValidPosition(position string) bool {
	normalized := strings.TrimSpace(position)
	return normalized == PositionDefender ||
		normalized == PositionForward ||
		normalized == PositionGoalkeeper
}

// IsValidHandedness проверяет валидность хвата
func IsValidHandedness(handedness string) bool {
	if handedness == "" {
		return true // nullable
	}
	normalized := strings.TrimSpace(handedness)
	return normalized == HandednessLeft || normalized == HandednessRight
}

// IsNewerSeason проверяет является ли newSeason более новым чем oldSeason
// Формат сезонов: "2024/2025"
// Возвращает true если newSeason > oldSeason или oldSeason пустой
func IsNewerSeason(newSeason, oldSeason string) bool {
	if oldSeason == "" {
		return true // Если старый сезон не установлен - всегда обновляем
	}
	if newSeason == "" {
		return false // Если новый сезон пустой - не обновляем
	}
	// Лексикографическое сравнение работает для формата "YYYY/YYYY"
	return newSeason > oldSeason
}
