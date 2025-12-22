package keyboard

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SearchResultsKeyboard создает клавиатуру результатов поиска
func (p *KeyboardPresenter) SearchResultsKeyboard(players []*services.PlayerWithTeam, currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Кнопки профилей игроков
	for _, player := range players {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("📋 "+player.Name, cb.Player(cb.PlayerProfile, player.ID)),
		})
	}

	// Пагинация
	var navRow []tgbotapi.InlineKeyboardButton
	if currentPage > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", cb.SearchPageDirection(cb.PagePrev)))
	}
	if currentPage < totalPages {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("Вперед ▶️", cb.SearchPageDirection(cb.PageNext)))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	// Кнопка изменения фильтров
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔄 Изменить фильтры", cb.Search(cb.SearchBackToFilters)),
	})

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// ProfileKeyboard создает клавиатуру профиля игрока
func (p *KeyboardPresenter) ProfileKeyboard(playerID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 Скачать отчёт", cb.Report(playerID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ К результатам", cb.Search(cb.SearchBackToResults)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", cb.Menu(cb.MenuMain)),
		),
	)
}
