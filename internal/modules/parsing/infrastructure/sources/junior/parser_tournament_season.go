package junior

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/parsing"
	"github.com/PuerkitoBio/goquery"
)

var yearFromTextRe = regexp.MustCompile(`\b(20\d{2})\b`)

// ExtractAllSeasons извлекает все сезоны из дропдауна на странице /tournaments/
func (c *Client) ExtractAllSeasons(domain string) ([]SeasonInfo, error) {
	tournamentURL := domain + "/tournaments/"

	resp, err := c.MakeRequest(tournamentURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	// Первый проход: собираем данные и birth years по parentURL
	type tournamentData struct {
		parentURL  string
		name       string
		startDate  string
		endDate    string
		isEnded    bool
		birthYears map[int]bool
	}

	tournamentsMap := make(map[string]*tournamentData) // parentURL → data
	var order []string                                 // порядок обнаружения

	doc.Find(`a.comp-age[href^="/tournaments/"]`).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || href == "/tournaments/" || href == "/tournaments" {
			return
		}

		parentURL := href
		if idx := strings.Index(href, "-year-"); idx != -1 {
			parentURL = href[:idx] + "/"
		}

		// Извлекаем год из текста ссылки ("2015 г.р." → 2015)
		var birthYear int
		if matches := yearFromTextRe.FindStringSubmatch(strings.TrimSpace(s.Text())); len(matches) > 1 {
			birthYear, _ = strconv.Atoi(matches[1])
		}

		td, exists := tournamentsMap[parentURL]
		if !exists {
			name := ""
			s.Parents().EachWithBreak(func(j int, parent *goquery.Selection) bool {
				if parent.HasClass("comp-card") {
					name = strings.TrimSpace(parent.Find(".comp-link").First().Text())
					return name == ""
				}
				return true
			})
			if name == "" {
				name = strings.TrimSpace(s.Text())
			}

			startDate, endDate, isEnded := parsing.ParseTournamentMetadata(s)
			td = &tournamentData{
				parentURL:  parentURL,
				name:       name,
				startDate:  startDate,
				endDate:    endDate,
				isEnded:    isEnded,
				birthYears: make(map[int]bool),
			}
			tournamentsMap[parentURL] = td
			order = append(order, parentURL)
		}

		if birthYear > 0 {
			td.birthYears[birthYear] = true
		}
	})

	// Второй проход: создаём TournamentDTO с заполненными BirthYears
	tournaments := make([]TournamentDTO, 0, len(order))
	for _, parentURL := range order {
		td := tournamentsMap[parentURL]

		tournamentID := parsing.ExtractTournamentID(parentURL)
		if tournamentID == "" {
			tournamentID = parentURL
		}

		var birthYears []int
		for y := range td.birthYears {
			birthYears = append(birthYears, y)
		}
		sort.Ints(birthYears)

		tournaments = append(tournaments, TournamentDTO{
			ID:         tournamentID,
			Name:       td.name,
			URL:        parentURL,
			Domain:     domain,
			Season:     season,
			StartDate:  td.startDate,
			EndDate:    td.endDate,
			IsEnded:    td.isEnded,
			BirthYears: birthYears,
		})
	}

	return tournaments, nil
}
