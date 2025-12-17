package app

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer"
)

// ParserApp расширение App для parser процесса
type ParserApp struct {
	*initializer.App
}

// NewParserApp создает parser приложение
func NewParserApp(ctx context.Context) (*ParserApp, error) {
	baseApp, err := initializer.New(ctx)
	if err != nil {
		return nil, err
	}

	return &ParserApp{App: baseApp}, nil
}

// Run запускает парсинг через OrchestratorService
func (a *ParserApp) Run(ctx context.Context) error {
	orchestrator := a.DiContainer.OrchestratorService(ctx)
	return orchestrator.Run(ctx)
}
