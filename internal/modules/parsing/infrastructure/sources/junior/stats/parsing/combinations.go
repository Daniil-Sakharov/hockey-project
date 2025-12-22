package parsing

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseCombinationsWithAjax парсит комбинации год+группа с двухуровневой логикой
func ParseCombinationsWithAjax(
	ctx context.Context,
	doc *goquery.Document,
	domain string,
	httpClient *http.Client,
) ([]StatsCombination, error) {
	combinationsMap := make(map[string]StatsCombination)

	years := extractYearsFromDoc(doc)

	if len(years) == 0 {
		return parseGroupsFromDoc(doc, "")
	}

	for _, year := range years {
		yearDoc, err := fetchYearPage(ctx, httpClient, domain, year.AjaxURL)
		if err != nil {
			continue
		}

		groups := extractGroupsFromDoc(yearDoc)

		if len(groups) == 0 {
			key := year.ID + "|all"
			combinationsMap[key] = StatsCombination{
				YearID:    year.ID,
				YearLabel: year.Label,
				GroupID:   "all",
				GroupName: "Общая статистика",
			}
		} else {
			for _, group := range groups {
				key := year.ID + "|" + group.ID
				combinationsMap[key] = StatsCombination{
					YearID:    year.ID,
					YearLabel: year.Label,
					GroupID:   group.ID,
					GroupName: group.Name,
				}
			}
		}
	}

	combinations := make([]StatsCombination, 0, len(combinationsMap))
	for _, combo := range combinationsMap {
		combinations = append(combinations, combo)
	}

	return combinations, nil
}

// fetchYearPage делает AJAX запрос для получения страницы года
func fetchYearPage(ctx context.Context, httpClient *http.Client, domain, ajaxURL string) (*goquery.Document, error) {
	fullURL := domain + ajaxURL
	if strings.HasPrefix(ajaxURL, "http") {
		fullURL = ajaxURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch year page: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// parseGroupsFromDoc парсит группы из текущего документа (fallback)
func parseGroupsFromDoc(doc *goquery.Document, defaultYearID string) ([]StatsCombination, error) {
	combinationsMap := make(map[string]StatsCombination)

	doc.Find(".filter-block .filter-btn[data-ajax-link], div.filter-btn[data-ajax-link]").Each(func(i int, s *goquery.Selection) {
		groupName := strings.TrimSpace(s.Text())
		if groupName == "" {
			return
		}

		ajaxLink, exists := s.Attr("data-ajax-link")
		if !exists {
			return
		}

		re := regexp.MustCompile(`params=([^&]+)`)
		matches := re.FindStringSubmatch(ajaxLink)
		if len(matches) < 2 {
			return
		}

		yearID, groupID := DecodeYearAndGroupID(matches[1])
		if yearID == "" {
			yearID = defaultYearID
		}

		yearLabel := ExtractYearLabelByID(doc, yearID)
		if yearLabel == "" && len(yearID) >= 4 {
			yearLabel = "Год_" + yearID[:4]
		}

		key := yearID + "|" + groupID
		combinationsMap[key] = StatsCombination{
			YearID:    yearID,
			YearLabel: yearLabel,
			GroupID:   groupID,
			GroupName: groupName,
		}
	})

	combinations := make([]StatsCombination, 0, len(combinationsMap))
	for _, combo := range combinationsMap {
		combinations = append(combinations, combo)
	}

	return combinations, nil
}

// ParseCombinations парсит комбинации (DEPRECATED: используйте ParseCombinationsWithAjax)
func ParseCombinations(doc *goquery.Document) ([]StatsCombination, error) {
	return parseGroupsFromDoc(doc, "")
}
