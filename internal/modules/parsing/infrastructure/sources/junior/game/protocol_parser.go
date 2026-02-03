package game

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parseProtocol парсит полный протокол матча из секции .match-protocol
func (p *Parser) parseProtocol(doc *goquery.Document, details *GameDetailsDTO) {
	doc.Find(".match-protocol").Each(func(i int, s *goquery.Selection) {
		p.parseProtocolEvent(s, details)
	})
}

// parseProtocolEvent парсит одно событие протокола
func (p *Parser) parseProtocolEvent(s *goquery.Selection, details *GameDetailsDTO) {
	center := s.Find(".match-protocol__center")
	eventType := strings.ToLower(strings.TrimSpace(center.Find(".protocol-time__top").First().Text()))
	timeText := strings.TrimSpace(center.Find(".protocol-time__bottom span").Last().Text())
	minutes, seconds := parseProtocolTime(timeText)

	// Определяем тип события
	switch {
	case strings.Contains(eventType, "вратарь"):
		p.parseGoalieEvent(s, details, minutes, seconds)
	case strings.Contains(eventType, "штраф"):
		p.parsePenaltyEvent(s, details, minutes, seconds)
	case strings.Contains(eventType, "тайм-аут"):
		p.parseTimeoutEvent(s, details, minutes, seconds)
	case isScoreEvent(eventType):
		p.parseGoalEvent(s, details, minutes, seconds, eventType)
	}
}

// isScoreEvent проверяет, является ли событие голом (формат "X - Y")
func isScoreEvent(text string) bool {
	matched, _ := regexp.MatchString(`\d+\s*-\s*\d+`, text)
	return matched
}

// parseProtocolTime парсит время из формата "MM:SS"
func parseProtocolTime(timeText string) (int, int) {
	timeText = strings.TrimSpace(timeText)
	parts := strings.Split(timeText, ":")
	if len(parts) != 2 {
		return 0, 0
	}
	minutes, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	seconds, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	return minutes, seconds
}

// CalculatePeriod вычисляет период по времени матча
func CalculatePeriod(minutes int) int {
	switch {
	case minutes < 20:
		return 1
	case minutes < 40:
		return 2
	case minutes < 60:
		return 3
	default:
		return 4 // OT
	}
}

// parseGoalEvent парсит событие гола из протокола
func (p *Parser) parseGoalEvent(s *goquery.Selection, details *GameDetailsDTO, minutes, seconds int, scoreText string) {
	goal := GoalDTO{
		TimeMinutes: minutes,
		TimeSeconds: seconds,
		Period:      CalculatePeriod(minutes),
	}

	// Определяем какая команда забила по highlighted классу
	center := s.Find(".match-protocol__center")
	homeHighlighted := center.Find(".protocol-time__top span.highlighted").First()
	isHomeGoal := homeHighlighted.Length() > 0 && homeHighlighted.Index() == 0

	// Выбираем нужную сторону
	var playerCard *goquery.Selection
	if isHomeGoal {
		playerCard = s.Find(".match-protocol__left .protocol-player-card").First()
		goal.GoalType = "home"
	} else {
		playerCard = s.Find(".match-protocol__right .protocol-player-card").First()
		goal.GoalType = "away"
	}

	if playerCard.Length() == 0 {
		return
	}

	// Проверяем, не пустые ли это ворота
	cardName := strings.TrimSpace(playerCard.Find(".protocol-player-card__name").Text())
	if strings.Contains(strings.ToLower(cardName), "пустые ворота") {
		// Это гол в пустые ворота, игрок на другой стороне
		if isHomeGoal {
			playerCard = s.Find(".match-protocol__right .protocol-player-card").First()
		} else {
			playerCard = s.Find(".match-protocol__left .protocol-player-card").First()
		}
		goal.GoalType = "en" // Empty net
	}

	// Извлекаем автора гола
	scorerLink := playerCard.Find("a.protocol-player-card__name")
	if scorerLink.Length() > 0 {
		goal.ScorerURL, _ = scorerLink.Attr("href")
		goal.ScorerName = strings.TrimSpace(scorerLink.Text())
	}

	// Извлекаем ассистентов из .protocol-player-card__label span
	playerCard.Find(".protocol-player-card__label span").Each(func(i int, span *goquery.Selection) {
		assistText := strings.TrimSpace(span.Text())
		if assistText == "" {
			return
		}
		// Формат: "55. Мухтаров Ленар"
		name := extractPlayerNameFromText(assistText)
		number := extractJerseyNumber(assistText)
		if i == 0 {
			goal.Assist1Name = name
			goal.Assist1Number = number
		} else if i == 1 {
			goal.Assist2Name = name
			goal.Assist2Number = number
		}
	})

	// Извлекаем "Кто был на льду"
	p.parsePlayersOnIceFromAccord(playerCard, &goal, isHomeGoal)

	details.Goals = append(details.Goals, goal)
}

// extractPlayerNameFromText извлекает имя игрока из текста формата "55. Мухтаров Ленар"
func extractPlayerNameFromText(text string) string {
	// Убираем номер в начале
	re := regexp.MustCompile(`^\d+\.\s*`)
	return re.ReplaceAllString(text, "")
}

// ExtractJerseyNumber извлекает номер игрока из текста формата "55. Мухтаров Ленар"
func ExtractJerseyNumber(text string) int {
	re := regexp.MustCompile(`^(\d+)\.`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		num, _ := strconv.Atoi(matches[1])
		return num
	}
	return 0
}

// extractJerseyNumber - внутренний alias для обратной совместимости
func extractJerseyNumber(text string) int {
	return ExtractJerseyNumber(text)
}

// parsePlayersOnIceFromAccord парсит игроков на льду из аккордеона
func (p *Parser) parsePlayersOnIceFromAccord(card *goquery.Selection, goal *GoalDTO, isHomeGoal bool) {
	accordBody := card.Find(".accord-body.js-accord-body")
	if accordBody.Length() == 0 {
		return
	}

	// Две ul - первая для команды автора гола, вторая для соперника
	accordBody.Find("ul").Each(func(i int, ul *goquery.Selection) {
		var players []string
		ul.Find("li").Each(func(j int, li *goquery.Selection) {
			// Пропускаем название команды (первый li)
			if j == 0 {
				return
			}
			// Пробуем найти ссылку на игрока - в HTML есть <a href="/player/...">
			link := li.Find("a")
			if link.Length() > 0 {
				if href, exists := link.Attr("href"); exists && href != "" {
					players = append(players, href)
					return
				}
			}
			// Fallback: текст (для обратной совместимости)
			text := strings.TrimSpace(li.Text())
			if text != "" {
				players = append(players, text)
			}
		})

		// Распределяем по командам
		if i == 0 {
			// Первая ul - команда автора гола
			if isHomeGoal {
				goal.HomePlayersOnIce = players
			} else {
				goal.AwayPlayersOnIce = players
			}
		} else if i == 1 {
			// Вторая ul - команда соперника
			if isHomeGoal {
				goal.AwayPlayersOnIce = players
			} else {
				goal.HomePlayersOnIce = players
			}
		}
	})
}

// parseGoalieEvent парсит событие вратаря
func (p *Parser) parseGoalieEvent(s *goquery.Selection, details *GameDetailsDTO, minutes, seconds int) {
	// Проверяем левую и правую сторону
	p.parseGoalieSide(s.Find(".match-protocol__left"), details, minutes, seconds, true)
	p.parseGoalieSide(s.Find(".match-protocol__right"), details, minutes, seconds, false)
}

func (p *Parser) parseGoalieSide(side *goquery.Selection, details *GameDetailsDTO, minutes, seconds int, isHome bool) {
	if side.HasClass("empty") {
		return
	}

	card := side.Find(".protocol-player-card").First()
	if card.Length() == 0 {
		return
	}

	cardName := strings.TrimSpace(card.Find(".protocol-player-card__name").Text())

	// Проверяем на пустые ворота
	if strings.Contains(strings.ToLower(cardName), "пустые ворота") {
		emptyNet := EmptyNetDTO{
			TimeMinutes: minutes,
			TimeSeconds: seconds,
			IsHome:      isHome,
		}
		details.EmptyNets = append(details.EmptyNets, emptyNet)
		return
	}

	// Это событие вратаря
	goalieEvent := GoalieEventDTO{
		TimeMinutes: minutes,
		TimeSeconds: seconds,
		IsHome:      isHome,
	}

	// Извлекаем URL и имя вратаря
	goalieLink := card.Find("a.protocol-player-card__name")
	if goalieLink.Length() > 0 {
		goalieEvent.PlayerURL, _ = goalieLink.Attr("href")
		goalieEvent.PlayerName = strings.TrimSpace(goalieLink.Text())
	}

	details.GoalieEvents = append(details.GoalieEvents, goalieEvent)
}

// parsePenaltyEvent парсит событие штрафа
func (p *Parser) parsePenaltyEvent(s *goquery.Selection, details *GameDetailsDTO, minutes, seconds int) {
	// Проверяем левую и правую сторону
	p.parsePenaltySide(s.Find(".match-protocol__left"), details, minutes, seconds, true)
	p.parsePenaltySide(s.Find(".match-protocol__right"), details, minutes, seconds, false)
}

func (p *Parser) parsePenaltySide(side *goquery.Selection, details *GameDetailsDTO, minutes, seconds int, isHome bool) {
	if side.HasClass("empty") {
		return
	}

	card := side.Find(".protocol-player-card").First()
	if card.Length() == 0 {
		return
	}

	penalty := PenaltyDTO{
		TimeMinutes: minutes,
		TimeSeconds: seconds,
		Period:      CalculatePeriod(minutes),
		IsHome:      isHome,
	}

	// Извлекаем URL и имя игрока
	playerLink := card.Find("a.protocol-player-card__name")
	if playerLink.Length() > 0 {
		penalty.PlayerURL, _ = playerLink.Attr("href")
		penalty.PlayerName = strings.TrimSpace(playerLink.Text())
	}

	// Извлекаем причину штрафа
	penalty.Reason = strings.TrimSpace(card.Find(".protocol-player-card__label").Text())

	// Извлекаем минуты штрафа из sidebar (формат "2'" или "ШБ'")
	sidebarText := strings.TrimSpace(card.Find(".protocol-player-card__sidebar").Text())
	penalty.Minutes = extractPenaltyMinutes(sidebarText)

	details.Penalties = append(details.Penalties, penalty)
}

// extractPenaltyMinutes извлекает минуты штрафа из текста
func extractPenaltyMinutes(text string) int {
	text = strings.ToLower(text)
	// Удаляем апостроф
	text = strings.ReplaceAll(text, "'", "")

	// Особые случаи
	if strings.Contains(text, "шб") {
		return 5 // Штрафной бросок
	}

	// Пробуем извлечь число
	re := regexp.MustCompile(`\d+`)
	if match := re.FindString(text); match != "" {
		mins, _ := strconv.Atoi(match)
		return mins
	}
	return 2 // По умолчанию
}

// parseTimeoutEvent парсит событие тайм-аута
func (p *Parser) parseTimeoutEvent(s *goquery.Selection, details *GameDetailsDTO, minutes, seconds int) {
	// Проверяем левую сторону
	left := s.Find(".match-protocol__left")
	if !left.HasClass("empty") {
		timeout := TimeoutDTO{
			TimeMinutes: minutes,
			TimeSeconds: seconds,
			IsHome:      true,
		}
		details.Timeouts = append(details.Timeouts, timeout)
	}

	// Проверяем правую сторону
	right := s.Find(".match-protocol__right")
	if !right.HasClass("empty") {
		timeout := TimeoutDTO{
			TimeMinutes: minutes,
			TimeSeconds: seconds,
			IsHome:      false,
		}
		details.Timeouts = append(details.Timeouts, timeout)
	}
}
