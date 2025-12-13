package message

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot"
)

// RenderSearchPagination рендерит сообщение пагинации
func (p *MessagePresenter) RenderSearchPagination(result *bot.SearchResult) string {
	if result.TotalCount == 0 {
		return "🔍 Игроки не найдены.\n\nПопробуйте изменить фильтры."
	}

	return fmt.Sprintf("━━━━━━━━━━━━━━━━\nСтраница %d из %d | Всего: %d игроков",
		result.CurrentPage, result.TotalPages, result.TotalCount)
}

// RenderPlayerCard рендерит карточку одного игрока (для отдельного сообщения)
func RenderPlayerCard(p *player.PlayerWithTeam, emoji string) string {
	year := p.Player.BirthDate.Year()

	// Первая строка: смайлик + полное ФИО
	line1 := fmt.Sprintf("%s %s", emoji, p.Player.Name)

	// Вторая строка: Год | Позиция | Рост, Вес | Команда
	height := "?"
	if p.Player.Height != nil {
		height = fmt.Sprintf("%d", *p.Player.Height)
	}

	weight := "?"
	if p.Player.Weight != nil {
		weight = fmt.Sprintf("%d", *p.Player.Weight)
	}

	line2 := fmt.Sprintf("%d | %s | %sсм, %sкг | %s",
		year, p.Player.Position, height, weight, p.TeamName)

	return fmt.Sprintf("%s\n%s", line1, line2)
}
