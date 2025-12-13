package fhspb

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/parsing"
)

// GetPlayer получает данные игрока
func (c *Client) GetPlayer(tournamentID int, teamID, playerID string) (*dto.PlayerDTO, error) {
	path := fmt.Sprintf("/Player?TournamentID=%d&TeamID=%s&PlayerID=%s", tournamentID, teamID, playerID)

	html, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("get player page: %w", err)
	}

	return parsing.ParsePlayer(html, playerID)
}
