package calendar

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	standingsDTO "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/standings"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processStandings(ctx context.Context, tournamentID, externalID string) error {
	// Загружаем страницу турнирной таблицы
	path := fmt.Sprintf("/Standings?TournamentID=%s", externalID)
	html, err := o.client.Get(path)
	if err != nil {
		return fmt.Errorf("fetch standings: %w", err)
	}

	// Парсим
	standingDTOs, err := o.standingsParser.Parse(html)
	if err != nil {
		return fmt.Errorf("parse standings: %w", err)
	}

	logger.Info(ctx, "Parsed standings",
		zap.String("tournament_id", tournamentID),
		zap.Int("count", len(standingDTOs)))

	// Сохраняем
	for _, dto := range standingDTOs {
		standing := o.convertStandingToEntity(dto, tournamentID)

		// Находим ID команды
		teamID := o.findTeamID(ctx, dto.TeamName, tournamentID)
		if teamID == nil {
			logger.Warn(ctx, "Team not found for standing",
				zap.String("team_name", dto.TeamName))
			continue
		}
		standing.TeamID = *teamID

		if err := o.standingRepo.Upsert(ctx, standing); err != nil {
			logger.Error(ctx, "Failed to save standing",
				zap.String("team_name", dto.TeamName),
				zap.Error(err))
			continue
		}
	}

	return nil
}

func (o *Orchestrator) convertStandingToEntity(dto standingsDTO.StandingDTO, tournamentID string) *entities.TeamStanding {
	standing := &entities.TeamStanding{
		TournamentID:   tournamentID,
		Source:         Source,
		Position:       &dto.Position,
		Points:         dto.Points,
		Games:          dto.Games,
		Wins:           dto.Wins,
		WinsOT:         dto.WinsOT,
		WinsSO:         dto.WinsSO,
		LossesSO:       dto.LossesSO,
		LossesOT:       dto.LossesOT,
		Losses:         dto.Losses,
		Draws:          dto.Draws,
		GoalsFor:       dto.GoalsFor,
		GoalsAgainst:   dto.GoalsAgainst,
		GoalDifference: dto.GoalDifference,
	}

	if dto.GroupName != "" {
		standing.GroupName = &dto.GroupName
	}

	if dto.BirthYear > 0 {
		standing.BirthYear = &dto.BirthYear
	}

	return standing
}
