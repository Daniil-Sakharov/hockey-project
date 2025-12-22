package di

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/infrastructure/persistence"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/presenter/keyboard"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/interfaces/bot/template"
)

// TelegramPlayerSearchRepository возвращает репозиторий поиска игроков для telegram
func (c *Container) TelegramPlayerSearchRepository(ctx context.Context) (services.PlayerRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return persistence.NewPlayerSearchRepository(db), nil
}

// TelegramProfileRepository возвращает репозиторий профилей для telegram
func (c *Container) TelegramProfileRepository(ctx context.Context) (services.ProfileRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return persistence.NewProfileRepository(db), nil
}

// TelegramReportRepository возвращает репозиторий отчётов для telegram
func (c *Container) TelegramReportRepository(ctx context.Context) (services.ReportRepository, error) {
	db, err := c.DB(ctx)
	if err != nil {
		return nil, err
	}
	return persistence.NewReportRepository(db), nil
}

// TelegramTemplateEngine возвращает template engine для telegram
func (c *Container) TelegramTemplateEngine() (template.Renderer, error) {
	return template.NewEngine()
}

// TelegramPresenter возвращает presenter для telegram
func (c *Container) TelegramPresenter() (*presenter.Presenter, error) {
	engine, err := c.TelegramTemplateEngine()
	if err != nil {
		return nil, err
	}
	return presenter.NewPresenter(engine), nil
}

// TelegramKeyboardPresenter возвращает keyboard presenter
func (c *Container) TelegramKeyboardPresenter() *keyboard.KeyboardPresenter {
	return keyboard.NewKeyboardPresenter()
}

// TelegramUserStateService возвращает сервис состояния пользователя
func (c *Container) TelegramUserStateService() *services.UserStateService {
	return services.NewUserStateService()
}

// TelegramPlayerSearchService возвращает сервис поиска игроков
func (c *Container) TelegramPlayerSearchService(ctx context.Context) (*services.PlayerSearchService, error) {
	repo, err := c.TelegramPlayerSearchRepository(ctx)
	if err != nil {
		return nil, err
	}
	return services.NewPlayerSearchService(repo), nil
}

// TelegramProfileService возвращает сервис профилей
func (c *Container) TelegramProfileService(ctx context.Context) (*services.ProfileService, error) {
	repo, err := c.TelegramProfileRepository(ctx)
	if err != nil {
		return nil, err
	}
	return services.NewProfileService(repo), nil
}

// TelegramReportService возвращает сервис отчётов
func (c *Container) TelegramReportService(ctx context.Context) (*services.ReportService, error) {
	repo, err := c.TelegramReportRepository(ctx)
	if err != nil {
		return nil, err
	}
	return services.NewReportService(repo), nil
}
