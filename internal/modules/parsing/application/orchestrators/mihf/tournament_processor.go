package mihf

import (
	"context"
	"fmt"
	"sync"

	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processTournamentSafe(ctx context.Context, path dto.TournamentPathDTO) runStats {
	logger.Info(ctx, "[TOURNAMENT] Starting",
		zap.String("season", path.SeasonYear),
		zap.Int("birth_year", path.BirthYear),
		zap.String("group", path.GroupName),
		zap.String("tournament_id", path.TournamentID),
		zap.String("sub_id", path.SubID),
	)

	tournamentID, err := o.saveTournament(ctx, path)
	if err != nil {
		logger.Error(ctx, "[TOURNAMENT] Failed to save",
			zap.String("tournament_id", path.TournamentID),
			zap.Error(err),
		)
		return runStats{}
	}

	logger.Info(ctx, "[TOURNAMENT] Saved to DB",
		zap.String("db_id", tournamentID),
	)

	stats, err := o.processTournament(ctx, path, tournamentID)
	if err != nil {
		logger.Error(ctx, "[TOURNAMENT] Failed to process",
			zap.String("tournament_id", tournamentID),
			zap.Error(err),
		)
		return runStats{}
	}

	logger.Info(ctx, "[TOURNAMENT] Completed",
		zap.String("db_id", tournamentID),
		zap.Int("teams", stats.teams),
		zap.Int("players", stats.players),
	)

	return stats
}

func (o *Orchestrator) saveTournament(ctx context.Context, path dto.TournamentPathDTO) (string, error) {
	externalID := fmt.Sprintf("%s-%s-%s", path.TournamentID, path.SubID, path.GroupID)
	season := path.SeasonYear + "-" + nextYear(path.SeasonYear)

	tournament := &mihfrepo.Tournament{
		ExternalID: externalID,
		Name:       path.GroupName,
		Season:     &season,
	}

	if path.BirthYear > 0 {
		tournament.BirthYear = &path.BirthYear
	}
	if path.GroupName != "" {
		tournament.GroupName = &path.GroupName
	}

	url := path.ScoreboardURL()
	tournament.URL = &url

	logger.Debug(ctx, "[TOURNAMENT] Upserting",
		zap.String("external_id", externalID),
		zap.String("name", path.GroupName),
		zap.String("season", season),
	)

	return o.tournamentRepo.Upsert(ctx, tournament)
}

func nextYear(year string) string {
	var y int
	fmt.Sscanf(year, "%d", &y)
	return fmt.Sprintf("%d", y+1)
}

func (o *Orchestrator) processTournament(ctx context.Context, path dto.TournamentPathDTO, tournamentID string) (runStats, error) {
	// Получаем scoreboard
	scoreboardURL := path.ScoreboardURL()
	logger.Info(ctx, "[TOURNAMENT] Fetching scoreboard",
		zap.String("url", scoreboardURL),
	)

	html, err := o.client.Get(scoreboardURL)
	if err != nil {
		return runStats{}, fmt.Errorf("get scoreboard: %w", err)
	}

	logger.Debug(ctx, "[TOURNAMENT] Scoreboard fetched",
		zap.Int("html_size", len(html)),
	)

	// Диагностика: показываем первые 1000 символов HTML
	preview := string(html)
	if len(preview) > 1000 {
		preview = preview[:1000]
	}
	logger.Debug(ctx, "[TOURNAMENT] HTML preview",
		zap.String("preview", preview),
	)

	teams, err := parsing.ParseScoreboard(html)
	if err != nil {
		return runStats{}, fmt.Errorf("parse scoreboard: %w", err)
	}

	logger.Info(ctx, "[TOURNAMENT] Teams parsed from scoreboard",
		zap.String("tournament_id", tournamentID),
		zap.Int("teams_count", len(teams)),
	)

	for _, t := range teams {
		logger.Debug(ctx, fmt.Sprintf("  Team: %s (id=%s, games=%d, points=%d)", t.Name, t.ID, t.Games, t.Points))
	}

	if len(teams) == 0 {
		logger.Warn(ctx, "[TOURNAMENT] No teams found in scoreboard")
		return runStats{}, nil
	}

	return o.processAllTeams(ctx, path, tournamentID, teams), nil
}

func (o *Orchestrator) processAllTeams(ctx context.Context, path dto.TournamentPathDTO, tournamentID string, teams []dto.TeamDTO) runStats {
	var (
		stats runStats
		mu    sync.Mutex
	)

	type teamJob struct {
		team dto.TeamDTO
		path dto.TournamentPathDTO
	}

	teamCh := make(chan teamJob, len(teams))
	for _, t := range teams {
		teamCh <- teamJob{team: t, path: path}
	}
	close(teamCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.TeamWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range teamCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				players := o.processTeamSafe(ctx, job.path, tournamentID, job.team)
				mu.Lock()
				stats.teams++
				stats.players += players
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return stats
}
