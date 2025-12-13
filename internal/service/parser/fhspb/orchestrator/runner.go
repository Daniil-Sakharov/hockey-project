package orchestrator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/client/fhspb/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает полный процесс парсинга
func (o *Orchestrator) Run(ctx context.Context) error {
	start := time.Now()
	logger.Info(ctx, "🚀 Starting FHSPB parser", zap.Int("max_birth_year", o.config.MaxBirthYear))

	tournaments, err := o.client.GetTournamentsByBirthYear(o.config.MaxBirthYear)
	if err != nil {
		return fmt.Errorf("get tournaments: %w", err)
	}

	logger.Info(ctx, "📋 Tournaments loaded", zap.Int("count", len(tournaments)))

	if len(tournaments) == 0 {
		logger.Warn(ctx, "⚠️ No tournaments found")
		return nil
	}

	totalTeams, totalPlayers := o.processAllTournaments(ctx, tournaments)

	logger.Info(ctx, "✅ FHSPB parser completed",
		zap.Duration("elapsed", time.Since(start)),
		zap.Int("tournaments", len(tournaments)),
		zap.Int("teams", totalTeams),
		zap.Int("players_saved", totalPlayers),
	)

	return nil
}

func (o *Orchestrator) processAllTournaments(ctx context.Context, tournaments []dto.TournamentDTO) (int, int) {
	var (
		totalTeams   int
		totalPlayers int
		mu           sync.Mutex
	)

	tournamentCh := make(chan dto.TournamentDTO, len(tournaments))
	for _, t := range tournaments {
		tournamentCh <- t
	}
	close(tournamentCh)

	var wg sync.WaitGroup
	for i := 0; i < o.config.TournamentWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for tournament := range tournamentCh {
				select {
				case <-ctx.Done():
					return
				default:
				}

				teams, players := o.processTournamentSafe(ctx, tournament)
				mu.Lock()
				totalTeams += teams
				totalPlayers += players
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return totalTeams, totalPlayers
}
