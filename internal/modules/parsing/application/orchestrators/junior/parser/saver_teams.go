package parser

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

var teamBirthYearRegex = regexp.MustCompile(`\b(20\d{2})\b`)

const minTeamBirthYear = 2008

// SaveTeams сохраняет команды в БД (с дедупликацией и привязкой к турниру)
func (s *orchestratorService) SaveTeams(ctx context.Context, teamsDTO []junior.TeamDTO, tournamentID string) ([]*entities.Team, error) {
	var saved []*entities.Team

	for _, dto := range teamsDTO {
		// Фильтрация по году рождения команды
		if birthYear := extractTeamBirthYear(dto.Name); birthYear > 0 && birthYear < minTeamBirthYear {
			logger.Info(ctx, fmt.Sprintf("    ⏭️  Skipping team %s (birth_year %d < %d)", dto.Name, birthYear, minTeamBirthYear))
			continue
		}

		t := convertTeamDTO(dto, tournamentID)

		if err := s.teamRepo.Upsert(ctx, t); err != nil {
			logger.Error(ctx, fmt.Sprintf("    ❌ UPSERT FAILED! ID=%s, Error: %v", t.ID, err))
			return nil, fmt.Errorf("failed to upsert team: %w", err)
		}

		saved = append(saved, t)
	}

	return saved, nil
}

// extractTeamBirthYear извлекает год рождения из названия команды (например "ЦСКА 2007" -> 2007)
func extractTeamBirthYear(name string) int {
	matches := teamBirthYearRegex.FindStringSubmatch(name)
	if len(matches) > 1 {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			return year
		}
	}
	return 0
}

func convertTeamDTO(dto junior.TeamDTO, tournamentID string) *entities.Team {
	t := &entities.Team{
		ID:        entities.ExtractTeamIDFromURLLegacy(dto.URL),
		URL:       dto.URL,
		Name:      dto.Name,
		City:      dto.City,
		Source:    entities.SourceJunior,
		CreatedAt: time.Now(),
	}
	if dto.LogoURL != "" {
		t.LogoURL = &dto.LogoURL
	}
	if tournamentID != "" {
		t.TournamentID = &tournamentID
	}
	return t
}
