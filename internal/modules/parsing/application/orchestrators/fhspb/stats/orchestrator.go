package stats

import (
	"context"
	"strconv"
	"strings"

	fhspbrepo "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories/fhspb"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Tournament алиас для удобства
type Tournament = fhspbrepo.Tournament

type Orchestrator struct {
	deps            Dependencies
	playerProcessor *PlayerProcessor
	goalieProcessor *GoalieProcessor
}

func NewOrchestrator(deps Dependencies) *Orchestrator {
	return &Orchestrator{
		deps:            deps,
		playerProcessor: NewPlayerProcessor(deps),
		goalieProcessor: NewGoalieProcessor(deps),
	}
}

func (o *Orchestrator) Run(ctx context.Context) error {
	logger.Info(ctx, "starting statistics parser")

	// Получаем все турниры SPB
	tournaments, err := o.deps.TournamentRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	// Конвертируем в указатели
	tournamentPtrs := make([]*Tournament, len(tournaments))
	for i := range tournaments {
		tournamentPtrs[i] = &tournaments[i]
	}

	return o.RunForTournaments(ctx, tournamentPtrs)
}

// RunForTournaments парсит статистику для указанных турниров
func (o *Orchestrator) RunForTournaments(ctx context.Context, tournaments []*Tournament) error {
	logger.Info(ctx, "processing tournaments", zap.Int("count", len(tournaments)))

	for _, t := range tournaments {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// external_id хранится как string, нужно конвертировать в int для API
		externalID, err := strconv.Atoi(t.ExternalID)
		if err != nil {
			logger.Error(ctx, "invalid external_id", zap.String("external_id", t.ExternalID))
			continue
		}

		logger.Info(ctx, "processing tournament statistics",
			zap.String("tournament_id", t.ID),
			zap.String("name", t.Name),
		)

		// Парсим статистику полевых игроков
		if err := o.playerProcessor.Process(ctx, externalID, t.ID); err != nil {
			// 500 - серверная ошибка fhspb.ru, пропускаем
			if strings.Contains(err.Error(), "500") {
				logger.Warn(ctx, "server error for player statistics",
					zap.String("tournament_id", t.ID),
					zap.Error(err),
				)
			} else {
				logger.Error(ctx, "failed to process player statistics",
					zap.String("tournament_id", t.ID),
					zap.Error(err),
				)
			}
			continue
		}

		// Парсим статистику вратарей
		if err := o.goalieProcessor.Process(ctx, externalID, t.ID); err != nil {
			// 500 - серверная ошибка fhspb.ru, пропускаем
			if strings.Contains(err.Error(), "500") {
				logger.Warn(ctx, "server error for goalie statistics",
					zap.String("tournament_id", t.ID),
					zap.Error(err),
				)
			} else {
				logger.Error(ctx, "failed to process goalie statistics",
					zap.String("tournament_id", t.ID),
					zap.Error(err),
				)
			}
			continue
		}
	}

	logger.Info(ctx, "statistics parser completed")
	return nil
}
