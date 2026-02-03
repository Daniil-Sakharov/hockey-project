package game

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseScore(doc *goquery.Document, d *GameDetailsDTO) {
	// Общий счёт из .match-score
	doc.Find(".match-score").First().Each(func(i int, s *goquery.Selection) {
		scores := s.Find(".text-score")
		if scores.Length() >= 2 {
			if h, err := strconv.Atoi(strings.TrimSpace(scores.Eq(0).Text())); err == nil {
				d.HomeScore = &h
			}
			if a, err := strconv.Atoi(strings.TrimSpace(scores.Eq(1).Text())); err == nil {
				d.AwayScore = &a
			}
		}
	})

	// Определение типа результата (ОТ/ПБ) - в отдельном span.text рядом с .match-score
	d.ResultType = "regular"
	doc.Find(".score-wrap span.text, .match-score + span.text, span.text").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "ОТ" {
			d.ResultType = "OT"
		} else if text == "ПБ" || text == "Б" {
			d.ResultType = "SO"
		}
	})

	// Счёт по периодам из .time-score-desk .period-score
	doc.Find(".time-score-desk").First().Each(func(i int, container *goquery.Selection) {
		container.Find(".period-score").Each(func(idx int, ps *goquery.Selection) {
			scores := ps.Find(".text-score")
			if scores.Length() < 2 {
				return
			}

			h, _ := strconv.Atoi(strings.TrimSpace(scores.Eq(0).Text()))
			a, _ := strconv.Atoi(strings.TrimSpace(scores.Eq(1).Text()))

			switch idx {
			case 0:
				d.HomeScoreP1, d.AwayScoreP1 = &h, &a
			case 1:
				d.HomeScoreP2, d.AwayScoreP2 = &h, &a
			case 2:
				d.HomeScoreP3, d.AwayScoreP3 = &h, &a
			case 3:
				d.HomeScoreOT, d.AwayScoreOT = &h, &a
			}
		})
	})
}

func (p *Parser) parseVideoURL(doc *goquery.Document) string {
	video := doc.Find("a[href*='youtube'], a[href*='vk.com/video'], iframe[src*='youtube']")
	if video.Length() > 0 {
		if href, exists := video.First().Attr("href"); exists {
			return href
		}
		if src, exists := video.First().Attr("src"); exists {
			return src
		}
	}
	return ""
}

func (p *Parser) parseTeams(doc *goquery.Document, d *GameDetailsDTO) {
	// Ищем ссылки на команды в разных форматах:
	// - /teams/team_name_ID/
	// - /tournaments/tournament_name/team_name_ID/
	selectors := []string{
		".team-info a[href*='/tournaments/']",
		".match-team a[href*='/tournaments/']",
		".team-info a[href*='/teams/']",
		".match-team a",
	}

	var homeURL, awayURL string

	for _, selector := range selectors {
		teams := doc.Find(selector)
		if teams.Length() >= 2 {
			homeURL, _ = teams.Eq(0).Attr("href")
			awayURL, _ = teams.Eq(1).Attr("href")
			if homeURL != "" && awayURL != "" {
				break
			}
		}
	}

	d.HomeTeamURL = homeURL
	d.AwayTeamURL = awayURL
}
