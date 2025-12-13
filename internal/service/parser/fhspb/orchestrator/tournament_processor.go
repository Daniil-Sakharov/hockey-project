package orchestrator

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// processTournamentSafe обрабатывает турнир с логированием ошибок
func (o *Orchestrator) processTournamentSafe(ctx context.Context, t dto.TournamentDTO) (int, int) {
	logger.Info(ctx, "🏆 Processing tournament",
		zap.Int("id", t.ID),
		zap.String("name", t.Name),
	)

	if err := o.saveTournament(ctx, t); err != nil {
		logger.Debug(ctx, "⚠️ Save tournament failed", zap.Error(err))
	}

	teams, players, err := o.processTournament(ctx, t)
	if err != nil {
		logger.Error(ctx, "❌ Tournament failed", zap.Int("id", t.ID), zap.Error(err))
		return 0, 0
	}

	logger.Info(ctx, "✅ Tournament done",
		zap.Int("id", t.ID),
		zap.Int("teams", teams),
		zap.Int("players", players),
	)

	return teams, players
}

func (o *Orchestrator) saveTournament(ctx context.Context, t dto.TournamentDTO) error {
	return o.tournamentRepo.Create(ctx, &tournament.Tournament{
		ID:        strconv.Itoa(t.ID),
		URL:       fmt.Sprintf("fhspb://tournament/%d", t.ID),
		Name:      t.Name,
		Domain:    "fhspb.ru",
		Season:    t.Season,
		StartDate: t.StartDate,
		EndDate:   t.EndDate,
		IsEnded:   t.IsEnded,
		CreatedAt: time.Now(),
	})
}

func (o *Orchestrator) processTournament(ctx context.Context, t dto.TournamentDTO) (int, int, error) {
	teams, err := o.client.GetTeamsByTournament(t.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("get teams: %w", err)
	}

	logger.Info(ctx, "👥 Teams found", zap.Int("count", len(teams)))

	// Воркер пул для команд
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
	for i := 0; i < o.config.TeamWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for team := range teamCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				players := o.processTeamSafe(ctx, t.ID, team)
				mu.Lock()
				totalPlayers += players
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return len(teams), totalPlayers, nil
}
