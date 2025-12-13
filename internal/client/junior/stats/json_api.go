package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// StatsAPIPath путь к AJAX endpoint для статистики
	StatsAPIPath = "/fhr-ajax/SXRwcm9maXRcQXBwXEFqYXhcQWpheFN0YXRz/Z2V0RGF0YVRhYmxlU2ltcGxlU3RhdHM=/"
)

// FetchStatistics делает GET запрос к JSON API и возвращает статистику
func FetchStatistics(
	ctx context.Context,
	httpClient *http.Client,
	domain string,
	tournamentID string,
	yearID string,
	groupID string,
) (*StatsResponse, error) {
	apiURL := buildAPIURL(domain, tournamentID, yearID, groupID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch statistics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var statsResp StatsResponse
	if err := json.Unmarshal(body, &statsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &statsResp, nil
}

// buildAPIURL формирует URL для запроса к API статистики
func buildAPIURL(domain, tournamentID, yearID, groupID string) string {
	baseURL := fmt.Sprintf("%s%s", domain, StatsAPIPath)

	params := url.Values{}
	params.Set("key", "scorers")
	params.Set("comp", tournamentID)
	params.Set("year", yearID)
	
	// КРИТИЧНО: "Общая статистика" = БЕЗ параметра group!
	// API НЕ поддерживает group=all, он возвращает 0 записей
	// Правильно: для "Общей статистики" вообще НЕ передавать параметр group
	if groupID != "all" && groupID != "" {
		params.Set("group", groupID)
	}
	
	params.Set("season", "2025-2026")
	
	// КРИТИЧНО: Пагинация для получения ВСЕХ записей
	params.Set("start", "0")
	params.Set("length", "9999")  // Максимальное количество записей

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}
