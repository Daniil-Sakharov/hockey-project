package calendar

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	fhspbClient "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb"
)

// Source константа источника
const Source = "fhspb"

// Orchestrator оркестратор парсинга календаря FHSPB
type Orchestrator struct {
	// HTTP клиент
	client *fhspbClient.Client

	// Парсеры
	calendarParser  CalendarParser
	matchParser     MatchParser
	standingsParser StandingsParser

	// Репозитории FHSPB
	tournamentRepo *fhspb.TournamentRepository
	teamRepo       *fhspb.TeamRepository
	playerRepo     *fhspb.PlayerRepository

	// Общие репозитории
	matchRepo          repositories.MatchRepository
	matchEventRepo     repositories.MatchEventRepository
	matchLineupRepo    repositories.MatchLineupRepository
	standingRepo       repositories.StandingRepository
	matchTeamStatsRepo repositories.MatchTeamStatsRepository

	// Конфигурация
	config CalendarConfig
}

// NewOrchestrator создает новый оркестратор
func NewOrchestrator(
	client *fhspbClient.Client,
	calendarParser CalendarParser,
	matchParser MatchParser,
	standingsParser StandingsParser,
	tournamentRepo *fhspb.TournamentRepository,
	teamRepo *fhspb.TeamRepository,
	playerRepo *fhspb.PlayerRepository,
	matchRepo repositories.MatchRepository,
	matchEventRepo repositories.MatchEventRepository,
	matchLineupRepo repositories.MatchLineupRepository,
	standingRepo repositories.StandingRepository,
	matchTeamStatsRepo repositories.MatchTeamStatsRepository,
	config CalendarConfig,
) *Orchestrator {
	return &Orchestrator{
		client:             client,
		calendarParser:     calendarParser,
		matchParser:        matchParser,
		standingsParser:    standingsParser,
		tournamentRepo:     tournamentRepo,
		teamRepo:           teamRepo,
		playerRepo:         playerRepo,
		matchRepo:          matchRepo,
		matchEventRepo:     matchEventRepo,
		matchLineupRepo:    matchLineupRepo,
		standingRepo:       standingRepo,
		matchTeamStatsRepo: matchTeamStatsRepo,
		config:             config,
	}
}
