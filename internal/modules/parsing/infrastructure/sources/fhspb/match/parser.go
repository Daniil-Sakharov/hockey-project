package match

import (
	"bytes"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

var (
	timeRegex   = regexp.MustCompile(`(\d{1,2}):(\d{2})`)
	scoreRegex  = regexp.MustCompile(`(\d+):(\d+)`)
	numberRegex = regexp.MustCompile(`^(\d+)\s+`)
)

// Parser парсер протокола матча FHSPB
type Parser struct{}

// NewParser создает новый парсер протокола
func NewParser() *Parser {
	return &Parser{}
}

// Parse парсит HTML страницы протокола матча
func (p *Parser) Parse(html []byte) (*MatchDetailsDTO, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	details := &MatchDetailsDTO{}

	// Парсим счёт по периодам из #ScoreGridView
	p.parsePeriodScores(doc, details)

	// Парсим броски из #ShotGridView
	p.parseShots(doc, details)

	// Парсим голы из секции h3:contains("Голы")
	details.Goals = p.parseGoals(doc, details)

	// Парсим штрафы из секции h3:contains("Удаления")
	details.Penalties = p.parsePenalties(doc)

	// Парсим составы из секций h5:contains("Полевые игроки")
	p.parseLineups(doc, details)

	// Парсим статистику вратарей из секций h5:contains("Вратари")
	p.parseGoalieStats(doc, details)

	return details, nil
}
