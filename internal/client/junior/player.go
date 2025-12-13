package junior

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParsePlayersFromTeam парсит игроков из команды
func (c *Client) ParsePlayersFromTeam(teamURL string) ([]PlayerDTO, error) {
	// teamURL уже содержит полный путь типа /tournaments/.../team_id/
	// Формируем полный URL
	fullURL := c.baseURL + teamURL

	resp, err := c.makeRequest(fullURL)
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

	var players []PlayerDTO

	// Ищем таблицу с игроками: table.team-table
	// Пропускаем таблицу тренеров (в ней нет номеров и позиций)
	doc.Find("table.team-table").Each(func(tableIndex int, table *goquery.Selection) {
		// Проверяем что это таблица игроков (по наличию колонки "Амплуа")
		isPlayerTable := false
		table.Find("thead th").Each(func(i int, th *goquery.Selection) {
			headerText := strings.TrimSpace(th.Text())
			if strings.Contains(headerText, "Амплуа") {
				isPlayerTable = true
			}
		})

		if !isPlayerTable {
			return
		}

		// Парсим игроков из tbody
		table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
			// Извлекаем все колонки <td>
			columns := row.Find("td")

			if columns.Length() < 6 {
				// Недостаточно колонок - пропускаем
				return
			}

			player := PlayerDTO{}

			// Колонка 1: Номер + ФИО
			firstCol := columns.Eq(0)
			player.Number = strings.TrimSpace(firstCol.Find("span.number").Text())

			// Ссылка на профиль + ФИО
			profileLink := firstCol.Find("a[href^='/player/']")
			if profileLink.Length() > 0 {
				player.ProfileURL, _ = profileLink.Attr("href")
				player.Name = strings.TrimSpace(profileLink.Text())
			}

			// Колонка 2: Дата рождения
			player.BirthDate = strings.TrimSpace(columns.Eq(1).Find("span.year").Text())

			// Колонка 3: Позиция (Амплуа)
			player.Position = strings.TrimSpace(columns.Eq(2).Find("div.cell").Text())

			// Колонка 4: Рост
			player.Height = strings.TrimSpace(columns.Eq(3).Find("div.cell").Text())

			// Колонка 5: Вес
			player.Weight = strings.TrimSpace(columns.Eq(4).Find("div.cell").Text())

			// Колонка 6: Хват
			player.Handedness = strings.TrimSpace(columns.Eq(5).Find("div.cell").Text())

			// Добавляем только если есть URL профиля (основной идентификатор)
			if player.ProfileURL != "" {
				players = append(players, player)
			}
		})
	})

	return players, nil
}
