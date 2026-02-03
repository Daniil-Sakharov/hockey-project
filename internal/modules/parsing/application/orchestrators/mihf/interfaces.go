package mihf

import (
	"time"

	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf"
	"github.com/jmoiron/sqlx"
)

// Config конфигурация парсера MIHF
type Config interface {
	MinBirthYear() int
	MaxBirthYear() int
	SeasonWorkers() int
	TournamentWorkers() int
	TeamWorkers() int
	PlayerWorkers() int
	RetryEnabled() bool
	RetryMaxAttempts() int
	RetryDelay() time.Duration
	MaxSeasons() int
	TestSeason() string
}

// Dependencies зависимости сервиса
type Dependencies struct {
	DB                   *sqlx.DB
	Client               *mihf.Client
	TournamentRepo       *mihfrepo.TournamentRepository
	TeamRepo             *mihfrepo.TeamRepository
	PlayerRepo           *mihfrepo.PlayerRepository
	PlayerTeamRepo       *mihfrepo.PlayerTeamRepository
	PlayerStatisticsRepo *mihfrepo.PlayerStatisticsRepository
	GoalieStatisticsRepo *mihfrepo.GoalieStatisticsRepository
}
