package app

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/parser/fhspb/stats"
)

// FHSPBStatsParserApp приложение для парсинга статистики fhspb.ru
type FHSPBStatsParserApp struct {
	*initializer.App
}

// NewFHSPBStatsParserApp создает приложение парсера статистики
func NewFHSPBStatsParserApp(ctx context.Context) (*FHSPBStatsParserApp, error) {
	baseApp, err := initializer.New(ctx)
	if err != nil {
		return nil, err
	}

	return &FHSPBStatsParserApp{App: baseApp}, nil
}

// Run запускает парсер статистики
func (a *FHSPBStatsParserApp) Run(ctx context.Context) error {
	cfg := config.AppConfig().FHSPB

	client := fhspb.NewClient()
	client.SetDelay(cfg.RequestDelay())

	deps := stats.Dependencies{
		Client:               client,
		TournamentRepo:       a.DiContainer.FHSPBTournament(ctx),
		TeamRepo:             a.DiContainer.FHSPBTeam(ctx),
		PlayerRepo:           a.DiContainer.FHSPBPlayer(ctx),
		PlayerStatisticsRepo: a.DiContainer.FHSPBPlayerStatistics(ctx),
		GoalieStatisticsRepo: a.DiContainer.FHSPBGoalieStatistics(ctx),
		StatisticsWorkers:    cfg.StatisticsWorkers(),
	}

	svc := stats.NewService(deps)
	return svc.Orchestrator().Run(ctx)
}
