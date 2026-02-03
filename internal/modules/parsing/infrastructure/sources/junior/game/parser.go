package game

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер деталей матча
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// Parse парсит детали матча
func (p *Parser) Parse(gameURL string) (*GameDetailsDTO, error) {
	resp, err := p.http.MakeRequest(gameURL)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	details := &GameDetailsDTO{}
	details.ExternalID = extractGameID(gameURL)

	// Парсинг года рождения и группы из заголовка
	p.parseTournamentHeader(doc, details)

	// Парсинг счёта по периодам
	p.parseScore(doc, details)

	// Парсинг видео
	details.VideoURL = p.parseVideoURL(doc)

	// Парсинг команд
	p.parseTeams(doc, details)

	// Парсинг протокола (голы, штрафы, вратари, таймауты)
	p.parseProtocol(doc, details)

	// Парсинг составов
	details.HomeLineup, details.AwayLineup = p.parseLineups(doc)

	return details, nil
}
