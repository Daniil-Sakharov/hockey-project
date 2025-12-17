package keyboard

import (
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/adapter/telegram/callback"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// FioMenu создает клавиатуру для панели ФИО
func (p *KeyboardPresenter) FioMenu(fio bot.TempFioData) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}

	// Ряд 1: Фамилия
	row1 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("❌ Очистить фамилию", cb.Filter(cb.FilterFio, cb.FioClearLast)),
		tgbotapi.NewInlineKeyboardButtonData("📝 Ввести фамилию", cb.Filter(cb.FilterFio, cb.FioLastName)),
	}
	rows = append(rows, row1)

	// Ряд 2: Имя
	row2 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("❌ Очистить имя", cb.Filter(cb.FilterFio, cb.FioClearFirst)),
		tgbotapi.NewInlineKeyboardButtonData("📝 Ввести имя", cb.Filter(cb.FilterFio, cb.FioFirstName)),
	}
	rows = append(rows, row2)

	// Ряд 3: Отчество
	row3 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("❌ Очистить отчество", cb.Filter(cb.FilterFio, cb.FioClearPatr)),
		tgbotapi.NewInlineKeyboardButtonData("📝 Ввести отчество", cb.Filter(cb.FilterFio, cb.FioPatronymic)),
	}
	rows = append(rows, row3)

	// Ряд 4: Назад и Применить
	row4 := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", cb.Filter(cb.FilterFio, cb.FioBack)),
		tgbotapi.NewInlineKeyboardButtonData("✅ Применить", cb.Filter(cb.FilterFio, cb.FioApply)),
	}
	rows = append(rows, row4)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// FioCancelButton создает кнопку отмены при вводе
func (p *KeyboardPresenter) FioCancelButton() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", cb.Filter(cb.FilterFio, cb.FioSelect)),
		),
	)
}
