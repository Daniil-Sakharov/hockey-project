package calendar

import (
	"context"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/game"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Orchestrator) saveLineups(ctx context.Context, matchID string, details *game.GameDetailsDTO) error {
	// Удаляем старые составы
	if err := o.matchLineupRepo.DeleteByMatchID(ctx, matchID); err != nil {
		return err
	}

	match, err := o.matchRepo.GetByID(ctx, matchID)
	if err != nil || match == nil {
		return err
	}

	// Сохраняем домашнюю команду
	if match.HomeTeamID != nil {
		for _, p := range details.HomeLineup {
			if err := o.savePlayerLineup(ctx, matchID, *match.HomeTeamID, p); err != nil {
				logger.Warn(ctx, "Failed to save home lineup", zap.Error(err))
			}
		}
	}

	// Сохраняем гостевую команду
	if match.AwayTeamID != nil {
		for _, p := range details.AwayLineup {
			if err := o.savePlayerLineup(ctx, matchID, *match.AwayTeamID, p); err != nil {
				logger.Warn(ctx, "Failed to save away lineup", zap.Error(err))
			}
		}
	}

	return nil
}

func (o *Orchestrator) savePlayerLineup(ctx context.Context, matchID, teamID string, p game.PlayerLineup) error {
	player := o.findOrCreatePlayer(ctx, p.PlayerURL, p.PlayerName)
	if player == nil {
		return nil
	}

	lineup := &entities.MatchLineup{
		ID:             uuid.New().String(),
		MatchID:        matchID,
		PlayerID:       player.ID,
		TeamID:         teamID,
		JerseyNumber:   intPtr(p.JerseyNumber),
		Position:       strPtr(p.Position),
		CaptainRole:    getCaptainRole(p.Role),
		Goals:          p.Goals,
		Assists:        p.Assists,
		PenaltyMinutes: p.PenaltyMinutes,
		PlusMinus:      p.PlusMinus,
		Saves:          p.Saves,
		GoalsAgainst:   p.GoalsAgainst,
		TimeOnIce:      p.TimeOnIce,
		Source:         Source,
	}

	return o.matchLineupRepo.Create(ctx, lineup)
}

func (o *Orchestrator) findOrCreatePlayer(ctx context.Context, profileURL, name string) *entities.Player {
	if profileURL == "" {
		return nil
	}

	// Проверяем, что это действительно URL, а не имя игрока
	// URL должен начинаться с / или http
	if !strings.HasPrefix(profileURL, "/") && !strings.HasPrefix(profileURL, "http") {
		return nil
	}

	// Пробуем найти по URL
	player, err := o.playerRepo.GetByProfileURL(ctx, profileURL)
	if err == nil && player != nil {
		return player
	}

	// Если profileParser не задан, создаём минимальную запись
	if o.profileParser == nil {
		newPlayer := &entities.Player{
			ID:         uuid.New().String(),
			ProfileURL: profileURL,
			Name:       name,
			Position:   "Нападающий", // По умолчанию
			Source:     Source,
		}

		if err := o.playerRepo.Create(ctx, newPlayer); err != nil {
			logger.Warn(ctx, "Failed to create player",
				zap.String("name", name),
				zap.Error(err))
			return nil
		}

		return newPlayer
	}

	// Парсим профиль и создаём
	info, err := o.profileParser.ParseProfile(ctx, profileURL)
	if err != nil {
		logger.Debug(ctx, "Failed to parse player profile",
			zap.String("url", profileURL),
			zap.Error(err))
		return nil
	}

	birthDate := parseBirthDate(info.BirthDate)

	newPlayer := &entities.Player{
		ID:          uuid.New().String(),
		ProfileURL:  profileURL,
		Name:        name,
		BirthDate:   birthDate,
		Position:    info.Position,
		Citizenship: strPtr(info.Citizenship),
		Source:      Source,
	}

	if err := o.playerRepo.Create(ctx, newPlayer); err != nil {
		logger.Warn(ctx, "Failed to create player",
			zap.String("name", name),
			zap.Error(err))
		return nil
	}

	return newPlayer
}

// getCaptainRole преобразует роль игрока в формат для БД
func getCaptainRole(role string) *string {
	if role == "" {
		return nil
	}
	// Роли: C (капитан), A (ассистент)
	// В БД храним: К (капитан), А (ассистент)
	switch role {
	case "C":
		r := "К"
		return &r
	case "A":
		r := "А"
		return &r
	default:
		return nil
	}
}

func parseBirthDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	formats := []string{
		"02.01.2006",
		"2006-01-02",
		"02/01/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}
