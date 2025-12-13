package parsing

import (
	"encoding/base64"
	"net/url"
	"regexp"
)

// DecodeYearAndGroupID декодирует YEAR_ID и GROUP_ID из base64 params
func DecodeYearAndGroupID(params string) (yearID string, groupID string) {
	// Декодируем URL encoding
	decoded, err := url.QueryUnescape(params)
	if err != nil {
		return "", ""
	}

	// Декодируем base64
	data, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		return "", ""
	}

	serialized := string(data)

	// Извлекаем YEAR_ID: s:7:"YEAR_ID";s:X:"VALUE"
	reYear := regexp.MustCompile(`s:7:"YEAR_ID";s:\d+:"([^"]+)"`)
	matchesYear := reYear.FindStringSubmatch(serialized)
	if len(matchesYear) > 1 {
		yearID = matchesYear[1]
	}

	// Извлекаем GROUP_ID: s:8:"GROUP_ID";s:X:"VALUE" или s:8:"GROUP_ID";N (NULL)
	reGroup := regexp.MustCompile(`s:8:"GROUP_ID";s:\d+:"([^"]+)"`)
	matchesGroup := reGroup.FindStringSubmatch(serialized)
	if len(matchesGroup) > 1 {
		groupID = matchesGroup[1]
	} else {
		// Проверяем NULL или отсутствие GROUP_ID (означает "Общая статистика")
		groupID = "all"
	}

	return yearID, groupID
}
