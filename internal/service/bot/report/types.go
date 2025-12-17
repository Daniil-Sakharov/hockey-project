package report

// FullPlayerReport полные данные для HTML отчета игрока
type FullPlayerReport struct {
	// Базовая информация
	Player PlayerInfo

	// Сводная статистика за всё время
	TotalStats TotalStatistics

	// Данные для графиков
	GoalsByType   GoalsBreakdown  // Распределение голов (ESG, PPG, SHG)
	GoalsByPeriod PeriodGoals     // Голы по периодам
	SeasonStats   []SeasonSummary // Статистика по сезонам (для линейного графика)

	// Полная история турниров (ВСЕ турниры)
	Tournaments []TournamentFullStats

	// Флаги для условного отображения секций
	HasStats           bool
	HasDetailedStats   bool
	HasMultipleSeasons bool
}

// PlayerInfo базовая информация об игроке
type PlayerInfo struct {
	ID        string
	Name      string
	BirthYear int
	Position  string
	Height    *int
	Weight    *int
	Team      string
	Region    string
}

// TotalStatistics сводная статистика за всё время
type TotalStatistics struct {
	TotalTournaments int
	TotalGames       int
	TotalGoals       int
	TotalAssists     int
	TotalPoints      int
	TotalPlusMinus   int
	TotalPenalties   int

	// Средние показатели
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64

	// Детальная статистика голов
	GoalsEvenStrength int
	GoalsPowerPlay    int
	GoalsShortHanded  int
	GoalsPeriod1      int
	GoalsPeriod2      int
	GoalsPeriod3      int
	GoalsOvertime     int

	// Достижения
	TotalHatTricks    int
	TotalWinningGoals int
}

// GoalsBreakdown распределение голов по типу (для круговой диаграммы)
type GoalsBreakdown struct {
	EvenStrength    int     // В равных составах
	PowerPlay       int     // В большинстве
	ShortHanded     int     // В меньшинстве
	EvenStrengthPct float64 // Процент
	PowerPlayPct    float64
	ShortHandedPct  float64
}

// PeriodGoals голы по периодам (для столбчатой диаграммы)
type PeriodGoals struct {
	Period1  int
	Period2  int
	Period3  int
	Overtime int
}

// SeasonSummary статистика за сезон (для линейного графика прогресса)
type SeasonSummary struct {
	Season  string
	Games   int
	Goals   int
	Assists int
	Points  int
}

// TournamentFullStats полная статистика по одному турниру
type TournamentFullStats struct {
	// Информация о турнире
	Season         string
	TournamentName string
	TournamentID   string
	GroupName      string
	TeamName       string

	// Основная статистика
	Games          int
	Goals          int
	Assists        int
	Points         int
	Plus           int
	Minus          int
	PlusMinus      int
	PenaltyMinutes int

	// Детальная статистика голов
	GoalsEvenStrength int
	GoalsPowerPlay    int
	GoalsShortHanded  int
	GoalsPeriod1      int
	GoalsPeriod2      int
	GoalsPeriod3      int
	GoalsOvertime     int

	// Достижения
	HatTricks        int
	GameWinningGoals int

	// Средние показатели
	GoalsPerGame     float64
	PointsPerGame    float64
	PenaltiesPerGame float64
}

// ChartData данные для JavaScript графиков
type ChartData struct {
	// Для круговой диаграммы голов по типу
	GoalTypeLabels []string `json:"goalTypeLabels"`
	GoalTypeValues []int    `json:"goalTypeValues"`
	GoalTypeColors []string `json:"goalTypeColors"`

	// Для столбчатой диаграммы голов по периодам
	PeriodLabels []string `json:"periodLabels"`
	PeriodValues []int    `json:"periodValues"`

	// Для линейного графика прогресса по сезонам
	SeasonLabels []string `json:"seasonLabels"`
	SeasonGoals  []int    `json:"seasonGoals"`
	SeasonPoints []int    `json:"seasonPoints"`

	// Для radar chart профиля игрока
	RadarLabels []string  `json:"radarLabels"`
	RadarValues []float64 `json:"radarValues"`
}
