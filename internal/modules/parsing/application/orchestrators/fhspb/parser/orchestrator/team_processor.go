package parser

import (
	"context"
	"fmt"
	"sync"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

func (o *Orchestrator) processTeamSafe(ctx context.Context, tournamentID string, t dto.TeamDTO) int {
	teamID, err := o.saveTeam(ctx, tournamentID, t)
	if err != nil {
		logger.Warn(ctx, "⚠️ Save team failed", zap.String("name", t.Name), zap.Error(err))
		return 0
	}

	players, err := o.processTeam(ctx, t.TournamentID, teamID, tournamentID, t)
	if err != nil {
		logger.Error(ctx, "❌ Team failed", zap.String("name", t.Name), zap.Error(err))
		return 0
	}

	logger.Info(ctx, "✅ Team done", zap.String("name", t.Name), zap.Int("players", players))
	return players
}

func (o *Orchestrator) saveTeam(ctx context.Context, tournamentID string, t dto.TeamDTO) (string, error) {
	url := fmt.Sprintf("https://www.fhspb.ru/Team?TournamentID=%d&TeamID=%s", t.TournamentID, t.ID)

	team := &fhspb.Team{
		ExternalID:   t.ID,
		TournamentID: tournamentID,
		Name:         t.Name,
		URL:          &url,
	}

	if t.City != "" {
		team.City = &t.City
	}

	return o.teamRepo.Upsert(ctx, team)
}

func (o *Orchestrator) processTeam(ctx context.Context, extTournamentID int, teamID, tournamentID string, t dto.TeamDTO) (int, error) {
	playerURLs, err := o.client.GetPlayerURLsFromTeam(extTournamentID, t.ID)
	if err != nil {
		return 0, fmt.Errorf("get player urls: %w", err)
	}

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
	for i := 0; i < o.config.PlayerWorkers(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pURL := range playerCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if o.processPlayerSafe(ctx, teamID, tournamentID, pURL) {
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
