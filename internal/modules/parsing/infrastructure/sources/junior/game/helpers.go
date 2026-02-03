package game

import (
	"regexp"
	"strconv"
	"strings"
)

var gameIDRegex = regexp.MustCompile(`/games/(\d+)`)

func extractGameID(url string) string {
	m := gameIDRegex.FindStringSubmatch(url)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func parseScoreText(text string) (home, away int, ok bool) {
	text = strings.TrimSpace(text)
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	h, e1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	a, e2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if e1 != nil || e2 != nil {
		return 0, 0, false
	}
	return h, a, true
}

var timeRegex = regexp.MustCompile(`(\d+):(\d+)`)
var periodRegex = regexp.MustCompile(`(\d)\s*период`)

func parseGameTime(text string) (period, min, sec int) {
	if m := periodRegex.FindStringSubmatch(text); len(m) >= 2 {
		period, _ = strconv.Atoi(m[1])
	}
	if m := timeRegex.FindStringSubmatch(text); len(m) >= 3 {
		min, _ = strconv.Atoi(m[1])
		sec, _ = strconv.Atoi(m[2])
	}
	return
}
