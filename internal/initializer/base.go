package initializer

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
	"github.com/Daniil-Sakharov/HockeyProject/internal/initializer/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/closer"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

// App базовая структура приложения (общая для parser и bot)
type App struct {
	Config      *config.Config
	DiContainer *di.Container
}

// New создает базовое приложение
func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// initDeps общие инициализации для обоих процессов
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initDI,
		a.initCloser,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	if err := config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.Config = config.AppConfig()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	if err := logger.Init(a.Config.Logger.Level(), a.Config.Logger.AsJson(), nil); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	closer.AddNamed("Logger", func(ctx context.Context) error {
		return logger.Shutdown(ctx)
	})

	logger.Info(ctx, "✅ Logger initialized")
	return nil
}

func (a *App) initDI(ctx context.Context) error {
	a.DiContainer = di.NewContainer(a.Config)
	logger.Info(ctx, "✅ DI container initialized")
	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
