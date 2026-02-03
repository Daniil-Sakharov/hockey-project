package dto

// ExploreOverviewResponse represents the explore dashboard overview.
type ExploreOverviewResponse struct {
	Players     int64 `json:"players"`
	Teams       int64 `json:"teams"`
	Tournaments int64 `json:"tournaments"`
	Matches     int64 `json:"matches"`
}

// GroupStatsDTO represents stats for a single group within a birth year.
type GroupStatsDTO struct {
	Name         string `json:"name"`
	TeamsCount   int    `json:"teamsCount"`
	MatchesCount int    `json:"matchesCount"`
}

// TournamentItemDTO represents a tournament in list.
type TournamentItemDTO struct {
	ID               string                       `json:"id"`
	Name             string                       `json:"name"`
	Domain           string                       `json:"domain"`
	Season           string                       `json:"season"`
	Source           string                       `json:"source"`
	BirthYearGroups  map[string][]GroupStatsDTO    `json:"birthYearGroups,omitempty"`
	TeamsCount       int                           `json:"teamsCount"`
	MatchesCount     int                           `json:"matchesCount"`
	IsEnded          bool                          `json:"isEnded"`
}

// TournamentListResponse represents the tournaments list.
type TournamentListResponse struct {
	Tournaments []TournamentItemDTO `json:"tournaments"`
}

// TournamentDetailResponse represents tournament detail info.
type TournamentDetailResponse struct {
	Tournament TournamentItemDTO `json:"tournament"`
}

// StandingDTO represents a team standing row.
type StandingDTO struct {
	Position     int    `json:"position"`
	Team         string `json:"team"`
	TeamID       string `json:"teamId"`
	LogoURL      string `json:"logoUrl,omitempty"`
	Games        int    `json:"games"`
	Wins         int    `json:"wins"`
	WinsOT       int    `json:"winsOt"`
	Losses       int    `json:"losses"`
	LossesOT     int    `json:"lossesOt"`
	Draws        int    `json:"draws"`
	GoalsFor     int    `json:"goalsFor"`
	GoalsAgainst int    `json:"goalsAgainst"`
	Points       int    `json:"points"`
	GroupName    string `json:"groupName,omitempty"`
}

// StandingsResponse represents tournament standings.
type StandingsResponse struct {
	Standings []StandingDTO `json:"standings"`
}

// ScorerDTO represents a tournament scorer.
type ScorerDTO struct {
	Position int    `json:"position"`
	PlayerID string `json:"playerId"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl,omitempty"`
	Team     string `json:"team"`
	TeamID   string `json:"teamId"`
	LogoURL  string `json:"logoUrl,omitempty"`
	Games    int    `json:"games"`
	Goals    int    `json:"goals"`
	Assists  int    `json:"assists"`
	Points   int    `json:"points"`
}

// ScorersResponse represents tournament scorers list.
type ScorersResponse struct {
	Scorers []ScorerDTO `json:"scorers"`
}

// SeasonsResponse represents available seasons.
type SeasonsResponse struct {
	Seasons []string `json:"seasons"`
}
