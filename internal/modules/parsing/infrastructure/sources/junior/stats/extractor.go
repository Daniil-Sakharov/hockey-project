package stats

import (
	"regexp"
	"strconv"
	"strings"
)

// ExtractPlayerID извлекает player_id из HTML href
// Пример: <a href="/player/tikhomirov-...-674476/">... → "674476"
func ExtractPlayerID(html string) string {
	re := regexp.MustCompile(`/player/[^/]+-(\d+)/`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// ExtractTeamID извлекает team_id из HTML href команды
// Пример: /tournaments/.../polet_1412116/ → "1412116"
func ExtractTeamID(html string) string {
	re := regexp.MustCompile(`/tournaments/[^/]+/[^_]+_(\d+)/`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// ParseInt парсит HTML в int
// Пример: <div class="cell vert-sort">11</div> → 11
func ParseInt(html string) int {
	// Извлекаем число из HTML
	re := regexp.MustCompile(`>(\d+)<`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		val, err := strconv.Atoi(matches[1])
		if err == nil {
			return val
		}
	}
	return 0
}

// ParseFloat парсит HTML в float64
// Пример: <div class="cell vert-sort">1,18</div> → 1.18
func ParseFloat(html string) float64 {
	// Извлекаем число из HTML
	re := regexp.MustCompile(`>([\d,]+)<`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		// Заменяем запятую на точку
		numStr := strings.ReplaceAll(matches[1], ",", ".")
		val, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			return val
		}
	}
	return 0.0
}
