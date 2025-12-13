package bot

// UserState состояние пользователя
type UserState struct {
	UserID                 int64         // ID пользователя
	Filters                SearchFilters // Выбранные фильтры
	LastMsgID              int           // ID последнего сообщения (для EditMessageText)
	CurrentView            string        // Текущий экран (filter_menu, year_select, position_select, etc.)
	WaitingForInput        string        // Ожидание ввода: "fio_last_name", "fio_first_name", "fio_patronymic", ""
	CurrentPage            int           // Текущая страница результатов поиска
	SearchResultMessageIDs []int         // ID сообщений результатов поиска (для удаления при пагинации)
	TempFioFilters         TempFioData   // Временные данные ФИО (до применения)
}

// TempFioData временное хранилище для ввода ФИО
type TempFioData struct {
	LastName   string // Фамилия
	FirstName  string // Имя
	Patronymic string // Отчество
}
