package calendar

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	"go.uber.org/zap"
)

// Run запускает полный цикл парсинга календаря FHSPB
func (o *Orchestrator) Run(ctx context.Context) error {
	logger.Info(ctx, "Starting FHSPB calendar parsing...")

	// Получаем все турниры FHSPB
	tournaments, err := o.tournamentRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("get tournaments: %w", err)
	}

	logger.Info(ctx, "Found tournaments for parsing",
		zap.Int("count", len(tournaments)))

	// Обрабатываем каждый турнир
	for _, t := range tournaments {
		if err := o.processTournament(ctx, t.ID, t.ExternalID); err != nil {
			logger.Error(ctx, "Failed to process tournament",
				zap.String("tournament", t.Name),
				zap.String("id", t.ID),
				zap.Error(err))
			continue
		}
	}

	// Парсим детали завершённых матчей
	if o.config.ParseProtocol() {
		if err := o.processUnparsedGames(ctx); err != nil {
			logger.Error(ctx, "Failed to process unparsed games", zap.Error(err))
		}
	}

	logger.Info(ctx, "FHSPB calendar parsing completed")
	return nil
}

func (o *Orchestrator) processTournament(ctx context.Context, tournamentID, externalID string) error {
	logger.Debug(ctx, "Processing tournament",
		zap.String("id", tournamentID),
		zap.String("external_id", externalID))

	// Парсим турнирную таблицу
	if err := o.processStandings(ctx, tournamentID, externalID); err != nil {
		logger.Warn(ctx, "Failed to parse standings",
			zap.String("tournament_id", tournamentID),
			zap.Error(err))
	}

	// Парсим календарь
	if err := o.processCalendar(ctx, tournamentID, externalID); err != nil {
		logger.Warn(ctx, "Failed to parse calendar",
			zap.String("tournament_id", tournamentID),
			zap.Error(err))
	}

	return nil
}

func (o *Orchestrator) processUnparsedGames(ctx context.Context) error {
	matches, err := o.matchRepo.GetUnparsedFinished(ctx, Source, 100)
	if err != nil {
		return fmt.Errorf("get unparsed matches: %w", err)
	}

	logger.Info(ctx, "Processing unparsed finished games",
		zap.Int("count", len(matches)))

	for _, m := range matches {
		// Извлекаем external ID турнира из внутреннего ID
		// Формат: "spb:6376" -> "6376"
		tournamentExternalID := ""
		if m.TournamentID != nil {
			tournamentExternalID = extractExternalID(*m.TournamentID)
		}

		if err := o.processGame(ctx, m.ID, m.ExternalID, tournamentExternalID); err != nil {
			logger.Error(ctx, "Failed to process game",
				zap.String("match_id", m.ID),
				zap.Error(err))
			continue
		}
	}

	return nil
}

// extractExternalID извлекает внешний ID из внутреннего формата
// "spb:6376" -> "6376"
func extractExternalID(internalID string) string {
	if idx := strings.Index(internalID, ":"); idx >= 0 {
		return internalID[idx+1:]
	}
	return internalID
}
