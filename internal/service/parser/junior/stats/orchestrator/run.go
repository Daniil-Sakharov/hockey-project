package stats_orchestrator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Run главный метод - парсит статистику всех турниров из БД параллельно
func (s *service) Run(ctx context.Context) error {
	startTime := time.Now()
	s.logger.Println("🏒 Запуск параллельного парсинга статистики турниров...")

	// 1. Очищаем таблицу статистики
	s.logger.Println("🗑️  Очистка таблицы player_statistics...")
	if err := s.statsRepo.DeleteAll(ctx); err != nil {
		return fmt.Errorf("failed to truncate statistics table: %w", err)
	}
	s.logger.Println("✅ Таблица очищена")

	// 2. Получаем все турниры из БД
	tournaments, err := s.tournamentRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tournaments: %w", err)
	}

	totalTournaments := len(tournaments)
	s.logger.Printf("📋 Найдено турниров: %d\n", totalTournaments)

	// 3. Создаём Worker Pool для параллельного парсинга
	const workerCount = 5
	s.logger.Printf("🚀 Запускаем %d параллельных воркеров...\n", workerCount)

	pool := NewTournamentWorkerPool(ctx, s, workerCount)
	pool.Start()

	// 4. Добавляем турниры в очередь
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

	// 5. Ждём завершения
	go pool.Wait()

	// 6. Собираем результаты
	totalStats := 0
	successCount := 0
	errorCount := 0

	type tournamentStats struct {
		name  string
		id    string
		count int
	}
	tournamentResults := make([]tournamentStats, 0, totalTournaments)

	for result := range pool.Results() {
		if result.Error != nil {
			errorCount++
			tournamentResults = append(tournamentResults, tournamentStats{
				name:  result.Tournament.Name,
				id:    result.Tournament.ID,
				count: 0,
			})
		} else {
			totalStats += result.Count
			successCount++
			tournamentResults = append(tournamentResults, tournamentStats{
				name:  result.Tournament.Name,
				id:    result.Tournament.ID,
				count: result.Count,
			})
		}
	}

	// 7. Финальная статистика
	duration := time.Since(startTime)
	s.logger.Println("\n" + strings.Repeat("=", 60))
	s.logger.Println("🎉 Парсинг статистики завершен!")
	s.logger.Printf("📊 Успешно обработано турниров: %d/%d", successCount, totalTournaments)
	s.logger.Printf("📊 Всего записей статистики: %d", totalStats)
	s.logger.Printf("❌ Ошибок: %d", errorCount)
	s.logger.Printf("⏱️  Время выполнения: %s (с %d воркерами)", duration.Round(time.Second), workerCount)
	s.logger.Println(strings.Repeat("=", 60))

	// 8. Проверяем итоговое количество записей в БД
	finalCount, err := s.statsRepo.CountAll(ctx)
	if err != nil {
		s.logger.Printf("⚠️  Не удалось подсчитать записи в БД: %v", err)
	} else {
		s.logger.Printf("✅ Записей в БД: %d", finalCount)
	}

	// 9. Детальная статистика по турнирам
	s.logger.Println("\n" + strings.Repeat("=", 60))
	s.logger.Println("📊 ДЕТАЛЬНАЯ СТАТИСТИКА ПО ТУРНИРАМ:")
	s.logger.Println(strings.Repeat("=", 60))

	for i := 0; i < len(tournamentResults)-1; i++ {
		for j := i + 1; j < len(tournamentResults); j++ {
			if tournamentResults[j].count > tournamentResults[i].count {
				tournamentResults[i], tournamentResults[j] = tournamentResults[j], tournamentResults[i]
			}
		}
	}

	for i, tr := range tournamentResults {
		if tr.count > 0 {
			s.logger.Printf("%2d. [%s] %s: %d записей", i+1, tr.id, tr.name, tr.count)
		} else {
			s.logger.Printf("%2d. [%s] %s: 0 записей ❌", i+1, tr.id, tr.name)
		}
	}
	s.logger.Println(strings.Repeat("=", 60))

	return nil
}
