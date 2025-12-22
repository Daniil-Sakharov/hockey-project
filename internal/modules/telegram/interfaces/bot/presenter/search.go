package presenter

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
)

// RenderSearchResults рендерит результаты поиска
func (p *Presenter) RenderSearchResults(result *services.SearchResult) string {
	if result.TotalCount == 0 {
		return "🔍 Игроки не найдены.\n\nПопробуйте изменить фильтры."
	}

	emojis := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}
	text := fmt.Sprintf("🔍 **Найдено: %d игроков**\n\n", result.TotalCount)

	for i, player := range result.Players {
		emoji := "▪️"
		if i < len(emojis) {
			emoji = emojis[i]
		}

		text += fmt.Sprintf("%s %s\n", emoji, player.Name)

		height := "?"
		if player.Height > 0 {
			height = fmt.Sprintf("%d", player.Height)
		}
		weight := "?"
		if player.Weight > 0 {
			weight = fmt.Sprintf("%d", player.Weight)
		}

		text += fmt.Sprintf("%s | %s | %sсм, %sкг | %s\n\n",
			player.BirthDate, player.Position, height, weight, player.TeamName)
	}

	text += "━━━━━━━━━━━━━━━━━━━━━━\n"
	text += fmt.Sprintf("Страница %d из %d | Всего: %d игроков",
		result.CurrentPage, result.TotalPages, result.TotalCount)

	return text
}
