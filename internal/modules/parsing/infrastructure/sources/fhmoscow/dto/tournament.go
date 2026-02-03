package dto

import (
	"regexp"
	"strconv"
)

// TournamentDTO представляет турнир из API
type TournamentDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"` // "ПМ 2009 г.р. 25/26"
}

var birthYearRegex = regexp.MustCompile(`(\d{4})\s*г\.?\s*р\.?`)
var seasonRegex = regexp.MustCompile(`(\d{2}/\d{2})$`)

// ParseBirthYear извлекает год рождения из названия турнира
func (t *TournamentDTO) ParseBirthYear() int {
	matches := birthYearRegex.FindStringSubmatch(t.Name)
	if len(matches) > 1 {
		year, _ := strconv.Atoi(matches[1])
		return year
	}
	return 0
}

// ParseSeason извлекает сезон из названия турнира
func (t *TournamentDTO) ParseSeason() string {
	matches := seasonRegex.FindStringSubmatch(t.Name)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
