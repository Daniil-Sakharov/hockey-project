package initializer

import (
	"context"
)

// ParserApp расширение App для parser процесса
type ParserApp struct {
	*App
}

// NewParserApp создает parser приложение
func NewParserApp(ctx context.Context) (*ParserApp, error) {
	baseApp, err := New(ctx)
	if err != nil {
		return nil, err
	}

	return &ParserApp{
		App: baseApp,
	}, nil
}

// Run запускает парсинг через OrchestratorService
func (a *ParserApp) Run(ctx context.Context) error {
	orchestrator := a.diContainer.OrchestratorService(ctx)
	return orchestrator.Run(ctx)
}
