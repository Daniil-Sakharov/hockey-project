package match

// MatchDetailsDTO протокол матча
type MatchDetailsDTO struct {
	ExternalID   string
	HomeTeamName string
	AwayTeamName string

	// Счёт по периодам
	HomeScoreP1 int
	AwayScoreP1 int
	HomeScoreP2 int
	AwayScoreP2 int
	HomeScoreP3 int
	AwayScoreP3 int
	HomeScoreOT int
	AwayScoreOT int

	// События
	Goals     []GoalDTO
	Penalties []PenaltyDTO

	// Составы
	HomeLineup []PlayerLineupDTO
	AwayLineup []PlayerLineupDTO

	// Вратари
	HomeGoalies []GoalieStatsDTO
	AwayGoalies []GoalieStatsDTO

	// Броски
	HomeShots ShotsDTO
	AwayShots ShotsDTO
}

// GoalDTO гол
type GoalDTO struct {
	Period       int
	TimeMinutes  int
	TimeSeconds  int
	ScorerURL    string // URL игрока
	ScorerName   string
	ScorerNumber int
	Assist1URL   string
	Assist1Name  string
	Assist2URL   string
	Assist2Name  string
	TeamName     string // Полное имя команды
	TeamAbbr     string // ГПБ, ОН и т.д.
	ScoreHome    int    // Счёт после гола
	ScoreAway    int
	GoalType     string // PP1, PP2, SH1, SH2, EN, PS, GWG или пусто
	IsHome       bool   // Гол домашней команды
}

// PenaltyDTO штраф
type PenaltyDTO struct {
	Period       int
	TimeMinutes  int
	TimeSeconds  int
	PlayerURL    string
	PlayerName   string
	PlayerNumber int
	TeamAbbr     string
	Minutes      int
	Reason       string // Полное описание нарушения
	ReasonCode   string // Код нарушения: ПОДН, ГРУБ, НП-АТ
	IsHome       bool   // Штраф домашней команды
}

// PlayerLineupDTO игрок в составе
type PlayerLineupDTO struct {
	PlayerURL      string
	PlayerName     string
	Number         int
	Position       string // Нп, Зщ, Вр → forward, defenseman, goalkeeper
	Played         bool   // Играл ли
	Points         int    // Очки (голы + передачи)
	Goals          int
	Assists        int
	PenaltyMinutes int
	PlusMinus      int
	CaptainRole    string // C = Captain, A = Assistant, или ""
}

// GoalieStatsDTO статистика вратаря
type GoalieStatsDTO struct {
	PlayerURL      string
	PlayerName     string
	Number         int
	Played         bool // Играл ли
	TimeOnIce      int  // в секундах
	GoalsAgainst   int
	ShotsAgainst   int
	SavePercentage float64
	PenaltyMinutes int
}

// ShotsDTO броски по периодам
type ShotsDTO struct {
	P1    int
	P2    int
	P3    int
	OT    int
	Total int
}
