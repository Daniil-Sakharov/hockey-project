package junior

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DiscoverAllDomains собирает все уникальные домены *.fhr.ru с главной страницы
func (c *Client) DiscoverAllDomains(mainURL string) ([]string, error) {
	resp, err := c.makeRequest(mainURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка HTTP запроса к %s: %w", mainURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус %d для %s", resp.StatusCode, mainURL)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	discoveredDomains := make(map[string]bool)

	// Ищем все ссылки, содержащие .fhr.ru
	doc.Find(`a[href*=".fhr.ru"]`).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Извлекаем чистый домен
		domain := extractDomain(href)
		if domain == "" {
			return
		}

		// Фильтруем: только домены, заканчивающиеся на .fhr.ru
		if strings.Contains(domain, ".fhr.ru") {
			discoveredDomains[domain] = true
		}
	})

	// Преобразуем map в slice
	var domains []string
	for domain := range discoveredDomains {
		domains = append(domains, domain)
	}

	// Добавляем сам junior.fhr.ru в список (он тоже содержит турниры)
	// Ставим его первым для приоритетной обработки
	domains = append([]string{mainURL}, domains...)

	return domains, nil
}

// extractDomain извлекает чистый домен из URL (схема + хост, без пути)
func extractDomain(urlStr string) string {
	// Если URL не содержит схему, добавляем https://
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}
