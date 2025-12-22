package logger

// StatsLogger - интерфейс для логирования процесса парсинга статистики
type StatsLogger interface {
	// LogTournamentStart логирует начало парсинга турнира
	LogTournamentStart(tournamentID, tournamentName, url string)

	// LogCombinationStart логирует начало обработки комбинации год+группа
	LogCombinationStart(year, yearID, group, groupID string)

	// LogCombinationResult логирует результат обработки комбинации
	LogCombinationResult(receivedCount int, playerIDs []string)

	// LogCombinationError логирует ошибку при обработке комбинации
	LogCombinationError(err error)

	// LogTournamentSummary логирует финальную сводку по турниру
	LogTournamentSummary(receivedCount, savedCount int)

	// LogValidationSkip логирует пропуск записи при валидации
	LogValidationSkip(playerID, teamID, reason string)

	// LogFKConstraintSkip логирует пропуск записи из-за FK constraint
	LogFKConstraintSkip(playerID, teamID, reason string)
}
