package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SaveTeams сохраняет команды в БД (с дедупликацией)
func (s *orchestratorService) SaveTeams(ctx context.Context, teamsDTO []junior.TeamDTO) ([]*team.Team, error) {
	var saved []*team.Team

	logger.Debug(ctx, fmt.Sprintf("  🔧 SaveTeams: начинаем сохранение %d команд", len(teamsDTO)))

	for i, dto := range teamsDTO {
		t := convertTeamDTO(dto)

		logger.Debug(ctx, fmt.Sprintf("    [%d/%d] Обработка команды: ID=%s, Name=%s, URL=%s",
			i+1, len(teamsDTO), t.ID, t.Name, dto.URL))

		if err := s.teamRepo.Upsert(ctx, t); err != nil {
			logger.Error(ctx, fmt.Sprintf("    ❌ UPSERT FAILED! ID=%s, Name=%s, URL=%s, Error: %v",
				t.ID, t.Name, dto.URL, err))
			return nil, fmt.Errorf("failed to upsert team: %w", err)
		}

		logger.Debug(ctx, fmt.Sprintf("    ✅ Команда upsert: ID=%s, Name=%s", t.ID, t.Name))
		saved = append(saved, t)
	}

	logger.Debug(ctx, fmt.Sprintf("  🔧 SaveTeams: завершено. Сохранено %d команд", len(saved)))
	return saved, nil
}

func convertTeamDTO(dto junior.TeamDTO) *team.Team {
	return &team.Team{
		ID:        team.ExtractIDFromURL(dto.URL),
		URL:       dto.URL,
		Name:      dto.Name,
		City:      dto.City,
		CreatedAt: time.Now(),
	}
}
