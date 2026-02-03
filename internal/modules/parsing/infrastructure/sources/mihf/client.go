package mihf

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

const (
	BaseURL        = "https://stats.mihf.ru"
	DefaultTimeout = 15 * time.Second
	DefaultDelay   = 150 * time.Millisecond
	MaxRetries     = 3
)

// Client HTTP клиент для работы с stats.mihf.ru
type Client struct {
	httpClient *http.Client
	baseURL    string
	delay      time.Duration
	mu         sync.Mutex
	lastReq    time.Time
}

// NewClient создает новый клиент для stats.mihf.ru
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
			Jar:     jar,
		},
		baseURL: BaseURL,
		delay:   DefaultDelay,
	}
}

// SetDelay устанавливает задержку между запросами
func (c *Client) SetDelay(d time.Duration) {
	c.delay = d
}

// Get выполняет GET запрос с rate limiting
func (c *Client) Get(path string) ([]byte, error) {
	url := c.baseURL + path
	return c.doRequestWithRateLimit(url)
}

// GetURL выполняет GET запрос по полному URL
func (c *Client) GetURL(url string) ([]byte, error) {
	return c.doRequestWithRateLimit(url)
}

// doRequestWithRateLimit выполняет запрос с блокировкой на время всего запроса
// Это гарантирует, что параллельные запросы выполняются последовательно
func (c *Client) doRequestWithRateLimit(url string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ждем минимальную задержку с момента последнего запроса
	elapsed := time.Since(c.lastReq)
	if elapsed < c.delay {
		time.Sleep(c.delay - elapsed)
	}

	// Выполняем запрос под мьютексом
	result, err := c.doRequest(url)

	// Обновляем время последнего запроса после завершения
	c.lastReq = time.Now()

	return result, err
}

// doRequest выполняет HTTP запрос с retry логикой
func (c *Client) doRequest(url string) ([]byte, error) {
	ctx := context.Background()

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
		req.Header.Set("Connection", "keep-alive")

		start := time.Now()
		logger.Debug(ctx, "mihf request", zap.String("url", url), zap.Int("attempt", attempt))

		resp, err := c.httpClient.Do(req)
		elapsed := time.Since(start)

		if err != nil {
			logger.Warn(ctx, "mihf request failed",
				zap.String("url", url),
				zap.Duration("elapsed", elapsed),
				zap.Int("attempt", attempt),
				zap.Error(err),
			)

			if isRetryable(err) && attempt < MaxRetries {
				backoff := time.Duration(1<<(attempt+1)) * time.Second
				time.Sleep(backoff)
				continue
			}
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		logger.Debug(ctx, "mihf response",
			zap.String("url", url),
			zap.Duration("elapsed", elapsed),
			zap.Int("size", len(body)),
		)

		return body, nil
	}

	return nil, fmt.Errorf("all %d attempts exhausted", MaxRetries)
}

func isRetryable(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout")
}
