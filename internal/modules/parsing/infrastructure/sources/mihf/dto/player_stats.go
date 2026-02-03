package dto

// PlayerStatsDTO статистика полевого игрока
type PlayerStatsDTO struct {
	ID         string // 23510
	Name       string // Голубев Кирилл
	Number     string // 17
	Position   string // З (защитник) / Н (нападающий)
	ProfileURL string // /players/info/23510

	// Основная статистика
	Games          int // И - игры
	Goals          int // Г - голы
	Assists        int // А - передачи
	Points         int // О - очки
	PenaltyMinutes int // Ш - штрафные минуты

	// Детальная статистика голов
	GoalsPowerPlay    int // ГБ - голы в большинстве
	GoalsShortHanded  int // ГМ - голы в меньшинстве
	GoalsEvenStrength int // ГР - голы в равных составах
}

// GoalieStatsDTO статистика вратаря
type GoalieStatsDTO struct {
	ID         string // 19980
	Name       string // Заболотный Никита
	Number     string // 30
	ProfileURL string // /players/info/19980

	// Статистика вратаря
	Games          int     // И - игры
	Goals          int     // Г - голы (забитые вратарём)
	Assists        int     // А - передачи
	Points         int     // О - очки
	SavePercentage float64 // КН - коэффициент надёжности
	GoalsAgainst   int     // ПШ - пропущено шайб
	MinutesPlayed  int     // Время на льду (в минутах)
}
