package app

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
)

// StatsParserApp расширение App для stats_parser процесса
type StatsParserApp struct {
	*initializer.App
}

// NewStatsParserApp создает stats parser приложение
func NewStatsParserApp(ctx context.Context) (*StatsParserApp, error) {
	baseApp, err := initializer.New(ctx)
	if err != nil {
		return nil, err
	}

	return &StatsParserApp{App: baseApp}, nil
}

// Run запускает парсинг статистики через StatsOrchestratorService
func (a *StatsParserApp) Run(ctx context.Context) error {
	orchestrator := a.DiContainer.StatsOrchestratorService(ctx)
	return orchestrator.Run(ctx)
}
