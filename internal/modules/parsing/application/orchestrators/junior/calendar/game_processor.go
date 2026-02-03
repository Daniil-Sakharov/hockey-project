package calendar

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/game"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Orchestrator) processGame(ctx context.Context, matchID, externalID string) error {
	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return fmt.Errorf("get match: %w", err)
	}
	if match == nil {
		return fmt.Errorf("match not found: %s", matchID)
	}

	// Строим URL с доменом из матча
	domain := "https://junior.fhr.ru"
	if match.Domain != nil && *match.Domain != "" {
		domain = *match.Domain
	}
	gameURL := fmt.Sprintf("%s/games/%s/", domain, externalID)

	details, err := o.gameParser.Parse(gameURL)
	if err != nil {
		return fmt.Errorf("parse game: %w", err)
	}

	// Обновляем матч
	if err := o.updateMatchDetails(ctx, matchID, details); err != nil {
		return fmt.Errorf("update match: %w", err)
	}

	// Удаляем ВСЕ события матча ОДИН раз (голы, штрафы, вратари)
	if err := o.matchEventRepo.DeleteByMatchID(ctx, matchID); err != nil {
		logger.Warn(ctx, "Failed to delete old events", zap.Error(err))
	}

	// Сохраняем составы ПЕРЕД голами, чтобы искать ассистентов по номеру
	if o.config.ParseLineups() {
		if err := o.saveLineups(ctx, matchID, details); err != nil {
			logger.Warn(ctx, "Failed to save lineups", zap.Error(err))
		}
	}

	// Сохраняем голы
	if err := o.saveGoals(ctx, match, details.Goals); err != nil {
		logger.Warn(ctx, "Failed to save goals", zap.Error(err))
	}

	// Сохраняем штрафы
	if err := o.savePenalties(ctx, match, details.Penalties); err != nil {
		logger.Warn(ctx, "Failed to save penalties", zap.Error(err))
	}

	// Сохраняем события вратарей
	if err := o.saveGoalieEvents(ctx, match, details.GoalieEvents); err != nil {
		logger.Warn(ctx, "Failed to save goalie events", zap.Error(err))
	}

	return o.matchRepo.MarkDetailsParsed(ctx, matchID)
}

func (o *Orchestrator) updateMatchDetails(ctx context.Context, matchID string, d *game.GameDetailsDTO) error {
	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return fmt.Errorf("match not found: %s", matchID)
	}

	match.HomeScore = d.HomeScore
	match.AwayScore = d.AwayScore
	match.HomeScoreP1 = d.HomeScoreP1
	match.AwayScoreP1 = d.AwayScoreP1
	match.HomeScoreP2 = d.HomeScoreP2
	match.AwayScoreP2 = d.AwayScoreP2
	match.HomeScoreP3 = d.HomeScoreP3
	match.AwayScoreP3 = d.AwayScoreP3
	match.HomeScoreOT = d.HomeScoreOT
	match.AwayScoreOT = d.AwayScoreOT
	match.ResultType = strPtr(d.ResultType)
	match.VideoURL = strPtr(d.VideoURL)
	match.Status = entities.MatchStatusFinished

	if d.BirthYear > 0 {
		match.BirthYear = &d.BirthYear
	}
	if d.GroupName != "" {
		match.GroupName = strPtr(d.GroupName)
	}

	return o.matchRepo.Update(ctx, match)
}

func (o *Orchestrator) saveGoals(ctx context.Context, match *entities.Match, goals []game.GoalDTO) error {
	// Сортируем голы хронологически (протокол может идти в обратном порядке)
	sort.Slice(goals, func(i, j int) bool {
		if goals[i].TimeMinutes != goals[j].TimeMinutes {
			return goals[i].TimeMinutes < goals[j].TimeMinutes
		}
		return goals[i].TimeSeconds < goals[j].TimeSeconds
	})

	scoreHome, scoreAway := 0, 0

	for _, g := range goals {
		isHome := g.GoalType == "home"
		if isHome {
			scoreHome++
		} else {
			scoreAway++
		}

		event := &entities.MatchEvent{
			ID:          uuid.New().String(),
			MatchID:     match.ID,
			EventType:   entities.EventTypeGoal,
			Period:      intPtr(g.Period),
			TimeMinutes: intPtr(g.TimeMinutes),
			TimeSeconds: intPtr(g.TimeSeconds),
			GoalType:    strPtr(g.GoalType),
			IsHome:      boolPtr(isHome),
			ScoreHome:   intPtrVal(scoreHome),
			ScoreAway:   intPtrVal(scoreAway),
			Source:      Source,
		}

		// team_id из матча
		if isHome && match.HomeTeamID != nil {
			event.TeamID = match.HomeTeamID
		} else if !isHome && match.AwayTeamID != nil {
			event.TeamID = match.AwayTeamID
		}

		// Находим автора гола
		if g.ScorerURL != "" {
			if p := o.findOrCreatePlayer(ctx, g.ScorerURL, g.ScorerName); p != nil {
				event.ScorerPlayerID = &p.ID
			}
		}

		// Ассистент 1: сначала по URL, затем по номеру в составе
		if g.Assist1URL != "" {
			if p := o.findOrCreatePlayer(ctx, g.Assist1URL, g.Assist1Name); p != nil {
				event.Assist1PlayerID = &p.ID
			}
		}
		if event.Assist1PlayerID == nil && g.Assist1Number > 0 {
			if playerID := o.findPlayerInLineup(ctx, match.ID, g.Assist1Number); playerID != nil {
				event.Assist1PlayerID = playerID
			}
		}

		// Ассистент 2: сначала по URL, затем по номеру в составе
		if g.Assist2URL != "" {
			if p := o.findOrCreatePlayer(ctx, g.Assist2URL, g.Assist2Name); p != nil {
				event.Assist2PlayerID = &p.ID
			}
		}
		if event.Assist2PlayerID == nil && g.Assist2Number > 0 {
			if playerID := o.findPlayerInLineup(ctx, match.ID, g.Assist2Number); playerID != nil {
				event.Assist2PlayerID = playerID
			}
		}

		// Вратарь, пропустивший гол
		if g.GoalieURL != "" {
			if p := o.findOrCreatePlayer(ctx, g.GoalieURL, ""); p != nil {
				event.GoaliePlayerID = &p.ID
			}
		}

		// Игроки на льду
		event.HomePlayersOnIce = o.resolvePlayersToIDs(ctx, match.ID, g.HomePlayersOnIce)
		event.AwayPlayersOnIce = o.resolvePlayersToIDs(ctx, match.ID, g.AwayPlayersOnIce)

		if err := o.matchEventRepo.Create(ctx, event); err != nil {
			logger.Warn(ctx, "Failed to save goal", zap.Error(err))
		}
	}

	return nil
}

func (o *Orchestrator) savePenalties(ctx context.Context, match *entities.Match, penalties []game.PenaltyDTO) error {
	for _, p := range penalties {
		event := &entities.MatchEvent{
			ID:             uuid.New().String(),
			MatchID:        match.ID,
			EventType:      entities.EventTypePenalty,
			Period:         intPtr(p.Period),
			TimeMinutes:    intPtr(p.TimeMinutes),
			TimeSeconds:    intPtr(p.TimeSeconds),
			PenaltyMinutes: intPtr(p.Minutes),
			PenaltyReason:  strPtr(p.Reason),
			IsHome:         boolPtr(p.IsHome),
			Source:         Source,
		}

		// team_id из матча
		if p.IsHome && match.HomeTeamID != nil {
			event.TeamID = match.HomeTeamID
		} else if !p.IsHome && match.AwayTeamID != nil {
			event.TeamID = match.AwayTeamID
		}

		if p.PlayerURL != "" {
			if player := o.findOrCreatePlayer(ctx, p.PlayerURL, p.PlayerName); player != nil {
				event.PenaltyPlayerID = &player.ID
			}
		}

		if err := o.matchEventRepo.Create(ctx, event); err != nil {
			logger.Warn(ctx, "Failed to save penalty", zap.Error(err))
		}
	}

	return nil
}

func (o *Orchestrator) saveGoalieEvents(ctx context.Context, match *entities.Match, events []game.GoalieEventDTO) error {
	for _, ge := range events {
		event := &entities.MatchEvent{
			ID:          uuid.New().String(),
			MatchID:     match.ID,
			EventType:   entities.EventTypeGoalieChange,
			Period:      intPtrVal(game.CalculatePeriod(ge.TimeMinutes)),
			TimeMinutes: intPtrVal(ge.TimeMinutes),
			TimeSeconds: intPtrVal(ge.TimeSeconds),
			IsHome:      boolPtr(ge.IsHome),
			Source:      Source,
		}

		// team_id из матча
		if ge.IsHome && match.HomeTeamID != nil {
			event.TeamID = match.HomeTeamID
		} else if !ge.IsHome && match.AwayTeamID != nil {
			event.TeamID = match.AwayTeamID
		}

		// Вратарь
		if ge.PlayerURL != "" {
			if p := o.findOrCreatePlayer(ctx, ge.PlayerURL, ge.PlayerName); p != nil {
				event.GoaliePlayerID = &p.ID
			}
		}

		if err := o.matchEventRepo.Create(ctx, event); err != nil {
			logger.Warn(ctx, "Failed to save goalie event", zap.Error(err))
		}
	}

	return nil
}

// resolvePlayersToIDs преобразует список URL или текстовых описаний игроков в список ID
func (o *Orchestrator) resolvePlayersToIDs(ctx context.Context, matchID string, items []string) []string {
	if len(items) == 0 {
		return nil
	}

	var ids []string
	for _, item := range items {
		var playerID string

		if strings.HasPrefix(item, "/") || strings.HasPrefix(item, "http") {
			if p := o.findOrCreatePlayer(ctx, item, ""); p != nil {
				playerID = p.ID
			}
		} else {
			jerseyNumber := game.ExtractJerseyNumber(item)
			if jerseyNumber > 0 {
				if pid := o.findPlayerInLineup(ctx, matchID, jerseyNumber); pid != nil {
					playerID = *pid
				}
			}
		}

		if playerID != "" {
			ids = append(ids, playerID)
		}
	}
	return ids
}

func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func intPtrVal(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

// findPlayerInLineup ищет игрока в составе матча по номеру
func (o *Orchestrator) findPlayerInLineup(ctx context.Context, matchID string, jerseyNumber int) *string {
	if jerseyNumber == 0 {
		return nil
	}
	lineup, err := o.matchLineupRepo.GetByMatchAndJersey(ctx, matchID, jerseyNumber)
	if err != nil || lineup == nil {
		return nil
	}
	return &lineup.PlayerID
}
