package calendar

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// processProtocol обрабатывает протокол матча
func (o *Orchestrator) processProtocol(ctx context.Context, match dto.MatchDTO,
	tournament dto.TournamentDTO, sub dto.SubTournamentDTO) (int, error) {

	if match.ProtoURL == "" {
		return 0, nil
	}

	html, err := o.client.Get(match.ProtoURL)
	if err != nil {
		return 0, fmt.Errorf("get protocol: %w", err)
	}

	proto, err := parsing.ParseMatchProtocol(html)
	if err != nil {
		return 0, fmt.Errorf("parse protocol: %w", err)
	}

	matchID := fmt.Sprintf("msk:%s", match.ExternalID)
	tournamentID := fmt.Sprintf("msk:%s-%s-%s", tournament.ID, sub.ID, tournament.GroupID)
	// Team ID format: msk:{tournament_id}-{sub_id}-{group_id}:{team_id}
	homeTeamID := fmt.Sprintf("msk:%s-%s-%s:%s", tournament.ID, sub.ID, tournament.GroupID, match.HomeTeamID)
	awayTeamID := fmt.Sprintf("msk:%s-%s-%s:%s", tournament.ID, sub.ID, tournament.GroupID, match.AwayTeamID)

	// Обновляем логотипы команд
	o.updateTeamLogos(ctx, homeTeamID, awayTeamID, proto)

	// Сохраняем события матча
	eventsCount, err := o.saveMatchEvents(ctx, matchID, proto)
	if err != nil {
		logger.Warn(ctx, "Failed to save match events", zap.Error(err))
	}

	// Сохраняем составы (с автосозданием недостающих игроков)
	if err := o.saveMatchLineups(ctx, matchID, homeTeamID, awayTeamID, tournamentID, tournament.BirthYear, proto); err != nil {
		logger.Warn(ctx, "Failed to save lineups", zap.Error(err))
	}

	// Обновляем счет по периодам
	o.updateMatchScoreByPeriods(ctx, matchID, proto)

	// Помечаем матч как распарсенный
	o.matchRepo.MarkDetailsParsed(ctx, matchID)

	logger.Debug(ctx, "Protocol processed",
		zap.String("match_id", matchID),
		zap.Int("events", eventsCount),
		zap.Int("home_lineup", len(proto.HomeLineup)),
		zap.Int("away_lineup", len(proto.AwayLineup)),
	)

	return eventsCount, nil
}

// updateTeamLogos обновляет логотипы команд
func (o *Orchestrator) updateTeamLogos(ctx context.Context, homeID, awayID string, proto *dto.MatchProtocolDTO) {
	if proto.HomeLogoURL != "" {
		o.teamRepo.UpdateLogoURL(ctx, homeID, proto.HomeLogoURL)
	}
	if proto.AwayLogoURL != "" {
		o.teamRepo.UpdateLogoURL(ctx, awayID, proto.AwayLogoURL)
	}
}

// updateMatchScoreByPeriods вычисляет и обновляет счет матча по периодам из голов
func (o *Orchestrator) updateMatchScoreByPeriods(ctx context.Context, matchID string, proto *dto.MatchProtocolDTO) {
	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil || match == nil {
		return
	}

	// Считаем счёт по периодам из голов (не берём с сайта - там ошибки)
	var homeP1, awayP1, homeP2, awayP2, homeP3, awayP3, homeOT, awayOT int

	for _, goal := range proto.Goals {
		period := calculatePeriodFromTime(goal.TimeMinutes)
		switch period {
		case 1:
			if goal.IsHome {
				homeP1++
			} else {
				awayP1++
			}
		case 2:
			if goal.IsHome {
				homeP2++
			} else {
				awayP2++
			}
		case 3:
			if goal.IsHome {
				homeP3++
			} else {
				awayP3++
			}
		case 4: // OT
			if goal.IsHome {
				homeOT++
			} else {
				awayOT++
			}
		}
	}

	match.HomeScoreP1 = scorePtr(homeP1)
	match.AwayScoreP1 = scorePtr(awayP1)
	match.HomeScoreP2 = scorePtr(homeP2)
	match.AwayScoreP2 = scorePtr(awayP2)
	match.HomeScoreP3 = scorePtr(homeP3)
	match.AwayScoreP3 = scorePtr(awayP3)

	// OT только если были голы после 60 минуты
	if homeOT > 0 || awayOT > 0 {
		match.HomeScoreOT = scorePtr(homeOT)
		match.AwayScoreOT = scorePtr(awayOT)
	}

	o.matchRepo.Update(ctx, match)
}

// calculatePeriodFromTime вычисляет период по минутам
// 0-19 = П1, 20-39 = П2, 40-59 = П3, 60+ = OT
func calculatePeriodFromTime(minutes int) int {
	switch {
	case minutes < 20:
		return 1
	case minutes < 40:
		return 2
	case minutes < 60:
		return 3
	default:
		return 4
	}
}
