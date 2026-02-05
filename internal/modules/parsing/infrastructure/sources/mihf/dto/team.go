package dto

// TeamDTO представляет команду из турнирной таблицы
type TeamDTO struct {
	ID          string // 8977
	Name        string // Динамо
	ExternalURL string // /championat/.../team/8977

	// Статистика команды
	Games        int // И - игры
	Wins         int // В - победы
	Draws        int // Н - ничьи
	Losses       int // П - поражения
	Points       int // О - очки
	GoalsFor     int // ШЗ - забитые
	GoalsAgainst int // ШП - пропущенные
	GoalsDiff    int // РШ - разница
}
