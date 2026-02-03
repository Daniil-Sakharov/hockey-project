package filters

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер фильтров турнира (года, группы)
type Parser struct {
	http types.HTTPRequester
}

// NewParser создает новый парсер фильтров
func NewParser(http types.HTTPRequester) *Parser {
	return &Parser{http: http}
}

// YearFilter год рождения с AJAX URL
type YearFilter struct {
	Year    int
	YearID  string
	AjaxURL string
}

// GroupFilter группа с AJAX URL
type GroupFilter struct {
	Name    string
	GroupID string
	AjaxURL string
}

// TournamentFilters все фильтры турнира
type TournamentFilters struct {
	Years  []YearFilter
	Groups []GroupFilter
}

// Parse парсит фильтры со страницы турнира
func (p *Parser) Parse(tournamentURL string) (*TournamentFilters, error) {
	resp, err := p.http.MakeRequest(tournamentURL)
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

	filters := &TournamentFilters{}
	filters.Years = p.parseYears(doc)
	filters.Groups = p.parseGroups(doc)

	return filters, nil
}

// ParseFromAjax парсит фильтры из AJAX ответа
func (p *Parser) ParseFromAjax(ajaxURL string) (*TournamentFilters, error) {
	resp, err := p.http.MakeRequestWithHeaders(ajaxURL, map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	})
	if err != nil {
		return nil, fmt.Errorf("ajax request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	filters := &TournamentFilters{}
	filters.Years = p.parseYears(doc)
	filters.Groups = p.parseGroups(doc)

	return filters, nil
}

func (p *Parser) parseYears(doc *goquery.Document) []YearFilter {
	var years []YearFilter

	doc.Find("select[data-ajax-select] option").Each(func(i int, opt *goquery.Selection) {
		yearText := strings.TrimSpace(opt.Text())
		year, err := strconv.Atoi(yearText)
		if err != nil {
			return
		}

		yearID, _ := opt.Attr("value")
		ajaxURL, _ := opt.Attr("data-ajax")

		years = append(years, YearFilter{
			Year:    year,
			YearID:  yearID,
			AjaxURL: ajaxURL,
		})
	})

	return years
}

func (p *Parser) parseGroups(doc *goquery.Document) []GroupFilter {
	var groups []GroupFilter

	doc.Find("div.filter-btn[data-ajax-link]").Each(func(i int, btn *goquery.Selection) {
		name := strings.TrimSpace(btn.Text())
		ajaxURL, _ := btn.Attr("data-ajax-link")

		groupID := extractGroupID(ajaxURL)

		groups = append(groups, GroupFilter{
			Name:    name,
			GroupID: groupID,
			AjaxURL: ajaxURL,
		})
	})

	return groups
}

var groupIDRegex = regexp.MustCompile(`GROUP_ID";s:\d+:"(\d+)"`)

func extractGroupID(ajaxURL string) string {
	// Извлекаем params из URL
	parsed, err := url.Parse(ajaxURL)
	if err != nil {
		return ""
	}

	params := parsed.Query().Get("params")
	if params == "" {
		return ""
	}

	// Декодируем base64
	decoded, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		// Попробуем URL-decoded версию
		unescaped, _ := url.QueryUnescape(params)
		decoded, err = base64.StdEncoding.DecodeString(unescaped)
		if err != nil {
			return ""
		}
	}

	// Ищем GROUP_ID в PHP serialized формате: s:8:"GROUP_ID";s:8:"17175707"
	matches := groupIDRegex.FindStringSubmatch(string(decoded))
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}
