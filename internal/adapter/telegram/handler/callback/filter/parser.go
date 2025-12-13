package filter

import (
	"strconv"
	"strings"

	domainBot "github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
)

func parseHeightRange(value string) *domainBot.HeightRange {
	if strings.Contains(value, "+") {
		min, _ := strconv.Atoi(strings.TrimSuffix(value, "+"))
		return &domainBot.HeightRange{Min: min, Max: 300}
	}

	parts := strings.Split(value, "-")
	if len(parts) == 2 {
		min, _ := strconv.Atoi(parts[0])
		max, _ := strconv.Atoi(parts[1])
		return &domainBot.HeightRange{Min: min, Max: max}
	}

	return nil
}

func parseWeightRange(value string) *domainBot.WeightRange {
	if strings.Contains(value, "+") {
		min, _ := strconv.Atoi(strings.TrimSuffix(value, "+"))
		return &domainBot.WeightRange{Min: min, Max: 200}
	}

	parts := strings.Split(value, "-")
	if len(parts) == 2 {
		min, _ := strconv.Atoi(parts[0])
		max, _ := strconv.Atoi(parts[1])
		return &domainBot.WeightRange{Min: min, Max: max}
	}

	return nil
}
