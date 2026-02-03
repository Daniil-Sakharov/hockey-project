package parser

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/types"
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

// ParseDomains возвращает все домены *.fhr.ru для парсинга
// TODO: вернуть полный список после тестирования
func (s *juniorService) ParseDomains(ctx context.Context) ([]string, error) {
	return []string{
		"https://ufo.fhr.ru",
		// "https://junior.fhr.ru",
		// "https://cfo.fhr.ru",
		// "https://dfo.fhr.ru",
		// "https://komi.fhr.ru",
		// "https://kuzbass.fhr.ru",
		// "https://len.fhr.ru",
		// "https://nsk.fhr.ru",
		// "https://pfo.fhr.ru",
		// "https://sam.fhr.ru",
		// "https://sfo.fhr.ru",
		// "https://spb.fhr.ru",
		// "https://szfo.fhr.ru",
		// "https://vrn.fhr.ru",
		// "https://yfo.fhr.ru",
	}, nil
}

// ParseTournaments парсит турниры с домена (текущий сезон)
func (s *juniorService) ParseTournaments(ctx context.Context, domain string) ([]types.TournamentDTO, error) {
	return s.client.ParseTournamentsFromDomain(domain)
}

// ParseAllSeasonsTournaments парсит турниры ВСЕХ сезонов домена через Worker Pool
func (s *juniorService) ParseAllSeasonsTournaments(ctx context.Context, domain string) ([]types.TournamentDTO, error) {
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
) ([]types.TournamentDTO, error) {
	return s.client.ParseSeasonTournaments(domain, season, ajaxURL)
}

// ParseTeams парсит команды из турнира с контекстом года/группы
func (s *juniorService) ParseTeams(ctx context.Context, domain, tournamentURL string, fallbackBirthYears ...int) ([]types.TeamWithContext, error) {
	return s.client.ParseTeamsFromTournament(ctx, domain, tournamentURL, fallbackBirthYears...)
}

// ParsePlayers парсит игроков из команды
func (s *juniorService) ParsePlayers(ctx context.Context, domain, teamURL string) ([]types.PlayerDTO, error) {
	return s.client.ParsePlayersFromTeam(domain, teamURL)
}

// ParsePlayerProfile парсит профиль игрока для получения дополнительных данных
func (s *juniorService) ParsePlayerProfile(ctx context.Context, domain, profileURL string) (*types.PlayerProfileDTO, error) {
	return s.client.ParsePlayerProfile(domain, profileURL)
}
