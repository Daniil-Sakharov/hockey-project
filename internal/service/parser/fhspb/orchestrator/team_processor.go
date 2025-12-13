package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processTeamSafe(ctx context.Context, tournamentID int, t dto.TeamDTO) int {
	if err := o.saveTeam(ctx, t); err != nil {
		logger.Debug(ctx, "⚠️ Save team failed", zap.Error(err))
	}

	players, err := o.processTeam(ctx, tournamentID, t)
	if err != nil {
		logger.Error(ctx, "❌ Team failed", zap.String("name", t.Name), zap.Error(err))
		return 0
	}

	logger.Info(ctx, "✅ Team done", zap.String("name", t.Name), zap.Int("players", players))
	return players
}

func (o *Orchestrator) saveTeam(ctx context.Context, t dto.TeamDTO) error {
	_, err := o.teamRepo.Upsert(ctx, &team.Team{
		ID:        t.ID,
		URL:       fmt.Sprintf("fhspb://team/%s", t.ID),
		Name:      t.Name,
		CreatedAt: time.Now(),
	})
	return err
}

func (o *Orchestrator) processTeam(ctx context.Context, tournamentID int, t dto.TeamDTO) (int, error) {
	playerURLs, err := o.client.GetPlayerURLsFromTeam(tournamentID, t.ID)
	if err != nil {
		return 0, fmt.Errorf("get player urls: %w", err)
	}

	// Воркер пул для игроков
	var (
		savedCount int
		mu         sync.Mutex
	)

	playerCh := make(chan dto.PlayerURLDTO, len(playerURLs))
	for _, p := range playerURLs {
		playerCh <- p
	}
	close(playerCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.PlayerWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pURL := range playerCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if o.processPlayerSafe(ctx, pURL) {
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
