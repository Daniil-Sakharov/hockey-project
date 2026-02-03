package calendar

import (
	"context"
	"regexp"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

var playerIDRegex = regexp.MustCompile(`PlayerID=(\d+)`)

// findOrCreatePlayer находит или создаёт игрока
func (o *Orchestrator) findOrCreatePlayer(ctx context.Context, playerURL, playerName string) *string {
	if playerURL == "" && playerName == "" {
		return nil
	}

	// Извлекаем PlayerID из URL
	externalID := extractPlayerID(playerURL)
	if externalID == "" {
		return nil
	}

	// Пробуем найти в БД
	player, err := o.playerRepo.GetByExternalID(ctx, externalID)
	if err == nil && player != nil {
		return &player.ID
	}

	// Если не найден - создаём минимальную запись
	newPlayer := &fhspb.Player{
		ExternalID: externalID,
		FullName:   cleanPlayerName(playerName),
	}

	if playerURL != "" {
		newPlayer.ProfileURL = &playerURL
	}

	playerID, err := o.playerRepo.Upsert(ctx, newPlayer)
	if err != nil {
		logger.Error(ctx, "Failed to create player",
			zap.String("name", playerName),
			zap.Error(err))
		return nil
	}

	return &playerID
}

func extractPlayerID(url string) string {
	matches := playerIDRegex.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func cleanPlayerName(name string) string {
	// Убираем номер из начала имени если есть
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`^\d+\s+`).ReplaceAllString(name, "")
	return name
}
