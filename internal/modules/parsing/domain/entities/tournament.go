package entities

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
)

// Tournament представляет турнир (Domain Entity)
type Tournament struct {
	ID                  string     `db:"id"`
	URL                 string     `db:"url"`
	Name                string     `db:"name"`
	Domain              string     `db:"domain"`
	Season              string     `db:"season"`
	StartDate           *time.Time `db:"start_date"`
	EndDate             *time.Time `db:"end_date"`
	IsEnded             bool       `db:"is_ended"`
	ExternalID          *string    `db:"external_id"`
	BirthYear           *int       `db:"birth_year"`
	GroupName           *string    `db:"group_name"`
	Region              *string    `db:"region"`
	ParentTournamentID  *string    `db:"parent_tournament_id"`
	Source              string     `db:"source"`
	CreatedAt           time.Time  `db:"created_at"`
	LastPlayersParsedAt *time.Time `db:"last_players_parsed_at"`
	LastStatsParsedAt   *time.Time `db:"last_stats_parsed_at"`
	BirthYearGroupsRaw *string    `db:"birth_year_groups"`

	// FallbackBirthYears — годы рождения со страницы списка турниров (не хранится в БД)
	FallbackBirthYears []int `db:"-"`
}

// BirthYearGroups возвращает map[birthYear][]groupName из JSONB поля
func (t *Tournament) BirthYearGroups() map[int][]string {
	if t.BirthYearGroupsRaw == nil || *t.BirthYearGroupsRaw == "" {
		return nil
	}
	var raw map[string][]string
	if err := json.Unmarshal([]byte(*t.BirthYearGroupsRaw), &raw); err != nil {
		return nil
	}
	result := make(map[int][]string, len(raw))
	for k, v := range raw {
		var year int
		if _, err := fmt.Sscanf(k, "%d", &year); err == nil {
			result[year] = v
		}
	}
	return result
}

// ExtractIDFromURL извлекает ID турнира из URL (deprecated, use id.ExtractTournamentIDFromURL)
func ExtractTournamentIDFromURLLegacy(url string) string {
	id, _ := ExtractTournamentIDFromURL(url)
	return id.String()
}

// ConvertJuniorTournament конвертирует DTO турнира junior.fhr.ru в entity
func ConvertJuniorTournament(dto junior.TournamentDTO, domain string) *Tournament {
	var startDate *time.Time
	if dto.StartDate != "" {
		if t, err := time.Parse("02.01.2006", dto.StartDate); err == nil {
			startDate = &t
		}
	}

	var endDate *time.Time
	if dto.EndDate != "" {
		if t, err := time.Parse("02.01.2006", dto.EndDate); err == nil {
			endDate = &t
		}
	}

	return &Tournament{
		ID:                 dto.ID,
		URL:                dto.URL,
		Name:               dto.Name,
		Domain:             domain,
		Season:             dto.Season,
		StartDate:          startDate,
		EndDate:            endDate,
		IsEnded:            dto.IsEnded,
		Source:             SourceJunior,
		CreatedAt:          time.Now(),
		FallbackBirthYears: dto.BirthYears,
	}
}
