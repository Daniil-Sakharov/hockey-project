package initializer

import (
	"context"
)

// StatsParserApp расширение App для stats_parser процесса
type StatsParserApp struct {
	*App
}

// NewStatsParserApp создает stats parser приложение
func NewStatsParserApp(ctx context.Context) (*StatsParserApp, error) {
	baseApp, err := New(ctx)
	if err != nil {
		return nil, err
	}

	return &StatsParserApp{
		App: baseApp,
	}, nil
}

// Run запускает парсинг статистики через StatsOrchestratorService
func (a *StatsParserApp) Run(ctx context.Context) error {
	orchestrator := a.diContainer.StatsOrchestratorService(ctx)
	return orchestrator.Run(ctx)
}
