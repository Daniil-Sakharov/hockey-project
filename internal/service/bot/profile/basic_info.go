package profile

import (
	"context"
)

// getPlayerTeamAndRegion получает команду и город игрока
func (s *Service) getPlayerTeamAndRegion(ctx context.Context, playerID string) (string, string, error) {
	// Получаем связи игрок-команда
	playerTeams, err := s.playerTeamRepo.GetByPlayer(ctx, playerID)
	if err != nil || len(playerTeams) == 0 {
		return "Неизвестно", "Неизвестно", nil
	}

	// Берем первую команду (обычно это текущая или последняя)
	// TODO: в будущем можно добавить логику выбора активной команды
	teamID := playerTeams[0].TeamID

	// Получаем информацию о команде
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil || team == nil {
		return "Неизвестно", "Неизвестно", nil
	}

	return team.Name, team.City, nil
}
