package parser

import (
	"context"
	"time"

	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
	"github.com/jmoiron/sqlx"
)

// Service интерфейс сервиса парсинга fhspb.ru
type Service interface {
	Run(ctx context.Context) error
}

// Config конфигурация парсера
type Config interface {
	MaxBirthYear() int
	TournamentWorkers() int
	TeamWorkers() int
	PlayerWorkers() int
	Mode() string
	RetryEnabled() bool
	RetryMaxAttempts() int
	RetryDelay() time.Duration
}

// Dependencies зависимости сервиса
type Dependencies struct {
	DB             *sqlx.DB
	Client         *fhspb.Client
	TournamentRepo *fhspbrepo.TournamentRepository
	TeamRepo       *fhspbrepo.TeamRepository
	PlayerRepo     *fhspbrepo.PlayerRepository
	PlayerTeamRepo *fhspbrepo.PlayerTeamRepository
}
