package dto

// Profile представляет полный профиль игрока для отображения (Presentation DTO)
type Profile struct {
	BasicInfo          PlayerBasicInfo
	AllTimeStats       *PlayerStats
	CurrentSeasonStats *SeasonStats
	RecentTournaments  []TournamentStats
}

// PlayerBasicInfo базовая информация об игроке
type PlayerBasicInfo struct {
	ID        string
	Name      string
	BirthYear int
	Position  string
	Height    *int
	Weight    *int
	Team      string
	Region    string
}

// PlayerStats общая статистика игрока
type PlayerStats struct {
	Tournaments      int
	Games            int
	Goals            int
	Assists          int
	Points           int
	PlusMinus        int
	Penalties        int
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64
	HatTricks        int
	GameWinningGoals int
}

// SeasonStats статистика за конкретный сезон
type SeasonStats struct {
	Season           string
	Tournaments      int
	Games            int
	Goals            int
	Assists          int
	Points           int
	PlusMinus        int
	Penalties        int
	GoalsPerGame     float64
	AssistsPerGame   float64
	PointsPerGame    float64
	PenaltiesPerGame float64
	HatTricks        int
	GameWinningGoals int
}

// TournamentStats статистика по одному турниру
type TournamentStats struct {
	TournamentName string
	GroupName      string
	Season         string
	IsChampionship bool
	Games          int
	Goals          int
	Assists        int
	Points         int
	PlusMinus      int
	Penalties      int
	HatTricks      int
	WinningGoals   int
}
