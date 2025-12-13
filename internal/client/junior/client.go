package junior

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client HTTP клиент для работы с junior.fhr.ru
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient создает новый клиент для junior.fhr.ru
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Увеличен таймаут т.к. сайт медленный
		},
		baseURL: "https://cfo.fhr.ru", // Пока хардкодим, потом вынесем в конфиг
	}
}

// makeRequest выполняет HTTP запрос с retry логикой
func (c *Client) makeRequest(url string) (*http.Response, error) {
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания запроса: %w", err)
		}

		// Реалистичные заголовки
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
		req.Header.Set("Cache-Control", "max-age=0")

		// Логируем запрос
		start := time.Now()
		log.Printf("→ Requesting: %s (attempt %d/%d)", url, attempt, maxRetries)

		resp, err := c.httpClient.Do(req)

		elapsed := time.Since(start)
		if err == nil {
			log.Printf("✓ Response received in %.2f seconds (status: %d)", elapsed.Seconds(), resp.StatusCode)
			return resp, nil
		}

		log.Printf("✗ Request failed after %.2f seconds: %v", elapsed.Seconds(), err)

		// Retry для временных ошибок
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
