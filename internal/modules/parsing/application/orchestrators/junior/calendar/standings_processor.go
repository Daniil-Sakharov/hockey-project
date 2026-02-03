package calendar

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/standings"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Orchestrator) processStandings(ctx context.Context, tournamentID, tournamentURL string) error {
	standingsList, err := o.standingsParser.Parse(tournamentURL)
	if err != nil {
		return err
	}

	logger.Debug(ctx, "Parsed standings",
		zap.Int("count", len(standingsList)))

	for _, s := range standingsList {
		standing := o.convertStanding(ctx, tournamentID, s)
		if standing == nil {
			continue
		}

		if err := o.standingRepo.Upsert(ctx, standing); err != nil {
			logger.Error(ctx, "Failed to save standing",
				zap.String("team", s.TeamName),
				zap.Error(err))
		}
	}

	return nil
}

func (o *Orchestrator) convertStanding(ctx context.Context, tournamentID string, s standings.StandingDTO) *entities.TeamStanding {
	// Находим команду по ID из URL
	teamID := entities.ExtractTeamIDFromURLLegacy(s.TeamURL)
	if teamID == "" {
		logger.Debug(ctx, "Team not found for standing",
			zap.String("team_name", s.TeamName),
			zap.String("team_url", s.TeamURL))
		return nil
	}

	// Проверяем существование команды
	team, err := o.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		logger.Error(ctx, "Error getting team by ID",
			zap.String("team_id", teamID),
			zap.Error(err))
		return nil
	}
	if team == nil {
		logger.Debug(ctx, "Team not found by ID",
			zap.String("team_id", teamID),
			zap.String("team_name", s.TeamName))
		return nil
	}
	logger.Debug(ctx, "Team FOUND by ID",
		zap.String("team_id", teamID),
		zap.String("team_name", team.Name))

	standing := &entities.TeamStanding{
		ID:             uuid.New().String(),
		TournamentID:   tournamentID,
		TeamID:         teamID,
		Position:       intPtr(s.Position),
		Points:         s.Points,
		Games:          s.Games,
		Wins:           s.Wins,
		WinsOT:         s.WinsOT,
		WinsSO:         s.WinsSO,
		LossesSO:       s.LossesSO,
		LossesOT:       s.LossesOT,
		Losses:         s.Losses,
		Draws:          s.Draws,
		GoalsFor:       s.GoalsFor,
		GoalsAgainst:   s.GoalsAgainst,
		GoalDifference: s.GoalDifference,
		GroupName:      strPtr(s.GroupName),
		Source:         Source,
	}

	if s.BirthYear > 0 {
		standing.BirthYear = &s.BirthYear
	}

	return standing
}
