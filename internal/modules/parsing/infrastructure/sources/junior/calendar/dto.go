package calendar

import "time"

// MatchDTO представляет матч из календаря
type MatchDTO struct {
	ExternalID  string     // ID из URL (17392431)
	GameURL     string     // URL матча (/games/17392431/)
	HomeTeam    TeamInfo   // Домашняя команда
	AwayTeam    TeamInfo   // Гостевая команда
	HomeScore   *int       // Счёт дома
	AwayScore   *int       // Счёт гостей
	ResultType  string     // regular, OT (овертайм), SO (буллиты)
	ScheduledAt *time.Time // Дата/время матча
	Status      string     // scheduled, finished
	Venue       string     // Арена
	GroupName   string     // Группа (А1, А2, Б...)
	BirthYear   int        // Год рождения
	MatchNumber *int       // Номер матча
}

// TeamInfo информация о команде
type TeamInfo struct {
	ID   string // ID команды (из лого или URL)
	URL  string // URL команды
	Name string // Название
}

// CalendarFilter фильтр для календаря
type CalendarFilter struct {
	GroupName string // Название группы
	BirthYear int    // Год рождения
}
