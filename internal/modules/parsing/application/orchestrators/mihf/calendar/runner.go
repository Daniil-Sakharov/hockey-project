package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает парсинг календаря MIHF
func (o *Orchestrator) Run(ctx context.Context) error {
	start := time.Now()

	logger.Info(ctx, "========================================")
	logger.Info(ctx, "MIHF Calendar Parser starting",
		zap.Int("min_birth_year", o.config.MinBirthYear()),
		zap.Bool("parse_protocol", o.config.ParseProtocol()),
		zap.Int("game_workers", o.config.GameWorkers()),
	)
	logger.Info(ctx, "========================================")

	// 1. Получаем сезоны
	seasons, err := o.fetchSeasons(ctx)
	if err != nil {
		return fmt.Errorf("fetch seasons: %w", err)
	}

	seasons = o.filterSeasons(ctx, seasons)
	logger.Info(ctx, "Seasons to process", zap.Int("count", len(seasons)))

	// 2. Для каждого сезона получаем турниры и парсим календари
	var totalMatches, totalEvents int
	for _, season := range seasons {
		matches, events, err := o.processSeason(ctx, season)
		if err != nil {
			logger.Error(ctx, "Season processing failed",
				zap.String("season", season.Year),
				zap.Error(err),
			)
			continue
		}
		totalMatches += matches
		totalEvents += events
	}

	elapsed := time.Since(start)
	logger.Info(ctx, "========================================")
	logger.Info(ctx, "MIHF Calendar Parser completed",
		zap.Duration("elapsed", elapsed),
		zap.Int("seasons_processed", len(seasons)),
		zap.Int("matches_saved", totalMatches),
		zap.Int("events_saved", totalEvents),
	)
	logger.Info(ctx, "========================================")

	return nil
}

// fetchSeasons получает список сезонов
func (o *Orchestrator) fetchSeasons(ctx context.Context) ([]dto.SeasonDTO, error) {
	html, err := o.client.Get("/")
	if err != nil {
		return nil, fmt.Errorf("get main page: %w", err)
	}

	return parsing.ParseSeasons(html)
}

// filterSeasons фильтрует сезоны по конфигурации
func (o *Orchestrator) filterSeasons(ctx context.Context, seasons []dto.SeasonDTO) []dto.SeasonDTO {
	if testSeason := o.config.TestSeason(); testSeason != "" {
		logger.Info(ctx, "Test mode: filtering to season", zap.String("season", testSeason))
		for _, s := range seasons {
			if s.Year == testSeason {
				return []dto.SeasonDTO{s}
			}
		}
		logger.Warn(ctx, "Test season not found")
	}

	if max := o.config.MaxSeasons(); max > 0 && len(seasons) > max {
		return seasons[:max]
	}
	return seasons
}
