package logger

import (
	"go.uber.org/zap"
)

// statsLogger реализует StatsLogger используя pkg/logger (zap)
type statsLogger struct {
	logger *zap.Logger
}

// New создает новый логгер статистики
func New(logger *zap.Logger) StatsLogger {
	return &statsLogger{
		logger: logger.With(zap.String("component", "stats_parser")),
	}
}

// LogTournamentStart логирует начало парсинга турнира
func (l *statsLogger) LogTournamentStart(tournamentID, tournamentName, url string) {
	l.logger.Info("Начало парсинга турнира",
		zap.String("tournament_id", tournamentID),
		zap.String("tournament_name", tournamentName),
		zap.String("url", url),
	)
}

// LogCombinationStart логирует начало обработки комбинации год+группа
func (l *statsLogger) LogCombinationStart(year, yearID, group, groupID string) {
	l.logger.Info("Обработка комбинации",
		zap.String("year", year),
		zap.String("year_id", yearID),
		zap.String("group", group),
		zap.String("group_id", groupID),
	)
}

// LogCombinationResult логирует результат обработки комбинации
func (l *statsLogger) LogCombinationResult(receivedCount int, playerIDs []string) {
	l.logger.Info("Получены данные из API",
		zap.Int("received_count", receivedCount),
		zap.Int("player_ids_count", len(playerIDs)),
		zap.Strings("player_ids", playerIDs),
	)
}

// LogCombinationError логирует ошибку при обработке комбинации
func (l *statsLogger) LogCombinationError(err error) {
	l.logger.Error("Ошибка при обработке комбинации",
		zap.Error(err),
	)
}

// LogTournamentSummary логирует финальную сводку по турниру
func (l *statsLogger) LogTournamentSummary(receivedCount, savedCount int) {
	lostCount := receivedCount - savedCount
	lostPercent := 0.0
	if receivedCount > 0 {
		lostPercent = float64(lostCount) * 100.0 / float64(receivedCount)
	}

	l.logger.Info("Итоги парсинга турнира",
		zap.Int("received_total", receivedCount),
		zap.Int("saved_total", savedCount),
		zap.Int("lost_total", lostCount),
		zap.Float64("lost_percent", lostPercent),
	)
}

// LogValidationSkip логирует пропуск записи при валидации
func (l *statsLogger) LogValidationSkip(playerID, teamID, reason string) {
	l.logger.Warn("Запись пропущена при валидации",
		zap.String("player_id", playerID),
		zap.String("team_id", teamID),
		zap.String("reason", reason),
	)
}

// LogFKConstraintSkip логирует пропуск записи из-за FK constraint
func (l *statsLogger) LogFKConstraintSkip(playerID, teamID, reason string) {
	l.logger.Warn("FK constraint: запись пропущена",
		zap.String("player_id", playerID),
		zap.String("team_id", teamID),
		zap.String("reason", reason),
	)
}
