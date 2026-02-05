package game

// GameDetailsDTO детальная информация о матче
type GameDetailsDTO struct {
	ExternalID   string // ID матча
	HomeScore    *int   // Счёт дома
	AwayScore    *int   // Счёт гостей
	HomeScoreP1  *int   // Счёт 1-го периода
	AwayScoreP1  *int
	HomeScoreP2  *int
	AwayScoreP2  *int
	HomeScoreP3  *int
	AwayScoreP3  *int
	HomeScoreOT  *int
	AwayScoreOT  *int
	ResultType   string           // regular, OT, SO
	VideoURL     string           // URL видео
	HomeTeamURL  string           // URL домашней команды
	AwayTeamURL  string           // URL гостевой команды
	BirthYear    int              // Год рождения (2008, 2009...)
	GroupName    string           // Группа (А1, Б2...)
	Goals        []GoalDTO        // Голы
	Penalties    []PenaltyDTO     // Штрафы
	GoalieEvents []GoalieEventDTO // События вратарей (смена, выход)
	EmptyNets    []EmptyNetDTO    // Пустые ворота
	Timeouts     []TimeoutDTO     // Тайм-ауты
	HomeLineup   []PlayerLineup   // Состав дома
	AwayLineup   []PlayerLineup   // Состав гостей
}

// GoalDTO информация о голе
type GoalDTO struct {
	Period        int    // Период (1, 2, 3, OT)
	TimeMinutes   int    // Минуты
	TimeSeconds   int    // Секунды
	ScorerURL     string // URL забившего
	ScorerName    string // Имя забившего
	Assist1URL    string // URL первого ассистента
	Assist1Name   string
	Assist1Number int    // Номер первого ассистента (для поиска в составе)
	Assist2URL    string // URL второго ассистента
	Assist2Name   string
	Assist2Number int    // Номер второго ассистента (для поиска в составе)
	TeamURL       string // URL команды
	GoalType      string // even, pp, sh, en

	// Кто был на льду
	GoalieURL        string   // URL вратаря (пропустившего)
	HomePlayersOnIce []string // URL игроков дома на льду
	AwayPlayersOnIce []string // URL игроков гостей на льду
}

// PenaltyDTO информация о штрафе
type PenaltyDTO struct {
	Period      int    // Период
	TimeMinutes int    // Минуты
	TimeSeconds int    // Секунды
	PlayerURL   string // URL игрока
	PlayerName  string // Имя игрока
	TeamURL     string // URL команды
	Minutes     int    // Минуты штрафа (2, 5, 10...)
	Reason      string // Причина (грубость, задержка...)
	IsHome      bool   // Домашняя команда
}

// PlayerLineup информация об игроке в составе
type PlayerLineup struct {
	PlayerURL      string // URL профиля
	PlayerName     string // Имя игрока
	JerseyNumber   int    // Номер
	Position       string // G, D, F
	Role           string // C (капитан), A (ассистент), пусто
	Goals          int    // Голы в матче
	Assists        int    // Передачи
	PenaltyMinutes int    // Штрафные минуты
	PlusMinus      int    // +/-
	Saves          *int   // Спасения (для вратарей)
	GoalsAgainst   *int   // Пропущено (для вратарей)
	TimeOnIce      *int   // Время на льду (секунды)
}

// GoalieEventDTO информация о событии вратаря (смена, выход на лёд)
type GoalieEventDTO struct {
	TimeMinutes int    // Минуты
	TimeSeconds int    // Секунды
	PlayerURL   string // URL вратаря
	PlayerName  string // Имя вратаря
	TeamURL     string // URL команды
	IsHome      bool   // Домашняя команда
}

// EmptyNetDTO информация о пустых воротах
type EmptyNetDTO struct {
	TimeMinutes int    // Минуты начала
	TimeSeconds int    // Секунды
	TeamURL     string // URL команды со снятым вратарём
	IsHome      bool   // Домашняя команда
}

// TimeoutDTO информация о тайм-ауте
type TimeoutDTO struct {
	TimeMinutes int    // Минуты
	TimeSeconds int    // Секунды
	TeamURL     string // URL команды
	IsHome      bool   // Домашняя команда
}
