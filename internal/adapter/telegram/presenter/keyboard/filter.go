package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// FilterMenu создает меню фильтров
func (p *KeyboardPresenter) FilterMenu(hasFilters bool) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("📝 Ввести ФИО", "filter:fio:select")},
		{tgbotapi.NewInlineKeyboardButtonData("🎂 Год рождения ▼", "filter:year:select")},
		{tgbotapi.NewInlineKeyboardButtonData("🏒 Позиция ▼", "filter:position:select")},
		{
			tgbotapi.NewInlineKeyboardButtonData("📏 Рост ▼", "filter:height:select"),
			tgbotapi.NewInlineKeyboardButtonData("⚖️ Вес ▼", "filter:weight:select"),
		},
		{tgbotapi.NewInlineKeyboardButtonData("🗺️ Регион ▼", "filter:region:select")},
	}

	// Кнопки действий
	if hasFilters {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🔍 НАЙТИ", "filter:apply"),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Сбросить", "filter:reset"),
		})
	} else {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🔍 НАЙТИ", "filter:apply"),
		})
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "menu:main"),
	})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// YearSelect создает клавиатуру выбора года
func (p *KeyboardPresenter) YearSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2005", "filter:year:2005"),
			tgbotapi.NewInlineKeyboardButtonData("2006", "filter:year:2006"),
			tgbotapi.NewInlineKeyboardButtonData("2007", "filter:year:2007"),
			tgbotapi.NewInlineKeyboardButtonData("2008", "filter:year:2008"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2009", "filter:year:2009"),
			tgbotapi.NewInlineKeyboardButtonData("2010", "filter:year:2010"),
			tgbotapi.NewInlineKeyboardButtonData("2011", "filter:year:2011"),
			tgbotapi.NewInlineKeyboardButtonData("2012", "filter:year:2012"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2013", "filter:year:2013"),
			tgbotapi.NewInlineKeyboardButtonData("2014", "filter:year:2014"),
			tgbotapi.NewInlineKeyboardButtonData("2015", "filter:year:2015"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", "filter:year:any"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", "filter:back"),
		),
	)
}

// PositionSelect создает клавиатуру выбора позиции
func (p *KeyboardPresenter) PositionSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎯 Нападающий", "filter:position:forward"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛡️ Защитник", "filter:position:defender"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥅 Вратарь", "filter:position:goalie"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любая", "filter:position:any"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", "filter:back"),
		),
	)
}

// HeightSelect создает клавиатуру выбора роста
func (p *KeyboardPresenter) HeightSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("150-160", "filter:height:150-160"),
			tgbotapi.NewInlineKeyboardButtonData("160-170", "filter:height:160-170"),
			tgbotapi.NewInlineKeyboardButtonData("170-180", "filter:height:170-180"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("180-190", "filter:height:180-190"),
			tgbotapi.NewInlineKeyboardButtonData("190-200", "filter:height:190-200"),
			tgbotapi.NewInlineKeyboardButtonData("200+", "filter:height:200+"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", "filter:height:any"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", "filter:back"),
		),
	)
}

// WeightSelect создает клавиатуру выбора веса
func (p *KeyboardPresenter) WeightSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("40-50", "filter:weight:40-50"),
			tgbotapi.NewInlineKeyboardButtonData("50-60", "filter:weight:50-60"),
			tgbotapi.NewInlineKeyboardButtonData("60-70", "filter:weight:60-70"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("70-80", "filter:weight:70-80"),
			tgbotapi.NewInlineKeyboardButtonData("80-90", "filter:weight:80-90"),
			tgbotapi.NewInlineKeyboardButtonData("90+", "filter:weight:90+"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", "filter:weight:any"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", "filter:back"),
		),
	)
}

// RegionSelect создает клавиатуру выбора региона
func (p *KeyboardPresenter) RegionSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ФХР", "filter:region:ФХР"),
			tgbotapi.NewInlineKeyboardButtonData("СПБ", "filter:region:СПБ"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ЦФО", "filter:region:ЦФО"),
			tgbotapi.NewInlineKeyboardButtonData("СЗФО", "filter:region:СЗФО"),
			tgbotapi.NewInlineKeyboardButtonData("ЮФО", "filter:region:ЮФО"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ПФО", "filter:region:ПФО"),
			tgbotapi.NewInlineKeyboardButtonData("УФО", "filter:region:УФО"),
			tgbotapi.NewInlineKeyboardButtonData("СФО", "filter:region:СФО"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ДВФО", "filter:region:ДВФО"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", "filter:region:any"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", "filter:back"),
		),
	)
}
