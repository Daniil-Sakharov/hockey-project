package calendar

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает полный цикл парсинга календаря
func (o *Orchestrator) Run(ctx context.Context) error {
	logger.Info(ctx, "🏒 Starting Junior calendar parsing...")

	baseTournaments, err := o.tournamentRepo.GetBySource(ctx, Source)
	if err != nil {
		return err
	}

	if o.config.MaxTournaments() > 0 && len(baseTournaments) > o.config.MaxTournaments() {
		logger.Info(ctx, "Applying global tournament limit",
			zap.Int("limit", o.config.MaxTournaments()),
			zap.Int("total", len(baseTournaments)))
		baseTournaments = baseTournaments[:o.config.MaxTournaments()]
	}

	logger.Info(ctx, "Found base tournaments for calendar parsing",
		zap.Int("count", len(baseTournaments)))

	// Параллельная обработка турниров
	o.processTournamentsParallel(ctx, baseTournaments)

	// Параллельный парсинг деталей матчей
	if err := o.processUnparsedGames(ctx); err != nil {
		logger.Error(ctx, "Failed to process unparsed games", zap.Error(err))
	}

	logger.Info(ctx, "✅ Junior calendar parsing completed")
	return nil
}

// processTournamentsParallel обрабатывает турниры параллельно через worker pool
func (o *Orchestrator) processTournamentsParallel(ctx context.Context, tournaments []*entities.Tournament) {
	workers := o.config.TournamentWorkers()
	if workers <= 0 {
		workers = 3
	}
	if workers > len(tournaments) {
		workers = len(tournaments)
	}

	logger.Info(ctx, "Starting tournament worker pool",
		zap.Int("workers", workers),
		zap.Int("tournaments", len(tournaments)))

	ch := make(chan tournamentTask, len(tournaments))
	for idx, t := range tournaments {
		ch <- tournamentTask{index: idx, total: len(tournaments), tournament: t}
	}
	close(ch)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch {
				if ctx.Err() != nil {
					return
				}
				o.processSingleTournament(ctx, task)
			}
		}()
	}
	wg.Wait()
}

type tournamentTask struct {
	index      int
	total      int
	tournament *entities.Tournament
}

func (o *Orchestrator) processSingleTournament(ctx context.Context, task tournamentTask) {
	logger.Info(ctx, "Processing tournament",
		zap.Int("index", task.index+1),
		zap.Int("total", task.total),
		zap.String("name", task.tournament.Name))

	if err := o.processCalendarWithFilters(ctx, task.tournament); err != nil {
		logger.Error(ctx, "Failed to process tournament calendar",
			zap.String("tournament", task.tournament.Name),
			zap.Error(err))
		return
	}

	if err := o.processStandingsWithFilters(ctx, task.tournament); err != nil {
		logger.Warn(ctx, "Failed to parse standings", zap.Error(err))
	}
}

// processUnparsedGames парсит детали завершённых матчей параллельно
func (o *Orchestrator) processUnparsedGames(ctx context.Context) error {
	if !o.config.ParseProtocol() {
		return nil
	}

	matches, err := o.matchRepo.GetUnparsedFinished(ctx, Source, 0)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return nil
	}

	workers := o.config.GameWorkers()
	if workers <= 0 {
		workers = 5
	}
	if workers > len(matches) {
		workers = len(matches)
	}

	logger.Info(ctx, "Processing unparsed finished games",
		zap.Int("count", len(matches)),
		zap.Int("workers", workers))

	ch := make(chan *entities.Match, len(matches))
	for _, m := range matches {
		ch <- m
	}
	close(ch)

	var processed int64
	var failed int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for m := range ch {
				if ctx.Err() != nil {
					return
				}
				if err := o.processGame(ctx, m.ID, m.ExternalID); err != nil {
					logger.Error(ctx, "Failed to process game",
						zap.String("match_id", m.ID),
						zap.Error(err))
					atomic.AddInt64(&failed, 1)
					continue
				}
				atomic.AddInt64(&processed, 1)
			}
		}()
	}
	wg.Wait()

	logger.Info(ctx, "Finished processing games",
		zap.Int64("processed", processed),
		zap.Int64("failed", failed))

	return nil
}
