package fhspb

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

const (
	BaseURL        = "https://www.fhspb.ru"
	DefaultTimeout = 30 * time.Second
	DefaultDelay   = 150 * time.Millisecond
)

// Client HTTP клиент для работы с fhspb.ru
type Client struct {
	httpClient *http.Client
	baseURL    string
	delay      time.Duration
	mu         sync.Mutex
	lastReq    time.Time
}

// NewClient создает новый клиент для fhspb.ru
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
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
	c.rateLimit()
	url := c.baseURL + path
	return c.doRequest(url)
}

// rateLimit обеспечивает задержку между запросами
func (c *Client) rateLimit() {
	c.mu.Lock()
	defer c.mu.Unlock()

	elapsed := time.Since(c.lastReq)
	if elapsed < c.delay {
		time.Sleep(c.delay - elapsed)
	}
	c.lastReq = time.Now()
}

// doRequest выполняет HTTP запрос с retry логикой
func (c *Client) doRequest(url string) ([]byte, error) {
	ctx := context.Background()
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")

		start := time.Now()
		logger.Debug(ctx, "fhspb request", zap.String("url", url), zap.Int("attempt", attempt))

		resp, err := c.httpClient.Do(req)
		elapsed := time.Since(start)

		if err != nil {
			logger.Warn(ctx, "fhspb request failed",
				zap.String("url", url),
				zap.Duration("elapsed", elapsed),
				zap.Error(err),
			)

			if isRetryable(err) && attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		logger.Debug(ctx, "fhspb response",
			zap.String("url", url),
			zap.Duration("elapsed", elapsed),
			zap.Int("size", len(body)),
		)

		return body, nil
	}

	return nil, fmt.Errorf("all %d attempts exhausted", maxRetries)
}

func isRetryable(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout")
}
