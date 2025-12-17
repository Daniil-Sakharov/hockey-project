package fhspb

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	fhspbRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
)

// Service интерфейс сервиса парсинга fhspb.ru
type Service interface {
	Run(ctx context.Context) error
}

// Config конфигурация парсера
type Config struct {
	MaxBirthYear      int
	TournamentWorkers int
	TeamWorkers       int
	PlayerWorkers     int
}

// Dependencies зависимости сервиса
type Dependencies struct {
	Client         *fhspb.Client
	TournamentRepo *fhspbRepo.TournamentRepository
	TeamRepo       *fhspbRepo.TeamRepository
	PlayerRepo     *fhspbRepo.PlayerRepository
	PlayerTeamRepo *fhspbRepo.PlayerTeamRepository
}
