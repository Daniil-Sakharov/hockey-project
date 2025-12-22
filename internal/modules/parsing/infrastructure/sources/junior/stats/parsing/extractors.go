package parsing

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// extractYearsFromDoc извлекает годы из dropdown
func extractYearsFromDoc(doc *goquery.Document) []YearInfo {
	var years []YearInfo
	seen := make(map[string]bool)

	doc.Find("select[data-ajax-select] option[data-ajax], select[name='tech'] option[data-ajax]").Each(func(i int, s *goquery.Selection) {
		value, hasValue := s.Attr("value")
		dataAjax, hasAjax := s.Attr("data-ajax")
		text := strings.TrimSpace(s.Text())

		if !hasValue || !hasAjax || value == "" {
			return
		}

		if !strings.Contains(dataAjax, "competitions-stats") {
			return
		}

		if seen[value] {
			return
		}
		seen[value] = true

		yearLabel := text
		if len(text) != 4 || !IsYear(text) {
			yearLabel = "Год_" + value
			if len(value) >= 4 {
				yearLabel = "Год_" + value[:4]
			}
		}

		years = append(years, YearInfo{
			ID:      value,
			Label:   yearLabel,
			AjaxURL: dataAjax,
		})
	})

	return years
}

// extractGroupsFromDoc извлекает группы из HTML документа (AJAX ответа)
func extractGroupsFromDoc(doc *goquery.Document) []GroupInfo {
	var groups []GroupInfo
	seen := make(map[string]bool)

	doc.Find(".filter-block .filter-btn[data-ajax-link], div.filter-btn[data-ajax-link]").Each(func(i int, s *goquery.Selection) {
		groupName := strings.TrimSpace(s.Text())
		if groupName == "" {
			return
		}

		ajaxLink, exists := s.Attr("data-ajax-link")
		if !exists || !strings.Contains(ajaxLink, "competitions-stats") {
			return
		}

		re := regexp.MustCompile(`params=([^&]+)`)
		matches := re.FindStringSubmatch(ajaxLink)
		if len(matches) < 2 {
			return
		}

		_, groupID := DecodeYearAndGroupID(matches[1])
		if groupID == "" {
			groupID = "all"
		}

		if seen[groupID] {
			return
		}
		seen[groupID] = true

		groups = append(groups, GroupInfo{
			ID:   groupID,
			Name: groupName,
		})
	})

	return groups
}
