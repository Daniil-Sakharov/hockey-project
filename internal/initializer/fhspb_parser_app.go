package initializer

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	svc "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb/orchestrator"
)

// RunFHSPBParser запускает парсер fhspb.ru
func (a *App) RunFHSPBParser(ctx context.Context) error {
	cfg := config.AppConfig().FHSPB

	client := fhspb.NewClient()
	client.SetDelay(cfg.RequestDelay())

	deps := svc.Dependencies{
		Client:         client,
		PlayerRepo:     a.diContainer.PlayerRepository(ctx),
		TeamRepo:       a.diContainer.TeamRepository(ctx),
		TournamentRepo: a.diContainer.TournamentRepository(ctx),
	}

	parserCfg := svc.Config{
		MaxBirthYear:      cfg.MaxBirthYear(),
		TournamentWorkers: cfg.TournamentWorkers(),
		TeamWorkers:       cfg.TeamWorkers(),
		PlayerWorkers:     cfg.PlayerWorkers(),
	}

	orch := orchestrator.New(deps, parserCfg)
	return orch.Run(ctx)
}
