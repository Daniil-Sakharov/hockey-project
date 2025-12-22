package telegram

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/config/modules"
	logctx "github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/logging/context"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/handlers"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/domain"
)

// Module represents the Telegram bot module.
type Module struct {
	config         *modules.TelegramConfig
	logger         *logctx.ContextualLogger
	filterService  *handlers.FilterService
	sessionService *handlers.SessionService
}

// NewModule creates a new Telegram module.
func NewModule(cfg *modules.TelegramConfig, sessions domain.SessionRepository) *Module {
	logger := logctx.NewContextualLogger("telegram")

	return &Module{
		config:         cfg,
		logger:         logger,
		filterService:  handlers.NewFilterService(sessions),
		sessionService: handlers.NewSessionService(sessions),
	}
}

// Name returns the module name.
func (m *Module) Name() string {
	return "telegram"
}

// Config returns telegram configuration.
func (m *Module) Config() *modules.TelegramConfig {
	return m.config
}

// FilterService returns filter use case.
func (m *Module) FilterService() *handlers.FilterService {
	return m.filterService
}

// SessionService returns session use case.
func (m *Module) SessionService() *handlers.SessionService {
	return m.sessionService
}
