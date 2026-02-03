package dto

// PlayerStatsDTO статистика игрока из таблицы на странице профиля
// Порядок колонок: Команда, Сезон, Турнир, И, Г, А, О, Ш, ШП, МИН/СЕК
type PlayerStatsDTO struct {
	TeamName       string // Команда - "ЦСКА 2009"
	TeamID         int    // ID команды из ссылки /team/{id}
	Season         string // Сезон - "24/25", "25/26"
	TournamentName string // Турнир - "ПМ 2009 г.р. 25/26"
	Games          int    // И - игры
	Goals          int    // Г - голы
	Assists        int    // А - передачи
	Points         int    // О - очки
	PenaltyCount   int    // Ш - количество штрафов
	PenaltyMinutes int    // ШП - штрафные минуты
	IceTime        string // МИН/СЕК - время на льду
	IceTimeSeconds int    // время на льду в секундах
}

// GoalieStatsDTO статистика вратаря
type GoalieStatsDTO struct {
	Season         string  // "2024/25"
	TeamName       string  // команда
	TeamID         int     // ID команды
	Games          int     // игры
	GoalsAgainst   int     // пропущенные шайбы
	SavePercentage float64 // процент отраженных
	MinutesPlayed  int     // минуты на льду
	GAA            float64 // КН - коэффициент надёжности
}
