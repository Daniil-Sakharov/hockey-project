package dto

// PlayerURLDTO URL игрока для парсинга
type PlayerURLDTO struct {
	URL          string
	PlayerID     string
	TeamID       string
	TournamentID int
}

// PlayerDTO полные данные игрока
type PlayerDTO struct {
	ExternalID  string
	FullName    string
	Position    string
	Number      int
	Role        string
	BirthDate   string
	BirthPlace  string
	Citizenship string
	Height      int
	Weight      int
	Stick       string
	School      string
}
