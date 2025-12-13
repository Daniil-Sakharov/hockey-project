package parsing

import (
	"regexp"
	"strconv"
	"strings"
)

// extractValue извлекает значение после метки
func extractValue(text, label string) string {
	idx := strings.Index(text, label)
	if idx == -1 {
		return ""
	}

	rest := text[idx+len(label):]
	rest = strings.TrimLeft(rest, " :\t│")

	endIdx := strings.IndexAny(rest, "│\n")
	if endIdx > 0 && endIdx < 100 {
		rest = rest[:endIdx]
	} else if len(rest) > 100 {
		rest = rest[:100]
	}

	return strings.TrimSpace(rest)
}

// extractIntValue извлекает числовое значение
func extractIntValue(text, label string) int {
	val := extractValue(text, label)
	if val == "" {
		return 0
	}

	numRegex := regexp.MustCompile(`(\d+)`)
	if m := numRegex.FindStringSubmatch(val); len(m) >= 2 {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}
