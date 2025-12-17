package di

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
)

// ContainerFactory создает контейнеры для разных типов приложений
type ContainerFactory struct {
	cfg *config.Config
}

func NewContainerFactory(cfg *config.Config) *ContainerFactory {
	return &ContainerFactory{cfg: cfg}
}

// CreateBotContainer создает контейнер для Telegram бота
func (f *ContainerFactory) CreateBotContainer() *BotContainer {
	return NewBotContainer(f.cfg)
}

// CreateParserContainer создает контейнер для парсеров
func (f *ContainerFactory) CreateParserContainer() *ParserContainer {
	return NewParserContainer(f.cfg)
}

// CreateMigrateContainer создает контейнер для миграций
func (f *ContainerFactory) CreateMigrateContainer() *BaseContainer {
	return NewBaseContainer(f.cfg)
}
