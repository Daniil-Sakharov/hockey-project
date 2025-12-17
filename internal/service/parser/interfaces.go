package parser

import (
	"context"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// OrchestratorService координирует весь процесс парсинга
type OrchestratorService interface {
	Run(ctx context.Context) error
	RunJuniorParsing(ctx context.Context) error
	RunRegistryParsing(ctx context.Context) error
	RunMerge(ctx context.Context) error
}

// StatsOrchestratorService координирует парсинг статистики турниров
type StatsOrchestratorService interface {
	Run(ctx context.Context) error
}

// JuniorParserService парсит junior.fhr.ru
type JuniorParserService interface {
	ParseDomains(ctx context.Context) ([]string, error)
	ParseTournaments(ctx context.Context, domain string) ([]junior.TournamentDTO, error)
	ParseAllSeasonsTournaments(ctx context.Context, domain string) ([]junior.TournamentDTO, error)
	ExtractAllSeasons(ctx context.Context, domain string) ([]junior.SeasonInfo, error)
	ParseSeasonTournaments(ctx context.Context, domain, season, ajaxURL string) ([]junior.TournamentDTO, error)
	ParseTeams(ctx context.Context, domain, tournamentURL string) ([]junior.TeamDTO, error)
	ParsePlayers(ctx context.Context, teamURL string) ([]junior.PlayerDTO, error)
}

// StatsParserService парсит статистику игроков из турниров
type StatsParserService interface {
	ParseTournamentStats(
		ctx context.Context,
		domain string,
		tournamentURL string,
		tournamentID string,
	) (int, error)
}

// RegistryParserService парсит registrynew.fhr.ru
type RegistryParserService interface {
	Auth(ctx context.Context) error
	ParseProfile(ctx context.Context, name string, birthDate time.Time) (*RegistryPlayerDTO, error)
}

// MergerService мержит данные из разных источников
type MergerService interface {
	FuzzyMatch(ctx context.Context, juniorPlayer *player.Player, registryPlayers []*RegistryPlayerDTO) (*RegistryPlayerDTO, float64, error)
}

// RegistryPlayerDTO - данные из registrynew.fhr.ru (для будущего)
type RegistryPlayerDTO struct {
	Name       string
	BirthDate  time.Time
	Photo      string
	Additional map[string]string
}
