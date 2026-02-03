package dto

// TopScorerResponse represents a single top scorer.
type TopScorerResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Team    string `json:"team"`
	Goals   int    `json:"goals"`
	Assists int    `json:"assists"`
	Games   int    `json:"games"`
}

// TopScorersResponse represents the top scorers list response.
type TopScorersResponse struct {
	Players []TopScorerResponse `json:"players"`
}
