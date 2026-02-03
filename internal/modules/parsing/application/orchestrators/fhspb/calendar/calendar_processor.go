package calendar

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	calendarDTO "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/calendar"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// matchIDPrefix - префикс для ID матчей FHSPB
const matchIDPrefix = "spb"

func (o *Orchestrator) processCalendar(ctx context.Context, tournamentID, externalID string) error {
	// Загружаем страницу календаря
	path := fmt.Sprintf("/Schedule?TournamentID=%s", externalID)
	html, err := o.client.Get(path)
	if err != nil {
		return fmt.Errorf("fetch calendar: %w", err)
	}

	// Парсим
	extID, _ := strconv.Atoi(externalID)
	matchDTOs, err := o.calendarParser.Parse(html, extID)
	if err != nil {
		return fmt.Errorf("parse calendar: %w", err)
	}

	logger.Info(ctx, "Parsed calendar matches",
		zap.String("tournament_id", tournamentID),
		zap.Int("count", len(matchDTOs)))

	// Сохраняем матчи
	for _, dto := range matchDTOs {
		match := o.convertCalendarMatchToEntity(dto, tournamentID)

		// Пропускаем существующие если нужно
		if o.config.SkipExisting() {
			existing, err := o.matchRepo.GetByExternalID(ctx, match.ExternalID, Source)
			if err == nil && existing != nil {
				continue
			}
		}

		// Находим ID команд
		match.HomeTeamID = o.findTeamID(ctx, dto.HomeTeamName, tournamentID)
		match.AwayTeamID = o.findTeamID(ctx, dto.AwayTeamName, tournamentID)

		if err := o.matchRepo.Upsert(ctx, match); err != nil {
			logger.Error(ctx, "Failed to save match",
				zap.String("external_id", dto.ExternalID),
				zap.Error(err))
			continue
		}
	}

	return nil
}

func (o *Orchestrator) convertCalendarMatchToEntity(dto calendarDTO.MatchDTO, tournamentID string) *entities.Match {
	match := &entities.Match{
		ID:           fmt.Sprintf("%s:%s", matchIDPrefix, dto.ExternalID),
		ExternalID:   dto.ExternalID,
		TournamentID: &tournamentID,
		Source:       Source,
		Status:       entities.MatchStatusScheduled,
	}

	if dto.MatchNumber > 0 {
		match.MatchNumber = &dto.MatchNumber
	}

	if dto.ScheduledAt != nil {
		match.ScheduledAt = dto.ScheduledAt
	}

	if dto.Venue != "" {
		match.Venue = &dto.Venue
	}

	if dto.IsFinished {
		match.Status = entities.MatchStatusFinished
		match.HomeScore = dto.HomeScore
		match.AwayScore = dto.AwayScore

		if dto.ResultType != "" {
			match.ResultType = &dto.ResultType
		}
	}

	return match
}

func (o *Orchestrator) findTeamID(ctx context.Context, teamName, tournamentID string) *string {
	if teamName == "" {
		return nil
	}

	// Ищем команду по названию в турнире
	team, err := o.teamRepo.GetByName(ctx, teamName, tournamentID)
	if err != nil || team == nil {
		return nil
	}

	return &team.ID
}
