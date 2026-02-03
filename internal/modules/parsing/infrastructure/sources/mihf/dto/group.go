package dto

// GroupDTO представляет группу турнира (по году рождения)
type GroupDTO struct {
	ID        string // 76, 77, 78
	Name      string // Первенство г. Москвы
	BirthYear int    // 2008, 2009, ...
	GroupName string // Группа А, Группа Б
	URL       string // /championat/2023/groups/76
}

// TournamentDTO представляет турнир (по году рождения внутри группы)
type TournamentDTO struct {
	ID        string // 330
	GroupID   string // 76
	Name      string // 2007 г.р.
	BirthYear int    // 2007
	URL       string // /championat/2023/groups/76/tournament/330
}

// SubTournamentDTO представляет подгруппу турнира (Группа А, Б, В)
type SubTournamentDTO struct {
	ID           string // 934
	TournamentID string // 330
	GroupID      string // 76
	Name         string // Группа А
	URL          string // /championat/2023/groups/76/tournament/330/sub/934
}

// TournamentPathDTO путь к конкретному турниру
type TournamentPathDTO struct {
	SeasonYear   string // 2023
	GroupID      string // 76
	TournamentID string // 330
	SubID        string // 934
	BirthYear    int    // 2008
	GroupName    string // Группа А
}

// ScoreboardURL возвращает URL турнирной таблицы
func (t TournamentPathDTO) ScoreboardURL() string {
	return "/championat/" + t.SeasonYear + "/groups/" + t.GroupID +
		"/tournament/" + t.TournamentID + "/sub/" + t.SubID + "/scoreboard"
}

// TeamURL возвращает URL команды
func (t TournamentPathDTO) TeamURL(teamID string) string {
	return "/championat/" + t.SeasonYear + "/" + t.GroupID +
		"/tournament/" + t.TournamentID + "/sub/" + t.SubID + "/team/" + teamID
}
