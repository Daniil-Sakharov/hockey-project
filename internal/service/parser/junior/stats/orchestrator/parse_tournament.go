package stats_orchestrator

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// parseTournamentStats парсит статистику одного турнира
// Возвращает количество обработанных записей
func (s *service) parseTournamentStats(
	ctx context.Context,
	t *tournament.Tournament,
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
	)
	if err != nil {
		return 0, fmt.Errorf("failed to parse tournament stats: %w", err)
	}

	return count, nil
}
