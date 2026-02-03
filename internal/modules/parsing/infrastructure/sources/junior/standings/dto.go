package standings

// StandingDTO позиция команды в турнирной таблице
type StandingDTO struct {
	Position       int    // Место
	TeamURL        string // URL команды
	TeamName       string // Название команды
	Games          int    // Игры
	Wins           int    // Победы (В)
	WinsOT         int    // Победы OT (ВО)
	WinsSO         int    // Победы SO (ВБ)
	LossesSO       int    // Поражения SO (ПБ)
	LossesOT       int    // Поражения OT (ПО)
	Losses         int    // Поражения (П)
	Draws          int    // Ничьи (Н)
	GoalsFor       int    // Забито (ШЗ)
	GoalsAgainst   int    // Пропущено (ШП)
	GoalDifference int    // Разница (+/-)
	Points         int    // Очки (О)
	GroupName      string // Группа (А1, А2...)
	BirthYear      int    // Год рождения
}

// StandingsFilter фильтр для турнирной таблицы
type StandingsFilter struct {
	GroupName string // Название группы
	BirthYear int    // Год рождения
}
