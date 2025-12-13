package parsing

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// StatsCombination представляет комбинацию год+группа для парсинга
type StatsCombination struct {
	YearID    string // ID года в системе
	YearLabel string // Текст года (например "2009")
	GroupID   string // ID группы ("all" для общей статистики)
	GroupName string // Название группы
}

// YearInfo информация о годе из dropdown
type YearInfo struct {
	ID       string // ID года (value из option)
	Label    string // Текст года (2009, 2010...)
	AjaxURL  string // URL для AJAX запроса
}

// ParseCombinationsWithAjax парсит комбинации год+группа с двухуровневой логикой
// 1. Извлекает все ГОДЫ из dropdown
// 2. Для КАЖДОГО года делает AJAX запрос
// 3. Из AJAX ответа извлекает ГРУППЫ для этого года
func ParseCombinationsWithAjax(
	ctx context.Context,
	doc *goquery.Document,
	domain string,
	httpClient *http.Client,
) ([]StatsCombination, error) {
	combinationsMap := make(map[string]StatsCombination)

	// Шаг 1: Извлекаем все ГОДЫ из dropdown
	years := extractYearsFromDoc(doc)
	
	if len(years) == 0 {
		// Нет dropdown годов - парсим группы с текущей страницы (старая логика)
		return parseGroupsFromDoc(doc, "")
	}

	// Шаг 2: Для каждого года делаем AJAX запрос и извлекаем группы
	for _, year := range years {
		// Делаем AJAX запрос для года
		yearDoc, err := fetchYearPage(ctx, httpClient, domain, year.AjaxURL)
		if err != nil {
			// Логируем ошибку но продолжаем
			continue
		}

		// Извлекаем группы из AJAX ответа
		groups := extractGroupsFromDoc(yearDoc)

		if len(groups) == 0 {
			// Нет групп для этого года - добавляем "Общую статистику"
			key := year.ID + "|all"
			combinationsMap[key] = StatsCombination{
				YearID:    year.ID,
				YearLabel: year.Label,
				GroupID:   "all",
				GroupName: "Общая статистика",
			}
		} else {
			// Добавляем все группы для этого года
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

	// Преобразуем map в slice
	combinations := make([]StatsCombination, 0, len(combinationsMap))
	for _, combo := range combinationsMap {
		combinations = append(combinations, combo)
	}

	return combinations, nil
}

// GroupInfo информация о группе
type GroupInfo struct {
	ID   string
	Name string
}

// extractYearsFromDoc извлекает годы из dropdown
// ВАЖНО: используем select[data-ajax-select] или select[name="tech"] - это dropdown годов рождения
func extractYearsFromDoc(doc *goquery.Document) []YearInfo {
	var years []YearInfo
	seen := make(map[string]bool)

	// Ищем только в правильном select (годы рождения турнира)
	// select[data-ajax-select] - основной селектор для годов рождения
	// select[name="tech"] - альтернативный селектор
	doc.Find("select[data-ajax-select] option[data-ajax], select[name='tech'] option[data-ajax]").Each(func(i int, s *goquery.Selection) {
		value, hasValue := s.Attr("value")
		dataAjax, hasAjax := s.Attr("data-ajax")
		text := strings.TrimSpace(s.Text())

		if !hasValue || !hasAjax || value == "" {
			return
		}

		// Проверяем что это competitions-stats
		if !strings.Contains(dataAjax, "competitions-stats") {
			return
		}

		// Дедупликация
		if seen[value] {
			return
		}
		seen[value] = true

		// Определяем label года
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
		if !exists {
			return
		}

		// Проверяем что это competitions-stats
		if !strings.Contains(ajaxLink, "competitions-stats") {
			return
		}

		// Извлекаем params из URL
		re := regexp.MustCompile(`params=([^&]+)`)
		matches := re.FindStringSubmatch(ajaxLink)
		if len(matches) < 2 {
			return
		}

		// Декодируем GROUP_ID
		_, groupID := DecodeYearAndGroupID(matches[1])
		if groupID == "" {
			groupID = "all"
		}

		// Дедупликация
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// parseGroupsFromDoc парсит группы из текущего документа (fallback для турниров без годов)
func parseGroupsFromDoc(doc *goquery.Document, defaultYearID string) ([]StatsCombination, error) {
	combinationsMap := make(map[string]StatsCombination)

	// Парсим комбинации из кнопок групп (.filter-btn)
	doc.Find(".filter-block .filter-btn[data-ajax-link], div.filter-btn[data-ajax-link]").Each(func(i int, s *goquery.Selection) {
		groupName := strings.TrimSpace(s.Text())
		if groupName == "" {
			return
		}

		ajaxLink, exists := s.Attr("data-ajax-link")
		if !exists {
			return
		}

		// Извлекаем params из URL
		re := regexp.MustCompile(`params=([^&]+)`)
		matches := re.FindStringSubmatch(ajaxLink)
		if len(matches) < 2 {
			return
		}

		// Декодируем и извлекаем YEAR_ID и GROUP_ID
		yearID, groupID := DecodeYearAndGroupID(matches[1])
		if yearID == "" {
			yearID = defaultYearID
		}

		// Пытаемся извлечь год из текста года рождения в HTML
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

// ParseCombinations парсит комбинации год+группа из начального HTML
// DEPRECATED: Используйте ParseCombinationsWithAjax для полного парсинга
// Эта функция оставлена для обратной совместимости
func ParseCombinations(doc *goquery.Document) ([]StatsCombination, error) {
	return parseGroupsFromDoc(doc, "")
}
