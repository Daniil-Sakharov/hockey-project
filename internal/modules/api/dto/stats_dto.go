package dto

// StatsOverviewResponse represents the stats overview response.
type StatsOverviewResponse struct {
	Players     int64 `json:"players"`
	Teams       int64 `json:"teams"`
	Tournaments int64 `json:"tournaments"`
}
