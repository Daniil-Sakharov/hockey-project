package player

// Profile представляет полный профиль игрока для отображения в Telegram
type Profile struct {
	// Базовая информация игрока
	BasicInfo PlayerBasicInfo

	// Общая статистика за всё время (может быть nil если нет данных)
	AllTimeStats *PlayerStats

	// Статистика за текущий/последний сезон (может быть nil если нет данных)
	CurrentSeasonStats *SeasonStats

	// Последние турниры игрока (обычно 5)
	RecentTournaments []TournamentStats
}

// PlayerBasicInfo базовая информация об игроке
type PlayerBasicInfo struct {
	ID       string
	Name     string
	BirthYear int    // Год рождения (из birth_date)
	Position string  // Защитник/Нападающий/Вратарь
	Height   *int
	Weight   *int
	Team     string  // Название команды
	Region   string  // Регион/Город команды
}

// PlayerStats общая статистика игрока
type PlayerStats struct {
	Tournaments int // Количество турниров
	Games       int
	Goals       int
	Assists     int
	Points      int
	PlusMinus   int
	Penalties   int // penalty_minutes

	// Средние показатели
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64

	// Достижения
	HatTricks        int
	GameWinningGoals int
}

// SeasonStats статистика за конкретный сезон
type SeasonStats struct {
	Season      string // "2024-2025"
	Tournaments int
	Games       int
	Goals       int
	Assists     int
	Points      int
	PlusMinus   int
	Penalties   int

	// Средние показатели
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64

	// Достижения в сезоне
	HatTricks        int
	GameWinningGoals int
}

// TournamentStats статистика по одному турниру
type TournamentStats struct {
	TournamentName string
	GroupName      string  // Может быть пустой если группа "Общая" и она единственная
	Season         string  // "2024-2025"
	IsChampionship bool    // true = Первенство, false = Кубок/другое
	Games          int
	Goals          int
	Assists        int
	Points         int
	PlusMinus      int
	Penalties      int
	HatTricks      int     // Хет-трики в этом турнире
	WinningGoals   int     // Победные голы в этом турнире
}
