package parsing

import (
	"encoding/json"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
)

// ParseSeasons парсит ответ API /api/filter/season
func ParseSeasons(data []byte) ([]dto.SeasonDTO, error) {
	var seasons []dto.SeasonDTO
	if err := json.Unmarshal(data, &seasons); err != nil {
		return nil, err
	}
	return seasons, nil
}

// ParseTournaments парсит ответ API /api/filter/tournament
func ParseTournaments(data []byte) ([]dto.TournamentDTO, error) {
	var tournaments []dto.TournamentDTO
	if err := json.Unmarshal(data, &tournaments); err != nil {
		return nil, err
	}
	return tournaments, nil
}

// ParseTeams парсит ответ API /api/filter/team
func ParseTeams(data []byte) ([]dto.TeamDTO, error) {
	var teams []dto.TeamDTO
	if err := json.Unmarshal(data, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// ParseFilterData парсит ответ API /api/filter/data
// API может вернуть пустой массив [] для турниров без групп
func ParseFilterData(data []byte) (*dto.FilterDataResponse, error) {
	// Проверяем если API вернул пустой массив или массив ошибок
	trimmed := string(data)
	if trimmed == "[]" || (len(trimmed) > 0 && trimmed[0] == '[') {
		// Пустой массив или массив ошибок = турнир без групп
		return &dto.FilterDataResponse{}, nil
	}

	var result dto.FilterDataResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
