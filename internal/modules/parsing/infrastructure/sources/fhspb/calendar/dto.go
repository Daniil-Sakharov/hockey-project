package calendar

import "time"

// MatchDTO представляет матч из календаря FHSPB
type MatchDTO struct {
	ExternalID   string     // MatchID из URL
	MatchNumber  int        // Номер матча в турнире
	ScheduledAt  *time.Time // Дата и время
	Venue        string     // Стадион (код арены)
	HomeTeamName string     // Название домашней команды
	AwayTeamName string     // Название гостевой команды
	HomeScore    *int       // Счёт дома (nil если не сыгран)
	AwayScore    *int       // Счёт гостей
	ResultType   string     // "", "OT", "SO" (ОТ, ПБ)
	IsFinished   bool
}

// MatchDetailsURL возвращает URL для получения деталей матча
func (m *MatchDTO) MatchDetailsURL(tournamentID string) string {
	return "/Match?TournamentID=" + tournamentID + "&MatchID=" + m.ExternalID
}
