package player

import (
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/types"
	"github.com/PuerkitoBio/goquery"
)

// Parser парсер игроков
type Parser struct {
	http    types.HTTPRequester
	baseURL string
}

// NewParser создает новый парсер игроков
func NewParser(http types.HTTPRequester, baseURL string) *Parser {
	return &Parser{http: http, baseURL: baseURL}
}

// ParseFromTeam парсит игроков из команды
func (p *Parser) ParseFromTeam(teamURL string) ([]types.PlayerDTO, error) {
	fullURL := p.baseURL + teamURL

	resp, err := p.http.MakeRequest(fullURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки страницы команды: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP статус %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	var players []types.PlayerDTO

	doc.Find("table.team-table").Each(func(tableIndex int, table *goquery.Selection) {
		isPlayerTable := false
		table.Find("thead th").Each(func(i int, th *goquery.Selection) {
			if strings.Contains(strings.TrimSpace(th.Text()), "Амплуа") {
				isPlayerTable = true
			}
		})

		if !isPlayerTable {
			return
		}

		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			columns := row.Find("td")
			if columns.Length() < 6 {
				return
			}

			player := types.PlayerDTO{}
			firstCol := columns.Eq(0)
			player.Number = strings.TrimSpace(firstCol.Find("span.number").Text())

			profileLink := firstCol.Find("a[href^='/player/']")
			if profileLink.Length() > 0 {
				player.ProfileURL, _ = profileLink.Attr("href")
				player.Name = strings.TrimSpace(profileLink.Text())
			}

			player.BirthDate = strings.TrimSpace(columns.Eq(1).Find("span.year").Text())
			player.Position = strings.TrimSpace(columns.Eq(2).Find("div.cell").Text())
			player.Height = strings.TrimSpace(columns.Eq(3).Find("div.cell").Text())
			player.Weight = strings.TrimSpace(columns.Eq(4).Find("div.cell").Text())
			player.Handedness = strings.TrimSpace(columns.Eq(5).Find("div.cell").Text())

			if player.ProfileURL != "" {
				players = append(players, player)
			}
		})
	})

	return players, nil
}
