package dto

// PlayerStatsDTO represents player statistics.
type PlayerStatsDTO struct {
	Games          int `json:"games"`
	Goals          int `json:"goals"`
	Assists        int `json:"assists"`
	Points         int `json:"points"`
	PlusMinus      int `json:"plusMinus"`
	PenaltyMinutes int `json:"penaltyMinutes"`
}

// PlayerItemDTO represents a player in search results.
type PlayerItemDTO struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Position     string          `json:"position"`
	BirthDate    string          `json:"birthDate"`
	BirthYear    int             `json:"birthYear"`
	Team         string          `json:"team"`
	TeamID       string          `json:"teamId"`
	JerseyNumber int             `json:"jerseyNumber"`
	PhotoURL     string          `json:"photoUrl,omitempty"`
	Stats        *PlayerStatsDTO `json:"stats,omitempty"`
}

// PlayersSearchResponse represents the players search result.
type PlayersSearchResponse struct {
	Players []PlayerItemDTO `json:"players"`
	Total   int             `json:"total"`
}

// PlayerProfileResponse represents a full player profile.
type PlayerProfileResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Position     string          `json:"position"`
	BirthDate    string          `json:"birthDate"`
	BirthYear    int             `json:"birthYear"`
	Team         string          `json:"team"`
	TeamID       string          `json:"teamId"`
	JerseyNumber int             `json:"jerseyNumber"`
	Height       *int            `json:"height,omitempty"`
	Weight       *int            `json:"weight,omitempty"`
	Handedness   string          `json:"handedness,omitempty"`
	City         string          `json:"city,omitempty"`
	PhotoURL     string          `json:"photoUrl,omitempty"`
	Stats        *PlayerStatsDTO `json:"stats,omitempty"`
}

// PlayerStatDTO represents a detailed stat entry for a player.
type PlayerStatDTO struct {
	Season         string `json:"season"`
	TournamentID   string `json:"tournamentId"`
	TournamentName string `json:"tournamentName"`
	GroupName      string `json:"groupName"`
	BirthYear      int    `json:"birthYear"`
	Games          int    `json:"games"`
	Goals          int    `json:"goals"`
	Assists        int    `json:"assists"`
	Points         int    `json:"points"`
	PlusMinus      int    `json:"plusMinus"`
	PenaltyMinutes int    `json:"penaltyMinutes"`
}

// PlayerStatsHistoryResponse represents detailed player stats.
type PlayerStatsHistoryResponse struct {
	Stats []PlayerStatDTO `json:"stats"`
}

// TeamStatsDTO represents aggregated team stats.
type TeamStatsDTO struct {
	Wins         int `json:"wins"`
	Losses       int `json:"losses"`
	Draws        int `json:"draws"`
	GoalsFor     int `json:"goalsFor"`
	GoalsAgainst int `json:"goalsAgainst"`
}

// TeamProfileResponse represents a team profile.
type TeamProfileResponse struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	City          string          `json:"city"`
	LogoURL       string          `json:"logoUrl,omitempty"`
	Tournaments   []string        `json:"tournaments"`
	PlayersCount  int             `json:"playersCount"`
	Roster        []PlayerItemDTO `json:"roster"`
	Stats         TeamStatsDTO    `json:"stats"`
	RecentMatches []MatchDTO      `json:"recentMatches"`
}
