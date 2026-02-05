package dto

import "time"

// MatchDTO представляет матч из календаря турнира
type MatchDTO struct {
	ExternalID  string // ID матча (54503)
	MatchNumber int    // Номер в календаре
	Round       int    // Тур

	HomeTeamID   string // ID команды A (домашняя)
	HomeTeamName string
	AwayTeamID   string // ID команды B (гостевая)
	AwayTeamName string

	HomeScore int // Счет домашней команды
	AwayScore int // Счет гостевой команды

	ScheduledAt time.Time // Дата + время матча
	Venue       string    // Стадион + город
	VenueCity   string    // Город из стадиона

	ProtoURL string // Ссылка на протокол матча
}

// MatchProtocolDTO представляет детали протокола матча
type MatchProtocolDTO struct {
	HomeLogoURL string // URL логотипа домашней команды
	AwayLogoURL string // URL логотипа гостевой команды

	// Счет по периодам [period][0=home, 1=away], period: 0-2 периоды, 3 - OT
	ScoreByPeriod [4][2]int

	Goals      []GoalEventDTO
	Penalties  []PenaltyEventDTO
	HomeLineup []LineupPlayerDTO
	AwayLineup []LineupPlayerDTO
}

// GoalEventDTO представляет гол в матче
type GoalEventDTO struct {
	Period      int // Период (1, 2, 3, 4=OT)
	TimeMinutes int
	TimeSeconds int

	ScorerID    string // ID забившего игрока
	ScorerName  string
	Assist1ID   string // ID первого ассистента
	Assist1Name string
	Assist2ID   string // ID второго ассистента
	Assist2Name string

	GoalType   string // even/pp/sh
	ScoreAfter [2]int // Счет после гола [home, away]
	IsHome     bool   // Гол домашней команды
}

// PenaltyEventDTO представляет удаление в матче
type PenaltyEventDTO struct {
	Period      int
	TimeMinutes int
	TimeSeconds int

	PlayerID   string
	PlayerName string
	Minutes    int    // Минуты штрафа (2, 5, 10)
	Reason     string // Причина удаления
	IsHome     bool   // Удаление игрока домашней команды
}

// LineupPlayerDTO представляет игрока в составе на матч
type LineupPlayerDTO struct {
	PlayerID    string
	PlayerName  string
	Number      int
	Position    string // Г/З/Н
	CaptainRole string // К/А/пусто
}
