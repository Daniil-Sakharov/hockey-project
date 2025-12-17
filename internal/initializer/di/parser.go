package di

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
)

// ParserContainer содержит зависимости только для парсеров (без Telegram)
type ParserContainer struct {
	*BaseContainer
}

func NewParserContainer(cfg *config.Config) *ParserContainer {
	return &ParserContainer{
		BaseContainer: NewBaseContainer(cfg),
	}
}
