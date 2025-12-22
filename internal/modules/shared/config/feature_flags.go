package config

// FeatureFlags управляет переключением между старой и новой архитектурой
type FeatureFlags struct {
	UseNewParsingModule  bool `env:"USE_NEW_PARSING" envDefault:"false"`
	UseNewTelegramModule bool `env:"USE_NEW_TELEGRAM" envDefault:"false"`
	UseEventSourcing     bool `env:"USE_EVENT_SOURCING" envDefault:"false"`
}

// NewFeatureFlags создает feature flags с дефолтными значениями
func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		UseNewParsingModule:  false,
		UseNewTelegramModule: false,
		UseEventSourcing:     false,
	}
}

// IsParsingModuleEnabled проверяет включен ли новый parsing module
func (f *FeatureFlags) IsParsingModuleEnabled() bool {
	return f.UseNewParsingModule
}

// IsTelegramModuleEnabled проверяет включен ли новый telegram module
func (f *FeatureFlags) IsTelegramModuleEnabled() bool {
	return f.UseNewTelegramModule
}

// IsEventSourcingEnabled проверяет включен ли event sourcing
func (f *FeatureFlags) IsEventSourcingEnabled() bool {
	return f.UseEventSourcing
}
