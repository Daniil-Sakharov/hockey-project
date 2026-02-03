package fhmoscow

import (
	"time"

	fhmoscowrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhmoscow"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow"
	"github.com/jmoiron/sqlx"
)

// Config конфигурация парсера FHMoscow
type Config interface {
	MinBirthYear() int
	SeasonWorkers() int
	TournamentWorkers() int
	TeamWorkers() int
	PlayerWorkers() int
	RetryEnabled() bool
	RetryMaxAttempts() int
	RetryDelay() time.Duration
	MaxSeasons() int
	TestSeason() string
	// Player scanning options (since team roster pages are JavaScript-rendered)
	ScanPlayers() bool
	MaxPlayerID() int
}

// Dependencies зависимости сервиса
type Dependencies struct {
	DB                   *sqlx.DB
	Client               *fhmoscow.Client
	TournamentRepo       *fhmoscowrepo.TournamentRepository
	TeamRepo             *fhmoscowrepo.TeamRepository
	PlayerRepo           *fhmoscowrepo.PlayerRepository
	PlayerTeamRepo       *fhmoscowrepo.PlayerTeamRepository
	PlayerStatisticsRepo *fhmoscowrepo.PlayerStatisticsRepository
	GoalieStatisticsRepo *fhmoscowrepo.GoalieStatisticsRepository
}
