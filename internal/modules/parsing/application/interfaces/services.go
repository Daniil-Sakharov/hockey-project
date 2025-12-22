package interfaces

import "context"

// StatsOrchestratorService координирует парсинг статистики турниров
type StatsOrchestratorService interface {
	Run(ctx context.Context) error
}

// StatsParserService парсит статистику игроков из турниров
type StatsParserService interface {
	ParseTournamentStats(
		ctx context.Context,
		domain string,
		tournamentURL string,
		tournamentID string,
		season string,
	) (int, error)
}

// JuniorConfig интерфейс конфигурации Junior парсера
type JuniorConfig interface {
	BaseURL() string
	DomainWorkers() int
	MinBirthYear() int
}
