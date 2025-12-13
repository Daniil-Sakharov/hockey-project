package filter

import (
	"strconv"
	"strings"

	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
)

// parseHeightRange парсит диапазон роста из строки
func parseHeightRange(value string) *domainBot.HeightRange {
	if strings.Contains(value, "+") {
		// 200+
		min, _ := strconv.Atoi(strings.TrimSuffix(value, "+"))
		return &domainBot.HeightRange{Min: min, Max: 300}
	}

	// 150-160
	parts := strings.Split(value, "-")
	if len(parts) == 2 {
		min, _ := strconv.Atoi(parts[0])
		max, _ := strconv.Atoi(parts[1])
		return &domainBot.HeightRange{Min: min, Max: max}
	}

	return nil
}

// parseWeightRange парсит диапазон веса из строки
func parseWeightRange(value string) *domainBot.WeightRange {
	if strings.Contains(value, "+") {
		// 90+
		min, _ := strconv.Atoi(strings.TrimSuffix(value, "+"))
		return &domainBot.WeightRange{Min: min, Max: 200}
	}

	// 50-60
	parts := strings.Split(value, "-")
	if len(parts) == 2 {
		min, _ := strconv.Atoi(parts[0])
		max, _ := strconv.Atoi(parts[1])
		return &domainBot.WeightRange{Min: min, Max: max}
	}

	return nil
}
