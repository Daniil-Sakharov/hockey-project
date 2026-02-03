package parser

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
)

// JuniorParserService интерфейс для парсинга junior.fhr.ru
type JuniorParserService interface {
	ParseDomains(ctx context.Context) ([]string, error)
	ParseTournaments(ctx context.Context, domain string) ([]types.TournamentDTO, error)
	ParseAllSeasonsTournaments(ctx context.Context, domain string) ([]types.TournamentDTO, error)
	ExtractAllSeasons(ctx context.Context, domain string) ([]junior.SeasonInfo, error)
	ParseSeasonTournaments(ctx context.Context, domain, season, ajaxURL string) ([]types.TournamentDTO, error)
	ParseTeams(ctx context.Context, domain, tournamentURL string, fallbackBirthYears ...int) ([]types.TeamWithContext, error)
	ParsePlayers(ctx context.Context, domain, teamURL string) ([]types.PlayerDTO, error)
	ParsePlayerProfile(ctx context.Context, domain, profileURL string) (*types.PlayerProfileDTO, error)
}

// JuniorConfig интерфейс конфигурации Junior парсера
type JuniorConfig interface {
	BaseURL() string
	DomainWorkers() int
	MinBirthYear() int
	MaxTournaments() int
}
