package bot

// UserState состояние пользователя
type UserState struct {
	UserID                 int64
	Filters                SearchFilters
	LastMsgID              int
	CurrentView            string
	WaitingForInput        string
	CurrentPage            int
	SearchResultMessageIDs []int
	TempFioFilters         TempFioData
}

// TempFioData временное хранилище для ввода ФИО
type TempFioData struct {
	LastName   string
	FirstName  string
	Patronymic string
}
