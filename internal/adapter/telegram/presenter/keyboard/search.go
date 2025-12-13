package keyboard

import (
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SearchPagination создает клавиатуру с пагинацией результатов поиска
func (p *KeyboardPresenter) SearchPagination(currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}

	// Кнопки пагинации
	paginationRow := []tgbotapi.InlineKeyboardButton{}

	if currentPage > 1 {
		paginationRow = append(paginationRow,
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", cb.SearchPageDirection(cb.PagePrev)))
	}

	if currentPage < totalPages {
		paginationRow = append(paginationRow,
			tgbotapi.NewInlineKeyboardButtonData("Вперед ▶️", cb.SearchPageDirection(cb.PageNext)))
	}

	if len(paginationRow) > 0 {
		rows = append(rows, paginationRow)
	}

	// Кнопка изменения фильтров
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔄 Изменить фильтры", cb.Search(cb.SearchBackToFilters)),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// PlayerProfile создает клавиатуру для карточки игрока
func (p *KeyboardPresenter) PlayerProfile(playerID string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Профиль", cb.Player(cb.PlayerProfile, playerID)),
		),
	)
}
