package keyboard

import (
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MainMenu создает главное меню
func (p *KeyboardPresenter) MainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Поиск игрока", cb.Menu(cb.MenuSearch)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Поиск по статистике 🚧", cb.Menu(cb.MenuStats)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏒 Поиск команды 🚧", cb.Menu(cb.MenuTeam)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Помощь", cb.Menu(cb.MenuHelp)),
		),
	)
}
