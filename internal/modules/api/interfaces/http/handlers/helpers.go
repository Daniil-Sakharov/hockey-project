package handlers

import "strings"

// Known abbreviations that should be preserved in titleCase.
var knownAbbreviations = map[string]string{
	"хк":   "ХК",
	"ск":   "СК",
	"ска":  "СКА",
	"цск":  "ЦСК",
	"цска": "ЦСКА",
	"сдюшор": "СДЮШОР",
	"уфо":  "УФО",
	"пфо":  "ПФО",
	"сфо":  "СФО",
	"юфо":  "ЮФО",
	"цфо":  "ЦФО",
	"сзфо": "СЗФО",
	"двфо": "ДВФО",
	"дфо":  "ДФО",
	"мо":   "МО",
	"ло":   "ЛО",
	"спб":  "СПб",
	"3х3":  "3х3",
}

// titleCase converts "ТОРПЕДО" → "Торпедо", preserving known abbreviations.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) == 0 {
			continue
		}
		lower := strings.ToLower(w)
		if abbr, ok := knownAbbreviations[lower]; ok {
			words[i] = abbr
			continue
		}
		runes := []rune(lower)
		runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// mapPositionToAPI maps DB position to API value.
func mapPositionToAPI(pos string) string {
	switch pos {
	case "Нападающий":
		return "forward"
	case "Защитник":
		return "defender"
	case "Вратарь":
		return "goalie"
	case "G":
		return "goalie"
	case "D":
		return "defender"
	case "F":
		return "forward"
	default:
		return pos
	}
}
