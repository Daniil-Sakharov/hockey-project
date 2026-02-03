package calendar

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	matchDTO "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/match"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processGame(ctx context.Context, matchID, externalID, tournamentExternalID string) error {
	// Загружаем страницу протокола
	path := fmt.Sprintf("/Match?TournamentID=%s&MatchID=%s", tournamentExternalID, externalID)
	html, err := o.client.Get(path)
	if err != nil {
		return fmt.Errorf("fetch match: %w", err)
	}

	// Парсим
	details, err := o.matchParser.Parse(html)
	if err != nil {
		return fmt.Errorf("parse match: %w", err)
	}

	// Обновляем счёт по периодам в матче
	if err := o.updateMatchPeriodScores(ctx, matchID, details); err != nil {
		logger.Warn(ctx, "Failed to update period scores", zap.Error(err))
	}

	// Сохраняем голы
	if err := o.saveGoals(ctx, matchID, details.Goals); err != nil {
		logger.Warn(ctx, "Failed to save goals", zap.Error(err))
	}

	// Сохраняем штрафы
	if err := o.savePenalties(ctx, matchID, details.Penalties); err != nil {
		logger.Warn(ctx, "Failed to save penalties", zap.Error(err))
	}

	// Сохраняем составы
	if o.config.ParseLineups() {
		if err := o.saveLineups(ctx, matchID, details); err != nil {
			logger.Warn(ctx, "Failed to save lineups", zap.Error(err))
		}
	}

	// Сохраняем статистику бросков
	if err := o.saveTeamStats(ctx, matchID, details); err != nil {
		logger.Warn(ctx, "Failed to save team stats", zap.Error(err))
	}

	// Помечаем матч как обработанный
	if err := o.matchRepo.MarkDetailsParsed(ctx, matchID); err != nil {
		return fmt.Errorf("mark details parsed: %w", err)
	}

	logger.Info(ctx, "Match details saved",
		zap.String("match_id", matchID),
		zap.Int("goals", len(details.Goals)),
		zap.Int("penalties", len(details.Penalties)))

	return nil
}

func (o *Orchestrator) updateMatchPeriodScores(ctx context.Context, matchID string, details *matchDTO.MatchDetailsDTO) error {
	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	match.HomeScoreP1 = &details.HomeScoreP1
	match.AwayScoreP1 = &details.AwayScoreP1
	match.HomeScoreP2 = &details.HomeScoreP2
	match.AwayScoreP2 = &details.AwayScoreP2
	match.HomeScoreP3 = &details.HomeScoreP3
	match.AwayScoreP3 = &details.AwayScoreP3

	if details.HomeScoreOT > 0 || details.AwayScoreOT > 0 {
		match.HomeScoreOT = &details.HomeScoreOT
		match.AwayScoreOT = &details.AwayScoreOT
	}

	return o.matchRepo.Update(ctx, match)
}

func (o *Orchestrator) saveGoals(ctx context.Context, matchID string, goals []matchDTO.GoalDTO) error {
	for _, g := range goals {
		event := &entities.MatchEvent{
			MatchID:     matchID,
			EventType:   "goal",
			Period:      &g.Period,
			TimeMinutes: &g.TimeMinutes,
			TimeSeconds: &g.TimeSeconds,
			Source:      Source,
		}

		// Находим игроков
		if g.ScorerURL != "" {
			scorerID := o.findOrCreatePlayer(ctx, g.ScorerURL, g.ScorerName)
			event.ScorerPlayerID = scorerID
		}

		if g.Assist1URL != "" {
			assist1ID := o.findOrCreatePlayer(ctx, g.Assist1URL, g.Assist1Name)
			event.Assist1PlayerID = assist1ID
		}

		if g.Assist2URL != "" {
			assist2ID := o.findOrCreatePlayer(ctx, g.Assist2URL, g.Assist2Name)
			event.Assist2PlayerID = assist2ID
		}

		// Добавляем счёт после гола
		event.ScoreHome = &g.ScoreHome
		event.ScoreAway = &g.ScoreAway

		// Добавляем тип гола (PP1, PP2, SH1, SH2, EN, PS, GWG)
		if g.GoalType != "" {
			event.GoalType = &g.GoalType
		}

		// Добавляем флаг домашней команды
		event.IsHome = &g.IsHome

		if err := o.matchEventRepo.Create(ctx, event); err != nil {
			logger.Error(ctx, "Failed to save goal", zap.Error(err))
		}
	}

	return nil
}

func (o *Orchestrator) savePenalties(ctx context.Context, matchID string, penalties []matchDTO.PenaltyDTO) error {
	for _, p := range penalties {
		event := &entities.MatchEvent{
			MatchID:        matchID,
			EventType:      "penalty",
			Period:         &p.Period,
			TimeMinutes:    &p.TimeMinutes,
			TimeSeconds:    &p.TimeSeconds,
			PenaltyMinutes: &p.Minutes,
			PenaltyReason:  &p.Reason,
			Source:         Source,
		}

		if p.PlayerURL != "" {
			playerID := o.findOrCreatePlayer(ctx, p.PlayerURL, p.PlayerName)
			event.PenaltyPlayerID = playerID
		}

		// Добавляем код нарушения (ПОДН, ГРУБ и т.д.)
		if p.ReasonCode != "" {
			event.PenaltyReasonCode = &p.ReasonCode
		}

		// Добавляем флаг домашней команды
		event.IsHome = &p.IsHome

		if err := o.matchEventRepo.Create(ctx, event); err != nil {
			logger.Error(ctx, "Failed to save penalty", zap.Error(err))
		}
	}

	return nil
}

func (o *Orchestrator) saveLineups(ctx context.Context, matchID string, details *matchDTO.MatchDetailsDTO) error {
	// Получаем матч для ID команд
	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	// Домашняя команда
	for _, p := range details.HomeLineup {
		lineup := o.convertLineupToEntity(p, matchID, match.HomeTeamID)
		if err := o.matchLineupRepo.Upsert(ctx, lineup); err != nil {
			logger.Error(ctx, "Failed to save lineup", zap.Error(err))
		}
	}

	// Гостевая команда
	for _, p := range details.AwayLineup {
		lineup := o.convertLineupToEntity(p, matchID, match.AwayTeamID)
		if err := o.matchLineupRepo.Upsert(ctx, lineup); err != nil {
			logger.Error(ctx, "Failed to save lineup", zap.Error(err))
		}
	}

	return nil
}

func (o *Orchestrator) convertLineupToEntity(dto matchDTO.PlayerLineupDTO, matchID string, teamID *string) *entities.MatchLineup {
	lineup := &entities.MatchLineup{
		MatchID:        matchID,
		JerseyNumber:   &dto.Number,
		Position:       &dto.Position,
		Goals:          dto.Goals,
		Assists:        dto.Assists,
		PenaltyMinutes: dto.PenaltyMinutes,
		PlusMinus:      dto.PlusMinus,
		Source:         Source,
	}

	if teamID != nil {
		lineup.TeamID = *teamID
	}

	if dto.CaptainRole != "" {
		lineup.CaptainRole = &dto.CaptainRole
	}

	return lineup
}

func (o *Orchestrator) saveTeamStats(ctx context.Context, matchID string, details *matchDTO.MatchDetailsDTO) error {
	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	// Домашняя команда
	if match.HomeTeamID != nil {
		homeStats := &entities.MatchTeamStats{
			MatchID:    matchID,
			TeamID:     *match.HomeTeamID,
			ShotsP1:    details.HomeShots.P1,
			ShotsP2:    details.HomeShots.P2,
			ShotsP3:    details.HomeShots.P3,
			ShotsOT:    details.HomeShots.OT,
			ShotsTotal: details.HomeShots.Total,
			Source:     Source,
		}
		homeStats.CalculateTotal()
		if err := o.matchTeamStatsRepo.Upsert(ctx, homeStats); err != nil {
			logger.Error(ctx, "Failed to save home team stats", zap.Error(err))
		}
	}

	// Гостевая команда
	if match.AwayTeamID != nil {
		awayStats := &entities.MatchTeamStats{
			MatchID:    matchID,
			TeamID:     *match.AwayTeamID,
			ShotsP1:    details.AwayShots.P1,
			ShotsP2:    details.AwayShots.P2,
			ShotsP3:    details.AwayShots.P3,
			ShotsOT:    details.AwayShots.OT,
			ShotsTotal: details.AwayShots.Total,
			Source:     Source,
		}
		awayStats.CalculateTotal()
		if err := o.matchTeamStatsRepo.Upsert(ctx, awayStats); err != nil {
			logger.Error(ctx, "Failed to save away team stats", zap.Error(err))
		}
	}

	return nil
}
