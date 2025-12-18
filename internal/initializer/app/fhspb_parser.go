package app

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
	svc "github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb/orchestrator"
)

// FHSPBParserApp расширение App для fhspb parser процесса
type FHSPBParserApp struct {
	*initializer.App
}

// NewFHSPBParserApp создает fhspb parser приложение
func NewFHSPBParserApp(ctx context.Context) (*FHSPBParserApp, error) {
	baseApp, err := initializer.New(ctx)
	if err != nil {
		return nil, err
	}

	return &FHSPBParserApp{App: baseApp}, nil
}

// Run запускает парсер fhspb.ru
func (a *FHSPBParserApp) Run(ctx context.Context) error {
	cfg := config.AppConfig().FHSPB

	client := fhspb.NewClient()
	client.SetDelay(cfg.RequestDelay())

	deps := svc.Dependencies{
		DB:             a.DiContainer.PostgresDB(ctx),
		Client:         client,
		TournamentRepo: a.DiContainer.FHSPBTournament(ctx),
		TeamRepo:       a.DiContainer.FHSPBTeam(ctx),
		PlayerRepo:     a.DiContainer.FHSPBPlayer(ctx),
		PlayerTeamRepo: a.DiContainer.FHSPBPlayerTeam(ctx),
	}

	orch := orchestrator.New(deps, cfg)
	return orch.Run(ctx)
}
