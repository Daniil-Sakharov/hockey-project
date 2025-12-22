package keyboard

import (
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MainMenu создает главное меню
func (p *KeyboardPresenter) MainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Начать поиск", cb.Menu(cb.MenuSearch)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Фильтры", cb.Menu(cb.MenuStats)),
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", cb.Menu(cb.MenuHelp)),
		),
	)
}
