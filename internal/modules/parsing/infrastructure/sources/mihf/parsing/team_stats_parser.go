package parsing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/PuerkitoBio/goquery"
)

var playerIDRegex = regexp.MustCompile(`/players/info/(\d+)`)

// tableType определяет тип таблицы по заголовку
type tableType int

const (
	tableTypeUnknown tableType = iota
	tableTypeGoalies           // Вратари
	tableTypeDefenders         // Защитники
	tableTypeForwards          // Нападающие
)

// ParseTeamStats парсит статистику игроков и вратарей команды
// Различает 3 таблицы: Вратари, Защитники, Нападающие
func ParseTeamStats(html []byte) ([]dto.PlayerStatsDTO, []dto.GoalieStatsDTO, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, nil, err
	}

	var players []dto.PlayerStatsDTO
	var goalies []dto.GoalieStatsDTO

	// Находим все таблицы со статистикой
	doc.Find("table.table-hover").Each(func(tableIdx int, table *goquery.Selection) {
		// Определяем тип таблицы по заголовку перед ней
		tblType := detectTableType(table)

		table.Find("tr").Each(func(i int, row *goquery.Selection) {
			// Пропускаем заголовки
			if row.Find("th").Length() > 0 {
				return
			}

			// Пропускаем строки-разделители (class="title")
			class, _ := row.Attr("class")
			if class == "title" {
				return
			}

			// Проверяем есть ли ссылка на игрока
			if row.Find("a[href*='/players/info/']").Length() == 0 {
				return
			}

			// Парсим в зависимости от типа таблицы
			switch tblType {
			case tableTypeGoalies:
				goalie := parseGoalieRow(row)
				if goalie != nil {
					goalies = append(goalies, *goalie)
				}
			case tableTypeDefenders:
				player := parsePlayerRow(row, "З")
				if player != nil {
					players = append(players, *player)
				}
			case tableTypeForwards:
				player := parsePlayerRow(row, "Н")
				if player != nil {
					players = append(players, *player)
				}
			default:
				// Неизвестный тип - парсим как полевого без позиции
				player := parsePlayerRow(row, "")
				if player != nil {
					players = append(players, *player)
				}
			}
		})
	})

	return players, goalies, nil
}

// detectTableType определяет тип таблицы по заголовку
func detectTableType(table *goquery.Selection) tableType {
	// Способ 1: структура stats.mihf.ru - <h3>Заголовок</h3><div class="table-responsive"><table>
	// table.Parent() = div, div.Prev() = h3
	parent := table.Parent()
	if parent.Length() > 0 {
		prev := parent.Prev()
		if prev.Length() > 0 {
			tagName := goquery.NodeName(prev)
			if tagName == "h3" || tagName == "h4" || tagName == "h5" {
				text := strings.ToLower(strings.TrimSpace(prev.Text()))
				if tblType := parseTableTypeFromText(text); tblType != tableTypeUnknown {
					return tblType
				}
			}
		}
	}

	// Способ 2: ищем h3 непосредственно перед таблицей (на случай другой структуры)
	prev := table.Prev()
	if prev.Length() > 0 {
		tagName := goquery.NodeName(prev)
		if tagName == "h3" || tagName == "h4" || tagName == "h5" {
			text := strings.ToLower(strings.TrimSpace(prev.Text()))
			if tblType := parseTableTypeFromText(text); tblType != tableTypeUnknown {
				return tblType
			}
		}
	}

	// Способ 3: проверяем первую строку таблицы с class="title"
	titleRow := table.Find("tr.title").First()
	if titleRow.Length() > 0 {
		titleText := strings.ToLower(titleRow.Text())
		if tblType := parseTableTypeFromText(titleText); tblType != tableTypeUnknown {
			return tblType
		}
		// Не нашли тип - продолжаем к способу 4
	}

	// Способ 4: анализируем колонки таблицы
	// Вратари имеют колонку "КН" (коэффициент надёжности)
	headerCells := table.Find("tr").First().Find("th")
	hasKN := false
	headerCells.Each(func(_ int, th *goquery.Selection) {
		text := strings.TrimSpace(th.Text())
		if text == "Кн" || text == "КН" || text == "кн" {
			hasKN = true
		}
	})
	if hasKN {
		return tableTypeGoalies
	}

	return tableTypeUnknown
}

// parseTableTypeFromText определяет тип таблицы по тексту
func parseTableTypeFromText(text string) tableType {
	text = strings.ToLower(text)
	if strings.Contains(text, "вратар") {
		return tableTypeGoalies
	}
	if strings.Contains(text, "защитник") {
		return tableTypeDefenders
	}
	if strings.Contains(text, "нападающ") {
		return tableTypeForwards
	}
	return tableTypeUnknown
}

// parsePlayerRow парсит строку полевого игрока (защитника или нападающего)
// Структура колонок: №, ФИО, И, Г, А, О, Ш, ГБ, ГМ, ГР
func parsePlayerRow(row *goquery.Selection, position string) *dto.PlayerStatsDTO {
	cells := row.Find("td")
	if cells.Length() < 5 {
		return nil
	}

	player := &dto.PlayerStatsDTO{
		Position: position,
	}

	// Ссылка на профиль
	playerLink := row.Find("a[href*='/players/info/']")
	if href, exists := playerLink.Attr("href"); exists {
		if matches := playerIDRegex.FindStringSubmatch(href); len(matches) > 1 {
			player.ID = matches[1]
			player.ProfileURL = href
		}
	}

	player.Name = strings.TrimSpace(playerLink.Text())

	// Структура столбцов для полевых игроков:
	// 0: №, 1: ФИО, 2: И, 3: Г, 4: А, 5: О, 6: Ш, 7: ГБ, 8: ГМ, 9: ГР
	cells.Each(func(j int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())
		val, _ := strconv.Atoi(text)

		switch j {
		case 0: // № - номер
			player.Number = text
		case 2: // И - игры
			player.Games = val
		case 3: // Г - голы
			player.Goals = val
		case 4: // А - передачи
			player.Assists = val
		case 5: // О - очки
			player.Points = val
		case 6: // Ш - штрафные минуты
			player.PenaltyMinutes = val
		case 7: // ГБ - голы в большинстве
			player.GoalsPowerPlay = val
		case 8: // ГМ - голы в меньшинстве
			player.GoalsShortHanded = val
		case 9: // ГР - голы в равных составах
			player.GoalsEvenStrength = val
		}
	})

	if player.ID == "" {
		return nil
	}
	return player
}

// parseGoalieRow парсит строку вратаря
// Структура колонок: №, ФИО, И, Г, А, О, КН, ШП, Мин:сек
func parseGoalieRow(row *goquery.Selection) *dto.GoalieStatsDTO {
	cells := row.Find("td")
	if cells.Length() < 6 {
		return nil
	}

	goalie := &dto.GoalieStatsDTO{}

	// Ссылка на профиль
	playerLink := row.Find("a[href*='/players/info/']")
	if href, exists := playerLink.Attr("href"); exists {
		if matches := playerIDRegex.FindStringSubmatch(href); len(matches) > 1 {
			goalie.ID = matches[1]
			goalie.ProfileURL = href
		}
	}

	goalie.Name = strings.TrimSpace(playerLink.Text())

	// Структура столбцов для вратарей:
	// 0: №, 1: ФИО, 2: И, 3: Г, 4: А, 5: О, 6: КН, 7: ШП, 8: Мин:сек
	cells.Each(func(j int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())
		val, _ := strconv.Atoi(text)

		switch j {
		case 0: // № - номер
			goalie.Number = text
		case 2: // И - игры
			goalie.Games = val
		case 3: // Г - голы (забитые вратарём)
			goalie.Goals = val
		case 4: // А - передачи
			goalie.Assists = val
		case 5: // О - очки
			goalie.Points = val
		case 6: // КН - коэффициент надёжности (save percentage)
			goalie.SavePercentage = parseFloatValue(text)
		case 7: // ШП - пропущено шайб (goals against)
			goalie.GoalsAgainst = val
		case 8: // Мин:сек - время на льду
			goalie.MinutesPlayed = parseMinutes(text)
		}
	})

	if goalie.ID == "" {
		return nil
	}
	return goalie
}

func parseFloatValue(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	s = strings.TrimSuffix(s, "%")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseMinutes(s string) int {
	s = strings.TrimSpace(s)
	// Формат может быть "45:30" (минуты:секунды)
	parts := strings.Split(s, ":")
	if len(parts) >= 1 {
		mins, _ := strconv.Atoi(parts[0])
		return mins
	}
	return 0
}
