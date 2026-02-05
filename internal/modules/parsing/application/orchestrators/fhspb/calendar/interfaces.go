package calendar

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/calendar"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/match"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/standings"
)

// CalendarConfig интерфейс конфигурации парсера календаря FHSPB
type CalendarConfig interface {
	RequestDelay() int   // Задержка между запросами (мс)
	GameWorkers() int    // Воркеры для матчей
	ParseProtocol() bool // Парсить протокол
	ParseLineups() bool  // Парсить составы
	SkipExisting() bool  // Пропускать существующие матчи
}

// CalendarParser интерфейс парсера календаря
type CalendarParser interface {
	Parse(html []byte, tournamentID int) ([]calendar.MatchDTO, error)
}

// MatchParser интерфейс парсера матча
type MatchParser interface {
	Parse(html []byte) (*match.MatchDetailsDTO, error)
}

// StandingsParser интерфейс парсера таблицы
type StandingsParser interface {
	Parse(html []byte) ([]standings.StandingDTO, error)
}
