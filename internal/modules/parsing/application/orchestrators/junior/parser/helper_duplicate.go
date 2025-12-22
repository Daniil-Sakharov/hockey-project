package parser

import (
	"context"
	"fmt"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// isDuplicateDomain проверяет является ли домен копией другого домена
// Проверка выполняется через парсинг ПЕРВОГО сезона (самого нового)
// Если ВСЕ турниры первого сезона - дубликаты, домен пропускается
func (s *orchestratorService) isDuplicateDomain(
	ctx context.Context,
	domain string,
	globalDedup *sync.Map,
) (bool, error) {
	logger.Info(ctx, "    🔍 Quick check: parsing latest season...")

	// 1. Извлекаем все сезоны домена
	seasons, err := s.juniorService.ExtractAllSeasons(ctx, domain)
	if err != nil {
		return false, fmt.Errorf("failed to extract seasons: %w", err)
	}

	if len(seasons) == 0 {
		logger.Warn(ctx, "    ⚠️  No seasons found - skipping domain")
		return false, nil // Не дубликат, просто пустой домен
	}

	// 2. Берем ПЕРВЫЙ сезон (самый новый, например 2025/2026)
	firstSeason := seasons[0]
	logger.Info(ctx, fmt.Sprintf("    📅 Checking season: %s", firstSeason.Name))

	// 3. Парсим турниры только первого сезона
	tournaments, err := s.juniorService.ParseSeasonTournaments(
		ctx,
		domain,
		firstSeason.Name,
		firstSeason.AjaxURL,
	)
	if err != nil {
		return false, fmt.Errorf("failed to parse season: %w", err)
	}

	if len(tournaments) == 0 {
		logger.Info(ctx, "    ℹ️  No tournaments in latest season")
		return false, nil // Пустой сезон, но не дубликат
	}

	// 4. Проверяем каждый турнир через globalDedup
	duplicateCount := 0
	var uniqueTournamentIDs []string

	logger.Info(ctx, fmt.Sprintf("    🔍 Checking %d tournaments for duplicates:", len(tournaments)))

	for _, t := range tournaments {
		if _, exists := globalDedup.Load(t.ID); exists {
			duplicateCount++
			logger.Debug(ctx, fmt.Sprintf("      🔁 DUPLICATE: ID=%s, Name=%s", t.ID, t.Name))
		} else {
			uniqueTournamentIDs = append(uniqueTournamentIDs, t.ID)
			logger.Debug(ctx, fmt.Sprintf("      ✨ UNIQUE: ID=%s, Name=%s", t.ID, t.Name))
		}
	}

	// 5. Вычисляем процент дубликатов
	totalTournaments := len(tournaments)
	duplicatePercentage := float64(duplicateCount) / float64(totalTournaments) * 100

	logger.Info(ctx, fmt.Sprintf("    📊 Check result: %d/%d duplicates (%.0f%%)",
		duplicateCount, totalTournaments, duplicatePercentage))

	// 6. Принимаем решение: если 100% дубликаты → пропускаем домен
	if duplicatePercentage >= 100.0 {
		logger.Info(ctx, "    🔁 100% duplicates - domain will be skipped")
		return true, nil
	}

	// 7. Домен уникальный - НЕ регистрируем здесь!
	// Регистрация произойдет в processDomain после ParseAllSeasonsTournaments
	// чтобы не потерять турниры первого сезона
	logger.Info(ctx, fmt.Sprintf("    ✅ Unique domain (%d unique tournaments found)", len(uniqueTournamentIDs)))
	return false, nil
}
