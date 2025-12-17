package domain

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер доменов
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер доменов
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// DiscoverAll собирает все уникальные домены *.fhr.ru
func (p *Parser) DiscoverAll(mainURL string) ([]string, error) {
	resp, err := p.http.MakeRequest(mainURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса к %s: %w", mainURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус %d для %s", resp.StatusCode, mainURL)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	discoveredDomains := make(map[string]bool)

	doc.Find(`a[href*=".fhr.ru"]`).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}
		domain := extractDomain(href)
		if domain != "" && strings.Contains(domain, ".fhr.ru") {
			discoveredDomains[domain] = true
		}
	})

	var domains []string
	for domain := range discoveredDomains {
		domains = append(domains, domain)
	}
	domains = append([]string{mainURL}, domains...)

	return domains, nil
}

func extractDomain(urlStr string) string {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}
