package keyboard

import (
	cb "github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/callback"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// FilterMenu создает меню фильтров
func (p *KeyboardPresenter) FilterMenu(hasFilters bool) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData("📝 Ввести ФИО", cb.Filter(cb.FilterFio, cb.FioSelect))},
		{tgbotapi.NewInlineKeyboardButtonData("🎂 Год рождения ▼", cb.Filter(cb.FilterYear, cb.SubCmdSelect))},
		{tgbotapi.NewInlineKeyboardButtonData("🏒 Позиция ▼", cb.Filter(cb.FilterPosition, cb.SubCmdSelect))},
		{
			tgbotapi.NewInlineKeyboardButtonData("📏 Рост ▼", cb.Filter(cb.FilterHeight, cb.SubCmdSelect)),
			tgbotapi.NewInlineKeyboardButtonData("⚖️ Вес ▼", cb.Filter(cb.FilterWeight, cb.SubCmdSelect)),
		},
		{tgbotapi.NewInlineKeyboardButtonData("🗺️ Регион ▼", cb.Filter(cb.FilterRegion, cb.SubCmdSelect))},
	}

	if hasFilters {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🔍 НАЙТИ", cb.Filter(cb.FilterApply, "")),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Сбросить", cb.Filter(cb.FilterReset, "")),
		})
	} else {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🔍 НАЙТИ", cb.Filter(cb.FilterApply, "")),
		})
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", cb.Menu(cb.MenuMain)),
	})

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

// YearSelect создает клавиатуру выбора года
func (p *KeyboardPresenter) YearSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2008", cb.Filter(cb.FilterYear, "2008")),
			tgbotapi.NewInlineKeyboardButtonData("2009", cb.Filter(cb.FilterYear, "2009")),
			tgbotapi.NewInlineKeyboardButtonData("2010", cb.Filter(cb.FilterYear, "2010")),
			tgbotapi.NewInlineKeyboardButtonData("2011", cb.Filter(cb.FilterYear, "2011")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("2012", cb.Filter(cb.FilterYear, "2012")),
			tgbotapi.NewInlineKeyboardButtonData("2013", cb.Filter(cb.FilterYear, "2013")),
			tgbotapi.NewInlineKeyboardButtonData("2014", cb.Filter(cb.FilterYear, "2014")),
			tgbotapi.NewInlineKeyboardButtonData("2015", cb.Filter(cb.FilterYear, "2015")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", cb.Filter(cb.FilterYear, cb.ValueAny)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", cb.Filter(cb.FilterBack, "")),
		),
	)
}

// PositionSelect создает клавиатуру выбора позиции
func (p *KeyboardPresenter) PositionSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎯 Нападающий", cb.Filter(cb.FilterPosition, "forward")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛡️ Защитник", cb.Filter(cb.FilterPosition, "defender")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥅 Вратарь", cb.Filter(cb.FilterPosition, "goalie")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любая", cb.Filter(cb.FilterPosition, cb.ValueAny)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", cb.Filter(cb.FilterBack, "")),
		),
	)
}

// HeightSelect создает клавиатуру выбора роста
func (p *KeyboardPresenter) HeightSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("150-160", cb.Filter(cb.FilterHeight, "150-160")),
			tgbotapi.NewInlineKeyboardButtonData("160-170", cb.Filter(cb.FilterHeight, "160-170")),
			tgbotapi.NewInlineKeyboardButtonData("170-180", cb.Filter(cb.FilterHeight, "170-180")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("180-190", cb.Filter(cb.FilterHeight, "180-190")),
			tgbotapi.NewInlineKeyboardButtonData("190-200", cb.Filter(cb.FilterHeight, "190-200")),
			tgbotapi.NewInlineKeyboardButtonData("200+", cb.Filter(cb.FilterHeight, "200-250")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", cb.Filter(cb.FilterHeight, cb.ValueAny)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", cb.Filter(cb.FilterBack, "")),
		),
	)
}

// WeightSelect создает клавиатуру выбора веса
func (p *KeyboardPresenter) WeightSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("40-50", cb.Filter(cb.FilterWeight, "40-50")),
			tgbotapi.NewInlineKeyboardButtonData("50-60", cb.Filter(cb.FilterWeight, "50-60")),
			tgbotapi.NewInlineKeyboardButtonData("60-70", cb.Filter(cb.FilterWeight, "60-70")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("70-80", cb.Filter(cb.FilterWeight, "70-80")),
			tgbotapi.NewInlineKeyboardButtonData("80-90", cb.Filter(cb.FilterWeight, "80-90")),
			tgbotapi.NewInlineKeyboardButtonData("90+", cb.Filter(cb.FilterWeight, "90-150")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", cb.Filter(cb.FilterWeight, cb.ValueAny)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", cb.Filter(cb.FilterBack, "")),
		),
	)
}

// RegionSelect создает клавиатуру выбора региона
func (p *KeyboardPresenter) RegionSelect() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ФХР", cb.Filter(cb.FilterRegion, "ФХР")),
			tgbotapi.NewInlineKeyboardButtonData("СПБ", cb.Filter(cb.FilterRegion, "СПБ")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ЦФО", cb.Filter(cb.FilterRegion, "ЦФО")),
			tgbotapi.NewInlineKeyboardButtonData("СЗФО", cb.Filter(cb.FilterRegion, "СЗФО")),
			tgbotapi.NewInlineKeyboardButtonData("ЮФО", cb.Filter(cb.FilterRegion, "ЮФО")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ПФО", cb.Filter(cb.FilterRegion, "ПФО")),
			tgbotapi.NewInlineKeyboardButtonData("УФО", cb.Filter(cb.FilterRegion, "УФО")),
			tgbotapi.NewInlineKeyboardButtonData("СФО", cb.Filter(cb.FilterRegion, "СФО")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ДВФО", cb.Filter(cb.FilterRegion, "ДВФО")),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Любой", cb.Filter(cb.FilterRegion, cb.ValueAny)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад к фильтрам", cb.Filter(cb.FilterBack, "")),
		),
	)
}
