package mihf

import (
	"context"
	"fmt"
	"strings"
	"sync"

	mihfrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/mihf"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/mihf/parsing"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processTeamSafe(ctx context.Context, path dto.TournamentPathDTO, tournamentID string, team dto.TeamDTO) int {
	logger.Info(ctx, "[TEAM] Starting",
		zap.String("name", team.Name),
		zap.String("team_id", team.ID),
		zap.String("tournament_id", tournamentID),
	)

	teamID, err := o.saveTeam(ctx, tournamentID, team)
	if err != nil {
		logger.Error(ctx, "[TEAM] Failed to save",
			zap.String("name", team.Name),
			zap.Error(err),
		)
		return 0
	}

	logger.Debug(ctx, "[TEAM] Saved to DB",
		zap.String("name", team.Name),
		zap.String("db_id", teamID),
	)

	players, err := o.processTeam(ctx, path, tournamentID, teamID, team)
	if err != nil {
		logger.Error(ctx, "[TEAM] Failed to process",
			zap.String("name", team.Name),
			zap.Error(err),
		)
		return 0
	}

	logger.Info(ctx, "[TEAM] Completed",
		zap.String("name", team.Name),
		zap.Int("players_saved", players),
	)

	return players
}

func (o *Orchestrator) saveTeam(ctx context.Context, tournamentID string, t dto.TeamDTO) (string, error) {
	team := &mihfrepo.Team{
		ExternalID:   t.ID,
		TournamentID: tournamentID,
		Name:         t.Name,
	}

	if t.ExternalURL != "" {
		team.URL = &t.ExternalURL
	}

	return o.teamRepo.Upsert(ctx, team)
}

func (o *Orchestrator) processTeam(ctx context.Context, path dto.TournamentPathDTO, tournamentID, teamID string, team dto.TeamDTO) (int, error) {
	// Получаем страницу команды со статистикой
	teamURL := path.TeamURL(team.ID)
	logger.Info(ctx, "[TEAM] Fetching team page",
		zap.String("name", team.Name),
		zap.String("url", teamURL),
	)

	html, err := o.client.Get(teamURL)
	if err != nil {
		return 0, fmt.Errorf("get team page: %w", err)
	}

	logger.Debug(ctx, "[TEAM] Team page fetched",
		zap.String("name", team.Name),
		zap.Int("html_size", len(html)),
	)

	// Диагностика: проверяем есть ли игроки в HTML
	hasPlayers := strings.Contains(string(html), "/players/info/")
	logger.Debug(ctx, "[TEAM] HTML check",
		zap.String("name", team.Name),
		zap.Bool("has_player_links", hasPlayers),
	)

	playerStats, goalieStats, err := parsing.ParseTeamStats(html)
	if err != nil {
		return 0, fmt.Errorf("parse team stats: %w", err)
	}

	logger.Info(ctx, "[TEAM] Stats parsed",
		zap.String("name", team.Name),
		zap.Int("field_players", len(playerStats)),
		zap.Int("goalies", len(goalieStats)),
	)

	// Логируем первых нескольких игроков для отладки
	for i, p := range playerStats {
		if i < 3 {
			logger.Debug(ctx, fmt.Sprintf("  Player: %s (id=%s, G=%d, A=%d, P=%d)", p.Name, p.ID, p.Goals, p.Assists, p.Points))
		}
	}
	for i, g := range goalieStats {
		if i < 2 {
			logger.Debug(ctx, fmt.Sprintf("  Goalie: %s (id=%s, GA=%d, SV%%=%.2f)", g.Name, g.ID, g.GoalsAgainst, g.SavePercentage))
		}
	}

	// Обрабатываем игроков
	var (
		savedCount int
		mu         sync.Mutex
	)

	type playerJob struct {
		playerStats *dto.PlayerStatsDTO
		goalieStats *dto.GoalieStatsDTO
	}

	allJobs := make([]playerJob, 0, len(playerStats)+len(goalieStats))
	for i := range playerStats {
		allJobs = append(allJobs, playerJob{playerStats: &playerStats[i]})
	}
	for i := range goalieStats {
		allJobs = append(allJobs, playerJob{goalieStats: &goalieStats[i]})
	}

	if len(allJobs) == 0 {
		logger.Warn(ctx, "[TEAM] No players found", zap.String("name", team.Name))
		return 0, nil
	}

	jobCh := make(chan playerJob, len(allJobs))
	for _, j := range allJobs {
		jobCh <- j
	}
	close(jobCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.PlayerWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				var success bool
				if job.playerStats != nil {
					success = o.processPlayerStatsSafe(ctx, tournamentID, teamID, path.BirthYear, *job.playerStats)
				} else if job.goalieStats != nil {
					success = o.processGoalieStatsSafe(ctx, tournamentID, teamID, path.BirthYear, *job.goalieStats)
				}

				if success {
					mu.Lock()
					savedCount++
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return savedCount, nil
}
