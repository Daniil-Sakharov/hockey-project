package junior

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/pool"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// ParseAllSeasonsTournaments парсит турниры ВСЕХ сезонов домена через Worker Pool
func (c *Client) ParseAllSeasonsTournaments(ctx context.Context, domain string) ([]TournamentDTO, error) {
	// 1. Извлекаем все сезоны из дропдауна
	seasons, err := c.ExtractAllSeasons(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to extract seasons: %w", err)
	}

	if len(seasons) == 0 {
		return nil, fmt.Errorf("no seasons found for domain %s", domain)
	}

	logger.Info(ctx, fmt.Sprintf("    📅 Found %d seasons", len(seasons)))

	// 2. Создаем Worker Pool (10 воркеров)
	workerPool := pool.NewWorkerPool(ctx, c, 10)
	workerPool.Start()

	// 3. Добавляем задачи в очередь (в отдельной goroutine)
	go func() {
		for _, season := range seasons {
			workerPool.AddTask(pool.SeasonTask{
				Domain:  domain,
				Season:  season.Name,
				AjaxURL: season.AjaxURL,
			})
		}
		workerPool.Close() // Закрываем очередь задач
	}()

	// 4. Ждем завершения воркеров (в отдельной goroutine)
	go workerPool.Wait() // Закроет results канал после завершения

	// 5. Собираем результаты
	var allTournaments []TournamentDTO
	dedup := make(map[string]bool) // Дедупликация по URL

	successCount := 0
	errorCount := 0

	for result := range workerPool.Results() {
		if result.Error != nil {
			// Логируем ошибку и продолжаем (не падаем)
			logger.Warn(ctx, fmt.Sprintf("    ⚠️  Worker error [%s]: %v",
				result.Task.Season, result.Error))
			errorCount++
			continue
		}

		// Добавляем турниры с дедупликацией
		addedCount := 0
		for _, tournament := range result.Tournaments {
			if !dedup[tournament.URL] {
				dedup[tournament.URL] = true
				allTournaments = append(allTournaments, tournament)
				addedCount++
			}
		}

		successCount++
		logger.Info(ctx, fmt.Sprintf("    ✅ Parsed %s: %d tournaments (added %d unique)",
			result.Task.Season, len(result.Tournaments), addedCount))
	}

	// 6. Финальная статистика
	logger.Info(ctx, fmt.Sprintf("    📊 Total: %d unique tournaments from %d/%d seasons",
		len(allTournaments), successCount, len(seasons)))

	if errorCount > 0 {
		logger.Warn(ctx, fmt.Sprintf("    ⚠️  Errors: %d seasons failed", errorCount))
	}

	return allTournaments, nil
}
