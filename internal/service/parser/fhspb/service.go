package fhspb

import (
	"context"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	fhspbRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
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
	TournamentRepo *fhspbRepo.TournamentRepository
	TeamRepo       *fhspbRepo.TeamRepository
	PlayerRepo     *fhspbRepo.PlayerRepository
	PlayerTeamRepo *fhspbRepo.PlayerTeamRepository
}
