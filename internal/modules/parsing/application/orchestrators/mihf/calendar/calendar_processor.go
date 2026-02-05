package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// processSeason обрабатывает один сезон
func (o *Orchestrator) processSeason(ctx context.Context, season dto.SeasonDTO) (int, int, error) {
	logger.Info(ctx, "[SEASON] Processing", zap.String("year", season.Year))

	// Получаем группы
	groups, err := o.fetchGroups(ctx, season)
	if err != nil {
		return 0, 0, fmt.Errorf("fetch groups: %w", err)
	}

	var totalMatches, totalEvents int

	for _, group := range groups {
		tournaments, err := o.fetchTournaments(ctx, season, group)
		if err != nil {
			logger.Warn(ctx, "Failed to fetch tournaments", zap.Error(err))
			continue
		}

		for _, tournament := range tournaments {
			if tournament.BirthYear < o.config.MinBirthYear() {
				continue
			}
			if max := o.config.MaxBirthYear(); max > 0 && tournament.BirthYear > max {
				continue
			}

			subTournaments, err := o.fetchSubTournaments(ctx, tournament)
			if err != nil {
				logger.Warn(ctx, "Failed to fetch sub-tournaments", zap.Error(err))
				continue
			}

			for _, sub := range subTournaments {
				m, e, err := o.processCalendar(ctx, season, tournament, sub)
				if err != nil {
					logger.Warn(ctx, "Calendar processing failed", zap.Error(err))
					continue
				}
				totalMatches += m
				totalEvents += e
			}
		}
	}

	return totalMatches, totalEvents, nil
}

// processCalendar обрабатывает календарь одного подтурнира
func (o *Orchestrator) processCalendar(ctx context.Context, season dto.SeasonDTO,
	tournament dto.TournamentDTO, sub dto.SubTournamentDTO,
) (int, int, error) {
	calendarURL := sub.URL + "/calendar"
	html, err := o.client.Get(calendarURL)
	if err != nil {
		return 0, 0, fmt.Errorf("get calendar: %w", err)
	}

	matches, err := parsing.ParseCalendar(html)
	if err != nil {
		return 0, 0, fmt.Errorf("parse calendar: %w", err)
	}

	logger.Debug(ctx, "Calendar parsed",
		zap.String("tournament", tournament.Name),
		zap.String("sub", sub.Name),
		zap.Int("matches", len(matches)),
	)

	// Определяем даты турнира
	startDate, endDate := parsing.FindTournamentDates(matches)
	if !startDate.IsZero() {
		tournamentID := fmt.Sprintf("msk:%s-%s-%s", tournament.ID, sub.ID, tournament.GroupID)
		o.updateTournamentDates(ctx, tournamentID, startDate, endDate)
	}

	// Сохраняем матчи
	var savedMatches, savedEvents int
	for _, match := range matches {
		if err := o.saveMatch(ctx, match, season, tournament, sub); err != nil {
			logger.Warn(ctx, "Failed to save match", zap.Error(err))
			continue
		}
		savedMatches++

		// Обновляем города команд из названий (формат: "Команда (Город)")
		if homeCity := extractCityFromTeamName(match.HomeTeamName); homeCity != "" {
			teamID := fmt.Sprintf("msk:%s-%s-%s:%s", tournament.ID, sub.ID, tournament.GroupID, match.HomeTeamID)
			o.teamRepo.UpdateCity(ctx, teamID, homeCity)
		}
		if awayCity := extractCityFromTeamName(match.AwayTeamName); awayCity != "" {
			teamID := fmt.Sprintf("msk:%s-%s-%s:%s", tournament.ID, sub.ID, tournament.GroupID, match.AwayTeamID)
			o.teamRepo.UpdateCity(ctx, teamID, awayCity)
		}

		// Парсим протокол если включено
		if o.config.ParseProtocol() && (match.HomeScore > 0 || match.AwayScore > 0) {
			events, err := o.processProtocol(ctx, match, tournament, sub)
			if err != nil {
				logger.Debug(ctx, "Protocol parsing failed", zap.Error(err))
				continue
			}
			savedEvents += events
		}
	}

	return savedMatches, savedEvents, nil
}

// updateTournamentDates обновляет даты турнира
func (o *Orchestrator) updateTournamentDates(ctx context.Context, id string, start, end time.Time) {
	tournament, err := o.tournamentRepo.GetByID(ctx, id)
	if err != nil || tournament == nil {
		return
	}

	tournament.StartDate = &start
	tournament.EndDate = &end
	o.tournamentRepo.Update(ctx, tournament)
}
