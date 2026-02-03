package fhmoscow

import (
	"bytes"
	"context"
	"encoding/json"
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
	BaseURL        = "https://www.fhmoscow.com"
	DefaultTimeout = 30 * time.Second
	DefaultDelay   = 150 * time.Millisecond
	MaxRetries     = 3
)

// Client HTTP клиент для работы с fhmoscow.com API
type Client struct {
	httpClient *http.Client
	baseURL    string
	delay      time.Duration
	mu         sync.Mutex
	lastReq    time.Time
}

// NewClient создает новый клиент для fhmoscow.com
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

// GetHTML выполняет GET запрос и возвращает HTML
func (c *Client) GetHTML(path string) ([]byte, error) {
	url := c.baseURL + path
	return c.doRequestWithRateLimit(http.MethodGet, url, nil, "text/html")
}

// GetAPI выполняет GET запрос к API и возвращает JSON
func (c *Client) GetAPI(path string) ([]byte, error) {
	url := c.baseURL + path
	return c.doRequestWithRateLimit(http.MethodGet, url, nil, "application/json")
}

// PostAPI выполняет POST запрос к API с JSON телом
func (c *Client) PostAPI(path string, body interface{}) ([]byte, error) {
	url := c.baseURL + path
	return c.doRequestWithRateLimit(http.MethodPost, url, body, "application/json")
}

// doRequestWithRateLimit выполняет запрос с блокировкой
func (c *Client) doRequestWithRateLimit(method, url string, body interface{}, accept string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elapsed := time.Since(c.lastReq)
	if elapsed < c.delay {
		time.Sleep(c.delay - elapsed)
	}

	result, err := c.doRequest(method, url, body, accept)
	c.lastReq = time.Now()

	return result, err
}

// doRequest выполняет HTTP запрос с retry логикой
func (c *Client) doRequest(method, url string, body interface{}, accept string) ([]byte, error) {
	ctx := context.Background()

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		var reqBody io.Reader
		if body != nil {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("marshal body: %w", err)
			}
			reqBody = bytes.NewReader(jsonBody)
		}

		req, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", accept)
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
		req.Header.Set("Connection", "keep-alive")

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		start := time.Now()
		logger.Debug(ctx, "fhmoscow request",
			zap.String("method", method),
			zap.String("url", url),
			zap.Int("attempt", attempt),
		)

		resp, err := c.httpClient.Do(req)
		requestElapsed := time.Since(start)

		if err != nil {
			logger.Warn(ctx, "fhmoscow request failed",
				zap.String("url", url),
				zap.Duration("elapsed", requestElapsed),
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
			if resp.StatusCode >= 500 && attempt < MaxRetries {
				backoff := time.Duration(1<<(attempt+1)) * time.Second
				time.Sleep(backoff)
				continue
			}
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read body: %w", err)
		}

		logger.Debug(ctx, "fhmoscow response",
			zap.String("url", url),
			zap.Duration("elapsed", requestElapsed),
			zap.Int("size", len(respBody)),
		)

		return respBody, nil
	}

	return nil, fmt.Errorf("all %d attempts exhausted", MaxRetries)
}

func isRetryable(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout")
}
