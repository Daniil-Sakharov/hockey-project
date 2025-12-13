package fhspb

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/parsing"
)

// GetAllTournaments получает все турниры со всех сезонов
func (c *Client) GetAllTournaments() ([]dto.TournamentDTO, error) {
	html, err := c.Get("/Tournaments?SeasonID=0")
	if err != nil {
		return nil, fmt.Errorf("get tournaments page: %w", err)
	}

	return parsing.ParseTournaments(html)
}

// GetTournamentsByBirthYear получает турниры с фильтром по году рождения
func (c *Client) GetTournamentsByBirthYear(maxYear int) ([]dto.TournamentDTO, error) {
	tournaments, err := c.GetAllTournaments()
	if err != nil {
		return nil, err
	}

	return parsing.FilterByBirthYear(tournaments, maxYear), nil
}
