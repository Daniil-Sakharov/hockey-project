package parser

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// SaveTeams сохраняет команды в БД (с дедупликацией)
func (s *orchestratorService) SaveTeams(ctx context.Context, teamsDTO []junior.TeamDTO) ([]*entities.Team, error) {
	var saved []*entities.Team

	for _, dto := range teamsDTO {
		t := convertTeamDTO(dto)

		if err := s.teamRepo.Upsert(ctx, t); err != nil {
			logger.Error(ctx, fmt.Sprintf("    ❌ UPSERT FAILED! ID=%s, Error: %v", t.ID, err))
			return nil, fmt.Errorf("failed to upsert team: %w", err)
		}

		saved = append(saved, t)
	}

	return saved, nil
}

func convertTeamDTO(dto junior.TeamDTO) *entities.Team {
	return &entities.Team{
		ID:        entities.ExtractTeamIDFromURLLegacy(dto.URL),
		URL:       dto.URL,
		Name:      dto.Name,
		City:      dto.City,
		CreatedAt: time.Now(),
	}
}
