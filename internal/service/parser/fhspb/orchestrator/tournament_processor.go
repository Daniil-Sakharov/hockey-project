package orchestrator

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	fhspbRepo "github.com/Daniil-Sakharov/HockeyProject/internal/repository/postgres/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// processTournamentSafe обрабатывает турнир с логированием ошибок
func (o *Orchestrator) processTournamentSafe(ctx context.Context, t dto.TournamentDTO) (int, int) {
	logger.Info(ctx, "🏆 Processing tournament",
		zap.Int("id", t.ID),
		zap.String("name", t.Name),
	)

	tournamentID, err := o.saveTournament(ctx, t)
	if err != nil {
		logger.Debug(ctx, "⚠️ Save tournament failed", zap.Error(err))
		return 0, 0
	}

	teams, players, err := o.processTournament(ctx, t, tournamentID)
	if err != nil {
		// 404 - турнир удалён или недоступен, это нормально
		if strings.Contains(err.Error(), "404") {
			logger.Warn(ctx, "⚠️ Tournament not found", zap.Int("id", t.ID), zap.Error(err))
		} else {
			logger.Error(ctx, "❌ Tournament failed", zap.Int("id", t.ID), zap.Error(err))
		}
		return 0, 0
	}

	logger.Info(ctx, "✅ Tournament done",
		zap.Int("id", t.ID),
		zap.Int("teams", teams),
		zap.Int("players", players),
	)

	return teams, players
}

func (o *Orchestrator) saveTournament(ctx context.Context, t dto.TournamentDTO) (string, error) {
	// Формируем URL турнира
	tournamentURL := fmt.Sprintf("https://www.fhspb.ru/Tournament?TournamentID=%d", t.ID)
	domain := "https://www.fhspb.ru"

	tournament := &fhspbRepo.Tournament{
		ExternalID: strconv.Itoa(t.ID),
		URL:        &tournamentURL,
		Name:       t.Name,
		Domain:     &domain,
		Season:     &t.Season,
		StartDate:  t.StartDate,
		EndDate:    t.EndDate,
		IsEnded:    t.IsEnded,
	}
	if t.BirthYear > 0 {
		tournament.BirthYear = &t.BirthYear
	}
	if t.GroupName != "" {
		tournament.GroupName = &t.GroupName
	}
	return o.tournamentRepo.Upsert(ctx, tournament)
}

func (o *Orchestrator) processTournament(ctx context.Context, t dto.TournamentDTO, tournamentID string) (int, int, error) {
	teams, err := o.client.GetTeamsByTournament(t.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("get teams: %w", err)
	}

	logger.Info(ctx, "👥 Teams found", zap.Int("count", len(teams)))

	var (
		totalPlayers int
		mu           sync.Mutex
	)

	teamCh := make(chan dto.TeamDTO, len(teams))
	for _, team := range teams {
		teamCh <- team
	}
	close(teamCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.TeamWorkers(); i++ {
		workerID := i + 1
		wg.Add(1)
		go func() {
			defer wg.Done()
			for team := range teamCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				logger.Debug(ctx, "👷 Team worker processing", zap.Int("worker_id", workerID), zap.String("team", team.Name))
				players := o.processTeamSafe(ctx, tournamentID, team)
				mu.Lock()
				totalPlayers += players
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return len(teams), totalPlayers, nil
}
