package profile

import (
	"fmt"
	"time"
)

// getCurrentSeason определяет текущий сезон
// Если месяц >= июль (7) - новый сезон начался
func getCurrentSeason() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()

	if month >= time.July {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}

// getPreviousSeason возвращает предыдущий сезон
func getPreviousSeason(currentSeason string) string {
	// Парсим текущий сезон "2024-2025" -> 2024
	var startYear int
	_, _ = fmt.Sscanf(currentSeason, "%d-", &startYear)
	return fmt.Sprintf("%d-%d", startYear-1, startYear)
}

// calculateAverage вычисляет среднее значение с точностью 2 знака
func calculateAverage(value, games int) float64 {
	if games == 0 {
		return 0
	}
	return float64(value) / float64(games)
}
