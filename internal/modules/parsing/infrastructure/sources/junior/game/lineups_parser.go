package game

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parseLineups(doc *goquery.Document) (home, away []PlayerLineup) {
	// Структура junior.fhr.ru: .team-lineups с .left-team и .right-team
	lineups := doc.Find(".team-lineups")
	if lineups.Length() > 0 {
		leftTeam := lineups.Find(".left-team")
		rightTeam := lineups.Find(".right-team")

		if leftTeam.Length() > 0 {
			home = p.parseTeamBlock(leftTeam)
		}
		if rightTeam.Length() > 0 {
			away = p.parseTeamBlock(rightTeam)
		}
	}

	return home, away
}

func (p *Parser) parseTeamBlock(block *goquery.Selection) []PlayerLineup {
	var lineup []PlayerLineup
	seen := make(map[string]bool)

	// Структура: одна таблица с несколькими tbody по позициям
	table := block.Find(".team-table")
	if table.Length() == 0 {
		return lineup
	}

	// Определяем порядок колонок из заголовка
	colMap := p.parseTableHeader(table)

	// Каждый tbody = отдельная позиция
	table.Find("tbody").Each(func(i int, tbody *goquery.Selection) {
		// Пропускаем персонал (staff-body)
		if tbody.HasClass("staff-body") {
			return
		}

		// Определяем позицию по data-id:
		// 1 = Вратари, 2 = Защитники, 3 = Нападающие
		dataID, _ := tbody.Attr("data-id")
		position := "F"
		switch dataID {
		case "1":
			position = "G"
		case "2":
			position = "D"
		case "3":
			position = "F"
		}

		// Парсим игроков из этого tbody
		tbody.Find("tr").Each(func(j int, row *goquery.Selection) {
			if player := p.parsePlayerTableRow(row, position, colMap); player != nil {
				if !seen[player.PlayerURL] {
					seen[player.PlayerURL] = true
					lineup = append(lineup, *player)
				}
			}
		})
	})

	return lineup
}

// parseTableHeader определяет порядок колонок статистики
func (p *Parser) parseTableHeader(table *goquery.Selection) map[string]int {
	colMap := make(map[string]int)

	table.Find("thead th, thead td").Each(func(i int, th *goquery.Selection) {
		text := strings.ToLower(strings.TrimSpace(th.Text()))

		// Маппинг названий колонок
		switch {
		case text == "ш" || text == "г" || strings.Contains(text, "гол"):
			colMap["goals"] = i
		case text == "п" || text == "а" || strings.Contains(text, "передач") || strings.Contains(text, "ассист"):
			colMap["assists"] = i
		case text == "о" || text == "очк" || strings.Contains(text, "очк"):
			colMap["points"] = i
		case text == "штр" || strings.Contains(text, "штраф"):
			colMap["penalty"] = i
		case text == "+/-" || text == "+-":
			colMap["plusminus"] = i
		case text == "св" || strings.Contains(text, "спас"):
			colMap["saves"] = i
		case text == "пр" || text == "пш" || strings.Contains(text, "пропущ"):
			colMap["goalsagainst"] = i
		case text == "вр" || strings.Contains(text, "время"):
			colMap["toi"] = i
		}
	})

	return colMap
}

func (p *Parser) parsePlayerTableRow(row *goquery.Selection, position string, colMap map[string]int) *PlayerLineup {
	player := &PlayerLineup{}

	// Ссылка на игрока (.link-tr)
	link := row.Find("a.link-tr, a[href*='/player/']").First()
	if link.Length() > 0 {
		player.PlayerURL, _ = link.Attr("href")
	}

	// Имя игрока из .text
	nameEl := row.Find(".cell.player .text")
	if nameEl.Length() > 0 {
		player.PlayerName = strings.TrimSpace(nameEl.Text())
	}

	// Если нет ссылки — пропускаем
	if player.PlayerURL == "" {
		return nil
	}

	// Номер из .number
	numText := row.Find(".number").Text()
	if n, err := strconv.Atoi(strings.TrimSpace(numText)); err == nil {
		player.JerseyNumber = n
	}

	// Роль: К (капитан) или А (ассистент)
	captainMark := row.Find(".captain-mark")
	if captainMark.Length() > 0 {
		markText := strings.TrimSpace(captainMark.Text())
		if markText == "К" {
			player.Role = "C"
		} else if markText == "А" {
			player.Role = "A"
		}
	}

	// Позиция
	player.Position = position
	if player.Position == "" {
		player.Position = "F"
	}

	// Парсим статистику из ячеек таблицы
	cells := row.Find("td")
	p.parsePlayerStats(player, cells, colMap)

	return player
}

// parsePlayerStats извлекает статистику из ячеек строки
func (p *Parser) parsePlayerStats(player *PlayerLineup, cells *goquery.Selection, colMap map[string]int) {
	// Если карта пустая, используем позиции по умолчанию
	// Типичный порядок: #, Имя, Ш, П, О/ШТР, +/-
	if len(colMap) == 0 {
		if cells.Length() >= 6 {
			// Пытаемся извлечь из фиксированных позиций
			player.Goals = parseInt(cells.Eq(2).Text())
			player.Assists = parseInt(cells.Eq(3).Text())
			player.PenaltyMinutes = parseInt(cells.Eq(4).Text())
			player.PlusMinus = parseInt(cells.Eq(5).Text())
		}
		return
	}

	// Используем карту колонок
	if idx, ok := colMap["goals"]; ok && idx < cells.Length() {
		player.Goals = parseInt(cells.Eq(idx).Text())
	}
	if idx, ok := colMap["assists"]; ok && idx < cells.Length() {
		player.Assists = parseInt(cells.Eq(idx).Text())
	}
	if idx, ok := colMap["penalty"]; ok && idx < cells.Length() {
		player.PenaltyMinutes = parseInt(cells.Eq(idx).Text())
	}
	if idx, ok := colMap["plusminus"]; ok && idx < cells.Length() {
		player.PlusMinus = parseInt(cells.Eq(idx).Text())
	}

	// Для вратарей
	if player.Position == "G" {
		if idx, ok := colMap["saves"]; ok && idx < cells.Length() {
			saves := parseInt(cells.Eq(idx).Text())
			player.Saves = &saves
		}
		if idx, ok := colMap["goalsagainst"]; ok && idx < cells.Length() {
			ga := parseInt(cells.Eq(idx).Text())
			player.GoalsAgainst = &ga
		}
		if idx, ok := colMap["toi"]; ok && idx < cells.Length() {
			toi := parseTimeOnIce(cells.Eq(idx).Text())
			if toi > 0 {
				player.TimeOnIce = &toi
			}
		}
	}
}

// parseInt парсит число из текста
func parseInt(text string) int {
	text = strings.TrimSpace(text)
	// Убираем + перед числом для +/-
	text = strings.TrimPrefix(text, "+")
	if n, err := strconv.Atoi(text); err == nil {
		return n
	}
	return 0
}

// parseTimeOnIce парсит время на льду в секунды (формат: "MM:SS" или "M:SS")
func parseTimeOnIce(text string) int {
	text = strings.TrimSpace(text)
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		return 0
	}

	minutes, _ := strconv.Atoi(parts[0])
	seconds, _ := strconv.Atoi(parts[1])

	return minutes*60 + seconds
}
