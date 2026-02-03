package dto

// SeasonDTO представляет сезон из API
type SeasonDTO struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`    // "25/26"
	Current bool   `json:"current"` // true для текущего сезона
}
