package stats

import (
	"context"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
)

// RunForTournaments парсит статистику для указанных турниров
func (s *service) RunForTournaments(ctx context.Context, tournaments []*entities.Tournament) error {
	if len(tournaments) == 0 {
		s.logger.Println("📋 Нет турниров для парсинга")
		return nil
	}

	startTime := time.Now()
	totalTournaments := len(tournaments)
	s.logger.Printf("🏒 Запуск парсинга статистики для %d турниров...\n", totalTournaments)

	const workerCount = 5
	s.logger.Printf("🚀 Запускаем %d параллельных воркеров...\n", workerCount)

	pool := NewTournamentWorkerPool(ctx, s, workerCount)
	pool.Start()

	go func() {
		for i, t := range tournaments {
			pool.AddTask(TournamentTask{
				Tournament: t,
				Index:      i + 1,
				Total:      totalTournaments,
			})
		}
		pool.Close()
	}()

	go pool.Wait()

	totalStats := 0
	successCount := 0
	errorCount := 0

	for result := range pool.Results() {
		if result.Error != nil {
			errorCount++
		} else {
			totalStats += result.Count
			successCount++
		}
	}

	duration := time.Since(startTime)
	s.logger.Println("\n" + strings.Repeat("=", 60))
	s.logger.Println("🎉 Парсинг статистики завершен!")
	s.logger.Printf("📊 Успешно: %d/%d, Ошибок: %d", successCount, totalTournaments, errorCount)
	s.logger.Printf("📊 Записей статистики: %d", totalStats)
	s.logger.Printf("⏱️  Время: %s", duration.Round(time.Second))
	s.logger.Println(strings.Repeat("=", 60))

	return nil
}
