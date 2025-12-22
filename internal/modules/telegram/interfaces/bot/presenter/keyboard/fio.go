package keyboard

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// FioMenu создает клавиатуру для панели ФИО
func (p *KeyboardPresenter) FioMenu(fio entities.TempFIOData) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Очистить фамилию", cb.Filter(cb.FilterFio, cb.FioClearLast)),
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести фамилию", cb.Filter(cb.FilterFio, cb.FioLastName)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Очистить имя", cb.Filter(cb.FilterFio, cb.FioClearFirst)),
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести имя", cb.Filter(cb.FilterFio, cb.FioFirstName)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Очистить отчество", cb.Filter(cb.FilterFio, cb.FioClearPatr)),
			tgbotapi.NewInlineKeyboardButtonData("📝 Ввести отчество", cb.Filter(cb.FilterFio, cb.FioPatronymic)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("◀️ Назад", cb.Filter(cb.FilterFio, cb.FioBack)),
			tgbotapi.NewInlineKeyboardButtonData("✅ Применить", cb.Filter(cb.FilterFio, cb.FioApply)),
		),
	)
}

// FioCancelButton создает кнопку отмены при вводе
func (p *KeyboardPresenter) FioCancelButton() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", cb.Filter(cb.FilterFio, cb.FioSelect)),
		),
	)
}
