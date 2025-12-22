package presenter

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain/entities"
)

// RenderFioMenu рендерит панель ввода ФИО
func (p *Presenter) RenderFioMenu(fio entities.TempFIOData) (string, error) {
	return p.renderer.Render("fio_menu.tmpl", fio)
}

// RenderFioInputRequest рендерит запрос на ввод поля
func (p *Presenter) RenderFioInputRequest(field string) string {
	var fieldName string
	switch field {
	case "last_name":
		fieldName = "фамилию"
	case "first_name":
		fieldName = "имя"
	case "patronymic":
		fieldName = "отчество"
	default:
		fieldName = "значение"
	}

	return "Введите " + fieldName + " (минимум 2 символа, кириллица):\n\nЕсли передумали, нажмите кнопку \"Отменить\" ниже."
}
