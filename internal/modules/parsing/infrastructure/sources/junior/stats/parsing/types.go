package parsing

// StatsCombination представляет комбинацию год+группа для парсинга
type StatsCombination struct {
	YearID    string // ID года в системе
	YearLabel string // Текст года (например "2009")
	GroupID   string // ID группы ("all" для общей статистики)
	GroupName string // Название группы
}

// YearInfo информация о годе из dropdown
type YearInfo struct {
	ID      string // ID года (value из option)
	Label   string // Текст года (2009, 2010...)
	AjaxURL string // URL для AJAX запроса
}

// GroupInfo информация о группе
type GroupInfo struct {
	ID   string
	Name string
}
