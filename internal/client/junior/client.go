package junior

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/domain"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/player"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/team"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/junior/tournament"
)

// Client HTTP клиент для работы с junior.fhr.ru
type Client struct {
	httpClient *http.Client
	baseURL    string

	// Парсеры
	Domain     *domain.Parser
	Player     *player.Parser
	Team       *team.Parser
	Tournament *tournament.Parser
}

// NewClient создает новый клиент для junior.fhr.ru
func NewClient() *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		baseURL: "https://cfo.fhr.ru",
	}

	// Инициализируем парсеры
	c.Domain = domain.NewParser(c)
	c.Player = player.NewParser(c, c.baseURL)
	c.Team = team.NewParser(c)
	c.Tournament = tournament.NewParser(c)

	return c
}

// MakeRequest выполняет HTTP запрос с retry логикой (implements types.HTTPRequester)
func (c *Client) MakeRequest(url string) (*http.Response, error) {
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания запроса: %w", err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
		req.Header.Set("Cache-Control", "max-age=0")

		start := time.Now()
		log.Printf("→ Requesting: %s (attempt %d/%d)", url, attempt, maxRetries)

		resp, err := c.httpClient.Do(req)

		elapsed := time.Since(start)
		if err == nil {
			log.Printf("✓ Response received in %.2f seconds (status: %d)", elapsed.Seconds(), resp.StatusCode)
			return resp, nil
		}

		log.Printf("✗ Request failed after %.2f seconds: %v", elapsed.Seconds(), err)

		isRetryable := strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "connection reset") ||
			strings.Contains(err.Error(), "broken pipe") ||
			strings.Contains(err.Error(), "timeout")

		if !isRetryable || attempt == maxRetries {
			return nil, err
		}

		log.Printf("⟳ Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("все %d попытки исчерпаны", maxRetries)
}

// Методы-обёртки для обратной совместимости

func (c *Client) DiscoverAllDomains(mainURL string) ([]string, error) {
	return c.Domain.DiscoverAll(mainURL)
}

func (c *Client) ParsePlayersFromTeam(teamURL string) ([]PlayerDTO, error) {
	return c.Player.ParseFromTeam(teamURL)
}

func (c *Client) ParseTeamsFromTournament(ctx context.Context, domain, tournamentURL string) ([]TeamDTO, error) {
	return c.Team.ParseFromTournament(ctx, domain, tournamentURL)
}

func (c *Client) ParseTournamentsFromDomain(domain string) ([]TournamentDTO, error) {
	return c.Tournament.ParseFromDomain(domain)
}
