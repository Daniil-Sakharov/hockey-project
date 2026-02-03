package calendar

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/calendar"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/game"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/standings"
	"github.com/PuerkitoBio/goquery"
)

// CalendarConfig интерфейс конфигурации парсера календаря
type CalendarConfig interface {
	RequestDelay() int           // Задержка между запросами (мс)
	TournamentWorkers() int      // Воркеры для турниров
	GameWorkers() int            // Воркеры для матчей
	ParseProtocol() bool         // Парсить протокол
	ParseLineups() bool          // Парсить составы
	SkipExisting() bool          // Пропускать существующие матчи
	MaxTournaments() int         // Глобальный лимит турниров (0 = без лимита)
}

// CalendarParser интерфейс парсера календаря
type CalendarParser interface {
	Parse(tournamentURL string) ([]calendar.MatchDTO, error)
	ParseWithFilter(baseURL, ajaxURL string, filter calendar.CalendarFilter) ([]calendar.MatchDTO, error)
	ParseFromDoc(doc *goquery.Document, filter calendar.CalendarFilter) ([]calendar.MatchDTO, error)
}

// GameParser интерфейс парсера матча
type GameParser interface {
	Parse(gameURL string) (*game.GameDetailsDTO, error)
}

// StandingsParser интерфейс парсера таблицы
type StandingsParser interface {
	Parse(tournamentURL string) ([]standings.StandingDTO, error)
	ParseWithFilter(baseURL, ajaxURL string, filter standings.StandingsFilter) ([]standings.StandingDTO, error)
	ParseFromDoc(doc *goquery.Document, filter standings.StandingsFilter) ([]standings.StandingDTO, error)
}

// PlayerProfileParser интерфейс парсера профиля игрока
type PlayerProfileParser interface {
	ParseProfile(ctx context.Context, profileURL string) (PlayerInfo, error)
}

// PlayerInfo данные игрока из профиля
type PlayerInfo struct {
	ID          string
	ProfileURL  string
	Name        string
	BirthDate   string
	Position    string
	Citizenship string
}
