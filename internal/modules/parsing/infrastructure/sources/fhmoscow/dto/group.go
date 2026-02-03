package dto

// GroupDTO представляет группу турнира
type GroupDTO struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`    // "Группа А"
	Current bool   `json:"current"` // true если текущая
}

// StageDTO представляет этап турнира
type StageDTO struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`    // "Первый этап"
	Current bool   `json:"current"` // true если текущий
}

// TourDTO представляет тур турнира
type TourDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"` // "1 тур"
}

// FilterDataResponse ответ API /api/filter/data
type FilterDataResponse struct {
	Tournament []TournamentDTO `json:"tournament"`
	Stage      []StageDTO      `json:"stage"`
	Group      []GroupDTO      `json:"group"`
	Tour       []TourDTO       `json:"tour"`
}
