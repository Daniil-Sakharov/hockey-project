package calendar

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/repositories"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
)

// Source константа источника (соответствует junior parser)
const Source = "junior"

// Orchestrator оркестратор парсинга календаря Junior
type Orchestrator struct {
	// HTTP клиент для AJAX-запросов
	http types.HTTPRequester

	// Парсеры
	calendarParser  CalendarParser
	gameParser      GameParser
	standingsParser StandingsParser
	profileParser   PlayerProfileParser

	// Репозитории
	matchRepo       repositories.MatchRepository
	matchEventRepo  repositories.MatchEventRepository
	matchLineupRepo repositories.MatchLineupRepository
	standingRepo    repositories.StandingRepository
	tournamentRepo  repositories.TournamentRepository
	teamRepo        repositories.TeamRepository
	playerRepo      repositories.PlayerRepository

	// Конфигурация
	config CalendarConfig
}

// NewOrchestrator создает новый оркестратор
func NewOrchestrator(
	http types.HTTPRequester,
	calendarParser CalendarParser,
	gameParser GameParser,
	standingsParser StandingsParser,
	profileParser PlayerProfileParser,
	matchRepo repositories.MatchRepository,
	matchEventRepo repositories.MatchEventRepository,
	matchLineupRepo repositories.MatchLineupRepository,
	standingRepo repositories.StandingRepository,
	tournamentRepo repositories.TournamentRepository,
	teamRepo repositories.TeamRepository,
	playerRepo repositories.PlayerRepository,
	config CalendarConfig,
) *Orchestrator {
	return &Orchestrator{
		http:            http,
		calendarParser:  calendarParser,
		gameParser:      gameParser,
		standingsParser: standingsParser,
		profileParser:   profileParser,
		matchRepo:       matchRepo,
		matchEventRepo:  matchEventRepo,
		matchLineupRepo: matchLineupRepo,
		standingRepo:    standingRepo,
		tournamentRepo:  tournamentRepo,
		teamRepo:        teamRepo,
		playerRepo:      playerRepo,
		config:          config,
	}
}
