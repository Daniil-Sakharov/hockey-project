package stats

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// parseTournamentStats парсит статистику одного турнира
func (s *service) parseTournamentStats(
	ctx context.Context,
	t *entities.Tournament,
) (int, error) {
	// Валидация данных турнира
	if t.Domain == "" {
		return 0, fmt.Errorf("tournament domain is empty")
	}
	if t.URL == "" {
		return 0, fmt.Errorf("tournament URL is empty")
	}
	if t.ID == "" {
		return 0, fmt.Errorf("tournament ID is empty")
	}

	// Парсим статистику через StatsParserService
	count, err := s.statsService.ParseTournamentStats(
		ctx,
		t.Domain,
		t.URL,
		t.ID,
		t.Season,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to parse tournament stats: %w", err)
	}

	return count, nil
}
