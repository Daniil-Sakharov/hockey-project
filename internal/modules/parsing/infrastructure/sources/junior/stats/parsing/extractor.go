package parsing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// IsYear проверяет что строка является годом
func IsYear(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}

// ExtractYearLabelByID находит текст года (например "2009") по его ID (например "16743907")
func ExtractYearLabelByID(doc *goquery.Document, yearID string) string {
	var yearLabel string

	doc.Find("select.select-el option[data-ajax]").Each(func(i int, s *goquery.Selection) {
		value, hasValue := s.Attr("value")
		if hasValue && value == yearID {
			text := strings.TrimSpace(s.Text())
			if len(text) == 4 && IsYear(text) {
				yearLabel = text
			}
		}
	})

	return yearLabel
}

// ExtractGroupNameByID находит название группы по её ID из data-ajax-link кнопок
func ExtractGroupNameByID(doc *goquery.Document, groupID string) string {
	var groupName string

	doc.Find(".filter-block .filter-btn[data-ajax-link]").Each(func(i int, s *goquery.Selection) {
		ajaxLink, exists := s.Attr("data-ajax-link")
		if !exists {
			return
		}

		// Извлекаем params
		re := regexp.MustCompile(`params=([^&]+)`)
		matches := re.FindStringSubmatch(ajaxLink)
		if len(matches) < 2 {
			return
		}

		// Декодируем и проверяем GROUP_ID
		_, gid := DecodeYearAndGroupID(matches[1])
		if gid == groupID {
			groupName = strings.TrimSpace(s.Text())
		}
	})

	return groupName
}
