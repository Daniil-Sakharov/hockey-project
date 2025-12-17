package message

import (
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/bot"
)

// RenderFioMenu рендерит панель ввода ФИО
func (p *MessagePresenter) RenderFioMenu(fio bot.TempFioData) (string, error) {
	lastName := fio.LastName
	if lastName == "" {
		lastName = "-"
	}

	firstName := fio.FirstName
	if firstName == "" {
		firstName = "-"
	}

	patronymic := fio.Patronymic
	if patronymic == "" {
		patronymic = "-"
	}

	text := fmt.Sprintf(`━━━━━━━━━━━━━━━━
*ФИО игрока*
━━━━━━━━━━━━━━━━
Фамилия: %s
Имя: %s
Отчество: %s

Выберите действие:`, lastName, firstName, patronymic)

	return text, nil
}

// RenderFioInputRequest рендерит запрос на ввод поля
func (p *MessagePresenter) RenderFioInputRequest(field string) string {
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

	return fmt.Sprintf(`Введите %s (минимум 2 символа, кириллица):

Если передумали, нажмите кнопку "Отменить" ниже.`, fieldName)
}
