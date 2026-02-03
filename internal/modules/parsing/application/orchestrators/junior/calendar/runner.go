package calendar

import (
	"context"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает полный цикл парсинга календаря
func (o *Orchestrator) Run(ctx context.Context) error {
	logger.Info(ctx, "🏒 Starting Junior calendar parsing...")

	// Получаем все турниры из БД (под-турниры удалены, только реальные)
	baseTournaments, err := o.tournamentRepo.GetBySource(ctx, Source)
	if err != nil {
		return err
	}

	// Применяем глобальный лимит турниров
	if o.config.MaxTournaments() > 0 && len(baseTournaments) > o.config.MaxTournaments() {
		logger.Info(ctx, "Applying global tournament limit",
			zap.Int("limit", o.config.MaxTournaments()),
			zap.Int("total", len(baseTournaments)))
		baseTournaments = baseTournaments[:o.config.MaxTournaments()]
	}

	logger.Info(ctx, "Found base tournaments for calendar parsing",
		zap.Int("count", len(baseTournaments)))

	// Обрабатываем каждый базовый турнир с AJAX-итерацией по годам/группам
	for idx, t := range baseTournaments {
		logger.Info(ctx, "Processing tournament",
			zap.Int("index", idx+1),
			zap.Int("total", len(baseTournaments)),
			zap.String("name", t.Name))

		// Используем новую логику с AJAX-фильтрацией
		if err := o.processCalendarWithFilters(ctx, t); err != nil {
			logger.Error(ctx, "Failed to process tournament calendar",
				zap.String("tournament", t.Name),
				zap.Error(err))
			continue
		}

		// Standings с AJAX-фильтрацией по годам/группам
		if err := o.processStandingsWithFilters(ctx, t); err != nil {
			logger.Warn(ctx, "Failed to parse standings", zap.Error(err))
		}
	}

	// Парсим детали завершённых матчей
	if err := o.processUnparsedGames(ctx); err != nil {
		logger.Error(ctx, "Failed to process unparsed games", zap.Error(err))
	}

	logger.Info(ctx, "✅ Junior calendar parsing completed")
	return nil
}

func (o *Orchestrator) processUnparsedGames(ctx context.Context) error {
	if !o.config.ParseProtocol() {
		return nil
	}

	matches, err := o.matchRepo.GetUnparsedFinished(ctx, Source, 0)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Processing unparsed finished games",
		zap.Int("count", len(matches)))

	for _, m := range matches {
		if err := o.processGame(ctx, m.ID, m.ExternalID); err != nil {
			logger.Error(ctx, "Failed to process game",
				zap.String("match_id", m.ID),
				zap.Error(err))
			continue
		}
	}

	return nil
}
