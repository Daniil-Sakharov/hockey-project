package fhspb

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/parsing"
)

// GetTeamsByTournament получает список команд турнира
func (c *Client) GetTeamsByTournament(tournamentID int) ([]dto.TeamDTO, error) {
	path := fmt.Sprintf("/Teams?TournamentID=%d", tournamentID)

	html, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("get teams page: %w", err)
	}

	return parsing.ParseTeams(html, tournamentID)
}

// GetPlayerURLsFromTeam получает URL игроков из страницы команды
func (c *Client) GetPlayerURLsFromTeam(tournamentID int, teamID string) ([]dto.PlayerURLDTO, error) {
	path := fmt.Sprintf("/Team?TournamentID=%d&TeamID=%s", tournamentID, teamID)

	html, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("get team page: %w", err)
	}

	return parsing.ParsePlayerURLs(html, tournamentID, teamID)
}
