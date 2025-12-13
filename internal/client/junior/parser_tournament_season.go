package junior

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/parsing"
)

// ExtractAllSeasons извлекает все сезоны из дропдауна на странице /tournaments/
func (c *Client) ExtractAllSeasons(domain string) ([]SeasonInfo, error) {
	tournamentURL := domain + "/tournaments/"

	resp, err := c.MakeRequest(tournamentURL)
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

	var seasons []SeasonInfo

	// Парсим все <option> из дропдауна сезонов
	doc.Find("select.js-ajax-select option").Each(func(i int, s *goquery.Selection) {
		seasonName := s.AttrOr("value", "")
		ajaxURL := s.AttrOr("data-ajax", "")

		if seasonName != "" && ajaxURL != "" {
			seasons = append(seasons, SeasonInfo{
				Name:    seasonName,
				AjaxURL: ajaxURL,
			})
		}
	})

	return seasons, nil
}

// ParseSeasonTournaments парсит турниры одного сезона через AJAX запрос
func (c *Client) ParseSeasonTournaments(domain, season, ajaxURL string) ([]TournamentDTO, error) {
	fullURL := domain + ajaxURL

	resp, err := c.MakeRequest(fullURL)
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

	var tournaments []TournamentDTO
	parsedURLs := make(map[string]bool) // Дедупликация

	// Парсим турниры (comp-card блоки) - аналогично ParseTournamentsFromDomain
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

		tournament := TournamentDTO{
			ID:        tournamentID,
			Name:      name,
			URL:       parentURL,
			Domain:    domain,
			Season:    season, // ← Сезон уже известен из параметра
			StartDate: startDate,
			EndDate:   endDate,
			IsEnded:   isEnded,
		}

		tournaments = append(tournaments, tournament)
	})

	return tournaments, nil
}
