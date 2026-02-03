package mihf

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает полный процесс парсинга MIHF
func (o *Orchestrator) Run(ctx context.Context) error {
	start := time.Now()

	logger.Info(ctx, "========================================")
	logger.Info(ctx, "MIHF Parser starting",
		zap.Int("min_birth_year", o.config.MinBirthYear()),
		zap.Int("max_seasons", o.config.MaxSeasons()),
		zap.String("test_season", o.config.TestSeason()),
		zap.Int("season_workers", o.config.SeasonWorkers()),
		zap.Int("tournament_workers", o.config.TournamentWorkers()),
		zap.Int("team_workers", o.config.TeamWorkers()),
		zap.Int("player_workers", o.config.PlayerWorkers()),
	)
	logger.Info(ctx, "========================================")

	// 1. Получаем список сезонов
	logger.Info(ctx, "[STEP 1] Fetching seasons from main page...")
	seasons, err := o.fetchSeasons(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to fetch seasons", zap.Error(err))
		return fmt.Errorf("fetch seasons: %w", err)
	}

	logger.Info(ctx, "[STEP 1] Seasons loaded",
		zap.Int("total_count", len(seasons)),
	)

	for i, s := range seasons {
		logger.Debug(ctx, fmt.Sprintf("  Season %d: %s (%s)", i+1, s.FullName, s.URL))
	}

	if len(seasons) == 0 {
		logger.Warn(ctx, "No seasons found, exiting")
		return nil
	}

	// Фильтруем сезоны если включен тестовый режим
	seasons = o.filterSeasons(ctx, seasons)

	logger.Info(ctx, "[STEP 2] Processing seasons...",
		zap.Int("seasons_to_process", len(seasons)),
	)

	// 2. Обрабатываем сезоны
	stats := o.processAllSeasons(ctx, seasons)

	elapsed := time.Since(start)
	logger.Info(ctx, "========================================")
	logger.Info(ctx, "MIHF Parser completed",
		zap.Duration("elapsed", elapsed),
		zap.Int("seasons_processed", len(seasons)),
		zap.Int("tournaments_processed", stats.tournaments),
		zap.Int("teams_processed", stats.teams),
		zap.Int("players_saved", stats.players),
	)
	logger.Info(ctx, "========================================")

	return nil
}

func (o *Orchestrator) filterSeasons(ctx context.Context, seasons []dto.SeasonDTO) []dto.SeasonDTO {
	// Если указан конкретный тестовый сезон
	if testSeason := o.config.TestSeason(); testSeason != "" {
		logger.Info(ctx, "Test mode: filtering to specific season", zap.String("test_season", testSeason))
		for _, s := range seasons {
			if s.Year == testSeason || s.FullName == testSeason {
				logger.Info(ctx, "Found test season", zap.String("season", s.FullName))
				return []dto.SeasonDTO{s}
			}
		}
		logger.Warn(ctx, "Test season not found, using all seasons", zap.String("requested", testSeason))
		return seasons
	}

	// Если указано максимальное количество сезонов
	if maxSeasons := o.config.MaxSeasons(); maxSeasons > 0 && len(seasons) > maxSeasons {
		logger.Info(ctx, "Limiting seasons",
			zap.Int("original", len(seasons)),
			zap.Int("limit", maxSeasons),
		)
		return seasons[:maxSeasons]
	}

	return seasons
}

type runStats struct {
	tournaments int
	teams       int
	players     int
}

func (o *Orchestrator) fetchSeasons(ctx context.Context) ([]dto.SeasonDTO, error) {
	logger.Debug(ctx, "Fetching main page...")
	html, err := o.client.Get("/")
	if err != nil {
		return nil, fmt.Errorf("get main page: %w", err)
	}
	logger.Debug(ctx, "Main page fetched", zap.Int("html_size", len(html)))

	seasons, err := parsing.ParseSeasons(html)
	if err != nil {
		return nil, fmt.Errorf("parse seasons: %w", err)
	}

	return seasons, nil
}

func (o *Orchestrator) processAllSeasons(ctx context.Context, seasons []dto.SeasonDTO) runStats {
	var (
		stats runStats
		mu    sync.Mutex
	)

	seasonCh := make(chan dto.SeasonDTO, len(seasons))
	for _, s := range seasons {
		seasonCh <- s
	}
	close(seasonCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.SeasonWorkers(); i++ {
		workerID := i + 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info(ctx, "Season worker started", zap.Int("worker_id", workerID))

			for season := range seasonCh {
				select {
				case <-ctx.Done():
					logger.Warn(ctx, "Context cancelled, stopping worker", zap.Int("worker_id", workerID))
					return
				default:
				}

				s := o.processSeasonSafe(ctx, season)
				mu.Lock()
				stats.tournaments += s.tournaments
				stats.teams += s.teams
				stats.players += s.players
				mu.Unlock()
			}

			logger.Info(ctx, "Season worker finished", zap.Int("worker_id", workerID))
		}()
	}

	wg.Wait()
	return stats
}
