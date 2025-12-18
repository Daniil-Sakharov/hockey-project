package profile

import (
	"context"
)

// getPlayerTeamAndRegion получает команду и город игрока используя исправленную логику latest_teams
func (s *Service) getPlayerTeamAndRegion(ctx context.Context, playerID string) (string, string, error) {
	// Получаем связи игрок-команда, отсортированные по started_at (как в CTE latest_teams)
	playerTeams, err := s.playerTeamRepo.GetByPlayer(ctx, playerID)
	if err != nil || len(playerTeams) == 0 {
		return "Неизвестно", "Неизвестно", nil
	}

	// Берем первую команду (она должна быть самой актуальной после сортировки)
	teamID := playerTeams[0].TeamID

	// Получаем информацию о команде
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil || team == nil {
		return "Неизвестно", "Неизвестно", nil
	}

	return team.Name, team.City, nil
}
