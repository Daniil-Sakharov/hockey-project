package fhmoscow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhmoscow/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает полный процесс парсинга FHMoscow
func (o *Orchestrator) Run(ctx context.Context) error {
	start := time.Now()

	logger.Info(ctx, "========================================")
	logger.Info(ctx, "FHMoscow Parser starting",
		zap.Int("min_birth_year", o.config.MinBirthYear()),
		zap.Int("max_seasons", o.config.MaxSeasons()),
		zap.String("test_season", o.config.TestSeason()),
		zap.Int("tournament_workers", o.config.TournamentWorkers()),
		zap.Int("team_workers", o.config.TeamWorkers()),
		zap.Int("player_workers", o.config.PlayerWorkers()),
	)
	logger.Info(ctx, "========================================")

	// 1. Получаем сезоны
	logger.Info(ctx, "[STEP 1] Fetching seasons from API...")
	seasons, err := o.fetchSeasons(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to fetch seasons", zap.Error(err))
		return fmt.Errorf("fetch seasons: %w", err)
	}

	logger.Info(ctx, "[STEP 1] Seasons loaded", zap.Int("total_count", len(seasons)))

	// 2. Получаем все турниры
	logger.Info(ctx, "[STEP 2] Fetching tournaments from API...")
	tournaments, err := o.fetchTournaments(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to fetch tournaments", zap.Error(err))
		return fmt.Errorf("fetch tournaments: %w", err)
	}

	// Фильтруем турниры по году рождения
	tournaments = o.filterTournaments(ctx, tournaments)
	logger.Info(ctx, "[STEP 2] Tournaments filtered",
		zap.Int("total_count", len(tournaments)),
		zap.Int("min_birth_year", o.config.MinBirthYear()),
	)

	if len(tournaments) == 0 {
		logger.Warn(ctx, "No tournaments found after filtering, exiting")
		return nil
	}

	// 3. Обрабатываем турниры
	logger.Info(ctx, "[STEP 3] Processing tournaments...")
	stats := o.processAllTournaments(ctx, seasons, tournaments)

	// 4. Сканируем игроков по ID (если включено)
	var scanStats scanStats
	if o.config.ScanPlayers() {
		scanStats = o.scanPlayers(ctx)
		stats.players += int(scanStats.saved)
	} else {
		logger.Info(ctx, "[STEP 4] Player scanning disabled, skipping...")
	}

	elapsed := time.Since(start)
	logger.Info(ctx, "========================================")
	logger.Info(ctx, "FHMoscow Parser completed",
		zap.Duration("elapsed", elapsed),
		zap.Int("tournaments_processed", stats.tournaments),
		zap.Int("groups_processed", stats.groups),
		zap.Int("teams_processed", stats.teams),
		zap.Int("players_saved", stats.players),
		zap.Int64("players_scanned", scanStats.scanned),
	)
	logger.Info(ctx, "========================================")

	return nil
}

type runStats struct {
	tournaments int
	groups      int
	teams       int
	players     int
}

func (o *Orchestrator) fetchSeasons(ctx context.Context) ([]dto.SeasonDTO, error) {
	data, err := o.client.GetAPI("/api/filter/season")
	if err != nil {
		return nil, fmt.Errorf("get seasons: %w", err)
	}

	seasons, err := parsing.ParseSeasons(data)
	if err != nil {
		return nil, fmt.Errorf("parse seasons: %w", err)
	}

	return seasons, nil
}

func (o *Orchestrator) fetchTournaments(ctx context.Context) ([]dto.TournamentDTO, error) {
	// POST /api/filter/tournament с пустым телом возвращает все турниры
	data, err := o.client.PostAPI("/api/filter/tournament", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("get tournaments: %w", err)
	}

	tournaments, err := parsing.ParseTournaments(data)
	if err != nil {
		return nil, fmt.Errorf("parse tournaments: %w", err)
	}

	return tournaments, nil
}

func (o *Orchestrator) filterTournaments(ctx context.Context, tournaments []dto.TournamentDTO) []dto.TournamentDTO {
	var filtered []dto.TournamentDTO
	minYear := o.config.MinBirthYear()

	for _, t := range tournaments {
		birthYear := t.ParseBirthYear()
		if birthYear >= minYear {
			filtered = append(filtered, t)
			logger.Debug(ctx, "Tournament included",
				zap.Int("id", t.ID),
				zap.String("name", t.Name),
				zap.Int("birth_year", birthYear),
			)
		} else if birthYear > 0 {
			logger.Debug(ctx, "Tournament excluded (birth_year < min)",
				zap.Int("id", t.ID),
				zap.String("name", t.Name),
				zap.Int("birth_year", birthYear),
				zap.Int("min_birth_year", minYear),
			)
		}
	}

	return filtered
}

func (o *Orchestrator) processAllTournaments(ctx context.Context, seasons []dto.SeasonDTO, tournaments []dto.TournamentDTO) runStats {
	var (
		stats runStats
		mu    sync.Mutex
	)

	// Находим текущий сезон
	var currentSeason dto.SeasonDTO
	for _, s := range seasons {
		if s.Current {
			currentSeason = s
			break
		}
	}
	if currentSeason.ID == 0 && len(seasons) > 0 {
		currentSeason = seasons[len(seasons)-1]
	}

	tournamentCh := make(chan dto.TournamentDTO, len(tournaments))
	for _, t := range tournaments {
		tournamentCh <- t
	}
	close(tournamentCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.TournamentWorkers(); i++ {
		workerID := i + 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info(ctx, "Tournament worker started", zap.Int("worker_id", workerID))

			for tournament := range tournamentCh {
				select {
				case <-ctx.Done():
					logger.Warn(ctx, "Context cancelled, stopping worker", zap.Int("worker_id", workerID))
					return
				default:
				}

				s := o.processTournamentSafe(ctx, currentSeason, tournament)
				mu.Lock()
				stats.tournaments++
				stats.groups += s.groups
				stats.teams += s.teams
				stats.players += s.players
				mu.Unlock()
			}

			logger.Info(ctx, "Tournament worker finished", zap.Int("worker_id", workerID))
		}()
	}

	wg.Wait()
	return stats
}
