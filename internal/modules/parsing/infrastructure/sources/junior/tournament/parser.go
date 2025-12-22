package tournament

import (
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер турниров
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер турниров
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// ParseFromDomain парсит турниры с конкретного домена
func (p *Parser) ParseFromDomain(domain string) ([]types.TournamentDTO, error) {
	tournamentURL := domain + "/tournaments/"

	resp, err := p.http.MakeRequest(tournamentURL)
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

	globalSeason := parsing.ExtractGlobalSeason(doc)

	var tournaments []types.TournamentDTO
	parsedURLs := make(map[string]bool)

	doc.Find(`a.comp-age[href^="/tournaments/"]`).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || href == "/tournaments/" || href == "/tournaments" {
			return
		}

		parentURL := href
		if strings.Contains(href, "-year-") {
			if yearIndex := strings.Index(href, "-year-"); yearIndex != -1 {
				parentURL = href[:yearIndex] + "/"
			}
		}

		if parsedURLs[parentURL] {
			return
		}
		parsedURLs[parentURL] = true

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

		if name == "" {
			name = strings.TrimSpace(s.Text())
		}

		tournamentID := parsing.ExtractTournamentID(parentURL)
		if tournamentID == "" {
			tournamentID = parentURL
		}

		startDate, endDate, isEnded := parsing.ParseTournamentMetadata(s)

		if globalSeason == "" {
			return
		}

		tournaments = append(tournaments, types.TournamentDTO{
			ID:        tournamentID,
			Name:      name,
			URL:       parentURL,
			Domain:    domain,
			Season:    globalSeason,
			StartDate: startDate,
			EndDate:   endDate,
			IsEnded:   isEnded,
		})
	})

	return tournaments, nil
}
