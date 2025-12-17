package fhspb

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

// GetPlayerStatsFirstPage загружает первую страницу статистики и возвращает данные + ViewState
func (c *Client) GetPlayerStatsFirstPage(ctx context.Context, tournamentID int) ([]dto.PlayerStatsDTO, dto.StatsPageDTO, error) {
	path := fmt.Sprintf("/StatsPlayer?TournamentID=%d", tournamentID)
	body, err := c.Get(path)
	if err != nil {
		return nil, dto.StatsPageDTO{}, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, dto.StatsPageDTO{}, fmt.Errorf("parse html: %w", err)
	}

	pageInfo := parsing.ParseStatsPage(doc)
	stats := parsing.ParsePlayerStats(doc)

	logger.Debug(ctx, "parsed player stats first page",
		zap.Int("tournament_id", tournamentID),
		zap.Int("total_pages", pageInfo.TotalPages),
		zap.Int("players", len(stats)),
	)

	return stats, pageInfo, nil
}

// GetPlayerStatsPage загружает конкретную страницу статистики используя ViewState
func (c *Client) GetPlayerStatsPage(ctx context.Context, tournamentID, page int, pageInfo dto.StatsPageDTO) ([]dto.PlayerStatsDTO, error) {
	c.rateLimit()

	reqURL := fmt.Sprintf("%s/StatsPlayer?TournamentID=%d", c.baseURL, tournamentID)

	data := url.Values{}
	data.Set("__EVENTTARGET", "ctl00$ctl00$MainContent$MainContent$StatsGridView")
	data.Set("__EVENTARGUMENT", fmt.Sprintf("Page$%d", page))
	data.Set("__VIEWSTATE", pageInfo.ViewState)
	data.Set("__VIEWSTATEGENERATOR", pageInfo.ViewStateGenerator)
	data.Set("__EVENTVALIDATION", pageInfo.EventValidation)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	stats := parsing.ParsePlayerStats(doc)

	logger.Debug(ctx, "parsed player stats page",
		zap.Int("tournament_id", tournamentID),
		zap.Int("page", page),
		zap.Int("players", len(stats)),
		zap.Duration("elapsed", time.Since(start)),
	)

	return stats, nil
}
