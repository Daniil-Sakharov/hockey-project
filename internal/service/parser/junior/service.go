package junior

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior"
)

type juniorService struct {
	client *junior.Client
}

// NewJuniorService создает сервис для парсинга junior.fhr.ru
func NewJuniorService(client *junior.Client) *juniorService {
	return &juniorService{
		client: client,
	}
}

// ParseDomains находит все домены *.fhr.ru
func (s *juniorService) ParseDomains(ctx context.Context) ([]string, error) {
	mainURL := "https://junior.fhr.ru"
	return s.client.DiscoverAllDomains(mainURL)
}

// ParseTournaments парсит турниры с домена (текущий сезон)
func (s *juniorService) ParseTournaments(ctx context.Context, domain string) ([]junior.TournamentDTO, error) {
	return s.client.ParseTournamentsFromDomain(domain)
}

// ParseAllSeasonsTournaments парсит турниры ВСЕХ сезонов домена через Worker Pool
func (s *juniorService) ParseAllSeasonsTournaments(ctx context.Context, domain string) ([]junior.TournamentDTO, error) {
	return s.client.ParseAllSeasonsTournaments(ctx, domain)
}

// ExtractAllSeasons извлекает все сезоны домена
func (s *juniorService) ExtractAllSeasons(ctx context.Context, domain string) ([]junior.SeasonInfo, error) {
	return s.client.ExtractAllSeasons(domain)
}

// ParseSeasonTournaments парсит турниры одного сезона
func (s *juniorService) ParseSeasonTournaments(
	ctx context.Context,
	domain, season, ajaxURL string,
) ([]junior.TournamentDTO, error) {
	return s.client.ParseSeasonTournaments(domain, season, ajaxURL)
}

// ParseTeams парсит команды из турнира
func (s *juniorService) ParseTeams(ctx context.Context, domain, tournamentURL string) ([]junior.TeamDTO, error) {
	return s.client.ParseTeamsFromTournament(ctx, domain, tournamentURL)
}

// ParsePlayers парсит игроков из команды
func (s *juniorService) ParsePlayers(ctx context.Context, teamURL string) ([]junior.PlayerDTO, error) {
	return s.client.ParsePlayersFromTeam(teamURL)
}
