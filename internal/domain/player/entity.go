package player

import (
	"strings"
	"time"
)

// Player представляет игрока (Domain Entity)
// db tags сохранены для совместимости с текущими репозиториями
type Player struct {
	ID         string    `db:"id"`
	ProfileURL string    `db:"profile_url"`
	Name       string    `db:"name"`
	BirthDate  time.Time `db:"birth_date"`
	Position   string    `db:"position"`
	Height     *int      `db:"height"`
	Weight     *int      `db:"weight"`
	Handedness *string   `db:"handedness"`

	DataSeason *string `db:"data_season"`

	ExternalID *string `db:"external_id"`
	BirthPlace *string `db:"birth_place"`

	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
		return true
	}
	normalized := strings.TrimSpace(handedness)
	return normalized == HandednessLeft || normalized == HandednessRight
}

// IsNewerSeason проверяет является ли newSeason более новым чем oldSeason
func IsNewerSeason(newSeason, oldSeason string) bool {
	if oldSeason == "" {
		return true
	}
	if newSeason == "" {
		return false
	}
	return newSeason > oldSeason
}
