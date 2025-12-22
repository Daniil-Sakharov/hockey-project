package services

// FullPlayerReport полные данные для HTML отчета игрока
type FullPlayerReport struct {
	Player        ReportPlayerInfo
	TotalStats    ReportTotalStats
	GoalsByType   GoalsBreakdown
	GoalsByPeriod PeriodGoals
	SeasonStats   []SeasonSummary
	Tournaments   []TournamentStats

	HasStats           bool
	HasDetailedStats   bool
	HasMultipleSeasons bool
}

// ReportPlayerInfo базовая информация об игроке
type ReportPlayerInfo struct {
	ID        string
	Name      string
	BirthYear int
	Position  string
	Height    *int
	Weight    *int
	Team      string
	Region    string
}

// ReportTotalStats сводная статистика за всё время
type ReportTotalStats struct {
	TotalTournaments int
	TotalGames       int
	TotalGoals       int
	TotalAssists     int
	TotalPoints      int
	TotalPlusMinus   int
	TotalPenalties   int

	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64

	GoalsEvenStrength int
	GoalsPowerPlay    int
	GoalsShortHanded  int
	GoalsPeriod1      int
	GoalsPeriod2      int
	GoalsPeriod3      int
	GoalsOvertime     int

	TotalHatTricks    int
	TotalWinningGoals int
}

// GoalsBreakdown распределение голов по типу
type GoalsBreakdown struct {
	EvenStrength int
	PowerPlay    int
	ShortHanded  int
}

// PeriodGoals голы по периодам
type PeriodGoals struct {
	Period1  int
	Period2  int
	Period3  int
	Overtime int
}

// SeasonSummary статистика за сезон
type SeasonSummary struct {
	Season  string
	Games   int
	Goals   int
	Assists int
	Points  int
}

// TournamentStats статистика по турниру
type TournamentStats struct {
	Season         string
	TournamentName string
	TournamentID   string
	GroupName      string
	TeamName       string

	Games            int
	Goals            int
	Assists          int
	Points           int
	PlusMinus        int
	PenaltyMinutes   int
	HatTricks        int
	GameWinningGoals int
}
