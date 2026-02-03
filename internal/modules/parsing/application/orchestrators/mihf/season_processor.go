package mihf

import (
	"context"
	"fmt"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processSeasonSafe(ctx context.Context, season dto.SeasonDTO) runStats {
	logger.Info(ctx, "----------------------------------------")
	logger.Info(ctx, "[SEASON] Starting processing",
		zap.String("year", season.Year),
		zap.String("name", season.FullName),
		zap.String("url", season.URL),
	)

	stats, err := o.processSeason(ctx, season)
	if err != nil {
		logger.Error(ctx, "[SEASON] Failed",
			zap.String("year", season.Year),
			zap.Error(err),
		)
		return runStats{}
	}

	logger.Info(ctx, "[SEASON] Completed",
		zap.String("year", season.Year),
		zap.Int("tournaments", stats.tournaments),
		zap.Int("teams", stats.teams),
		zap.Int("players", stats.players),
	)
	logger.Info(ctx, "----------------------------------------")

	return stats
}

func (o *Orchestrator) processSeason(ctx context.Context, season dto.SeasonDTO) (runStats, error) {
	// Получаем группы турниров
	logger.Info(ctx, "[SEASON] Fetching groups...", zap.String("url", season.URL))
	groups, err := o.fetchGroups(ctx, season)
	if err != nil {
		return runStats{}, fmt.Errorf("fetch groups: %w", err)
	}

	logger.Info(ctx, "[SEASON] Groups loaded",
		zap.String("season", season.Year),
		zap.Int("total_groups", len(groups)),
	)

	for _, g := range groups {
		logger.Debug(ctx, fmt.Sprintf("  Group: %s (birth_year=%d, url=%s)", g.Name, g.BirthYear, g.URL))
	}

	// Собираем все tournament paths через двухэтапный процесс
	var allPaths []dto.TournamentPathDTO

	for _, group := range groups {
		logger.Info(ctx, "[SEASON] Processing group",
			zap.String("group_id", group.ID),
			zap.String("group_name", group.Name),
		)

		// Шаг 1: Получаем турниры (по году рождения) из страницы группы
		tournaments, err := o.fetchTournaments(ctx, season, group)
		if err != nil {
			logger.Warn(ctx, "[SEASON] Failed to fetch tournaments",
				zap.String("group_id", group.ID),
				zap.Error(err),
			)
			continue
		}

		logger.Info(ctx, "[SEASON] Tournaments loaded",
			zap.String("group_id", group.ID),
			zap.Int("count", len(tournaments)),
		)

		// Фильтруем по году рождения
		var filteredTournaments []dto.TournamentDTO
		for _, t := range tournaments {
			if t.BirthYear < o.config.MinBirthYear() {
				continue
			}
			if max := o.config.MaxBirthYear(); max > 0 && t.BirthYear > max {
				continue
			}
			filteredTournaments = append(filteredTournaments, t)
			logger.Debug(ctx, fmt.Sprintf("  Tournament: %s (birth_year=%d)", t.Name, t.BirthYear))
		}

		logger.Info(ctx, "[SEASON] Tournaments filtered",
			zap.String("group_id", group.ID),
			zap.Int("before", len(tournaments)),
			zap.Int("after", len(filteredTournaments)),
			zap.Int("min_birth_year", o.config.MinBirthYear()),
			zap.Int("max_birth_year", o.config.MaxBirthYear()),
		)

		// Шаг 2: Для каждого турнира получаем подтурниры (Группа А, Б, В)
		for _, tournament := range filteredTournaments {
			subTournaments, err := o.fetchSubTournaments(ctx, season, group, tournament)
			if err != nil {
				logger.Warn(ctx, "[SEASON] Failed to fetch sub-tournaments",
					zap.String("tournament_id", tournament.ID),
					zap.Error(err),
				)
				continue
			}

			logger.Debug(ctx, "[SEASON] Sub-tournaments loaded",
				zap.String("tournament_id", tournament.ID),
				zap.Int("count", len(subTournaments)),
			)

			// Строим TournamentPathDTO для каждого подтурнира
			for _, sub := range subTournaments {
				path := parsing.BuildTournamentPath(season.Year, tournament, sub)
				allPaths = append(allPaths, path)
				logger.Debug(ctx, fmt.Sprintf("  Path: %s -> %s", tournament.Name, sub.Name))
			}
		}
	}

	logger.Info(ctx, "[SEASON] Tournament paths collected",
		zap.String("season", season.Year),
		zap.Int("total_paths", len(allPaths)),
	)

	if len(allPaths) == 0 {
		logger.Warn(ctx, "[SEASON] No tournament paths found", zap.String("season", season.Year))
		return runStats{}, nil
	}

	return o.processAllTournaments(ctx, allPaths), nil
}

func (o *Orchestrator) fetchGroups(ctx context.Context, season dto.SeasonDTO) ([]dto.GroupDTO, error) {
	html, err := o.client.Get(season.URL)
	if err != nil {
		return nil, fmt.Errorf("get season page: %w", err)
	}
	logger.Debug(ctx, "Season page fetched", zap.Int("html_size", len(html)))

	groups, err := parsing.ParseGroups(html, season.Year)
	if err != nil {
		return nil, fmt.Errorf("parse groups: %w", err)
	}

	return groups, nil
}

func (o *Orchestrator) fetchTournaments(ctx context.Context, season dto.SeasonDTO, group dto.GroupDTO) ([]dto.TournamentDTO, error) {
	html, err := o.client.Get(group.URL)
	if err != nil {
		return nil, fmt.Errorf("get group page: %w", err)
	}

	logger.Debug(ctx, "[FETCH] Group page fetched",
		zap.String("group_id", group.ID),
		zap.Int("html_size", len(html)),
	)

	tournaments, err := parsing.ParseTournaments(html, season.Year, group.ID)
	if err != nil {
		return nil, fmt.Errorf("parse tournaments: %w", err)
	}

	return tournaments, nil
}

func (o *Orchestrator) fetchSubTournaments(ctx context.Context, season dto.SeasonDTO, group dto.GroupDTO, tournament dto.TournamentDTO) ([]dto.SubTournamentDTO, error) {
	html, err := o.client.Get(tournament.URL)
	if err != nil {
		return nil, fmt.Errorf("get tournament page: %w", err)
	}

	logger.Debug(ctx, "[FETCH] Tournament page fetched",
		zap.String("tournament_id", tournament.ID),
		zap.Int("html_size", len(html)),
	)

	subTournaments, err := parsing.ParseSubTournaments(html, season.Year, group.ID, tournament.ID)
	if err != nil {
		return nil, fmt.Errorf("parse sub-tournaments: %w", err)
	}

	return subTournaments, nil
}

func (o *Orchestrator) processAllTournaments(ctx context.Context, paths []dto.TournamentPathDTO) runStats {
	var (
		stats runStats
		mu    sync.Mutex
	)

	pathCh := make(chan dto.TournamentPathDTO, len(paths))
	for _, p := range paths {
		pathCh <- p
	}
	close(pathCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.TournamentWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				s := o.processTournamentSafe(ctx, path)
				mu.Lock()
				stats.tournaments++
				stats.teams += s.teams
				stats.players += s.players
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return stats
}
