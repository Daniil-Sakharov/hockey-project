package parsing

import (
	"regexp"
)

// Общие regex для парсинга
var (
	PlayerIDRegex = regexp.MustCompile(`PlayerID=([a-f0-9-]+)`)
	TeamIDRegex   = regexp.MustCompile(`TeamID=([a-f0-9-]+)`)
)
