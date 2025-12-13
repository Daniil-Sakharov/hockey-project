package fhspb

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
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
	PlayerRepo     player.Repository
	TeamRepo       team.Repository
	TournamentRepo tournament.Repository
}
