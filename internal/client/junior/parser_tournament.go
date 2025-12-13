package junior

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/parsing"
)

// ParseTournamentsFromDomain парсит турниры с конкретного домена
func (c *Client) ParseTournamentsFromDomain(domain string) ([]TournamentDTO, error) {
	tournamentURL := domain + "/tournaments/"

	resp, err := c.makeRequest(tournamentURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	// Извлекаем глобальный сезон (один для всех турниров на странице)
	globalSeason := parsing.ExtractGlobalSeason(doc)

	var tournaments []TournamentDTO
	parsedURLs := make(map[string]bool) // Дедупликация

	// Парсим турниры (comp-card блоки)
	doc.Find(`a.comp-age[href^="/tournaments/"]`).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		// Пропускаем саму страницу /tournaments/
		if href == "/tournaments/" || href == "/tournaments" {
			return
		}

		// Убираем часть -year-XXXXX чтобы получить родительский URL
		parentURL := href
		if strings.Contains(href, "-year-") {
			yearIndex := strings.Index(href, "-year-")
			if yearIndex != -1 {
				parentURL = href[:yearIndex] + "/"
			}
		}

		// Дедупликация
		if parsedURLs[parentURL] {
			return
		}
		parsedURLs[parentURL] = true

		// Извлекаем название из родительского блока comp-card
		name := ""
		s.Parents().EachWithBreak(func(j int, parent *goquery.Selection) bool {
			if parent.HasClass("comp-card") {
				name = strings.TrimSpace(parent.Find(".comp-link").First().Text())
				if name != "" {
					return false
				}
			}
			return true
		})

		// Fallback: текст ссылки
		if name == "" {
			name = strings.TrimSpace(s.Text())
		}

		// Извлекаем ID турнира
		tournamentID := parsing.ExtractTournamentID(parentURL)
		if tournamentID == "" {
			tournamentID = parentURL
		}

		// Парсим даты из comp-period
		startDate, endDate, isEnded := parsing.ParseTournamentMetadata(s)

		// Если глобальный сезон не найден - пропускаем турнир
		if globalSeason == "" {
			return
		}

		tournament := TournamentDTO{
			ID:        tournamentID,
			Name:      name,
			URL:       parentURL,
			Domain:    domain,
			Season:    globalSeason, // ← Глобальный сезон
			StartDate: startDate,
			EndDate:   endDate,
			IsEnded:   isEnded,
		}

		tournaments = append(tournaments, tournament)
	})

	return tournaments, nil
}
