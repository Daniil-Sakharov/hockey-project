package dto

// MatchDTO represents a match item.
type MatchDTO struct {
	ID          string `json:"id"`
	HomeTeam    string `json:"homeTeam"`
	AwayTeam    string `json:"awayTeam"`
	HomeTeamID  string `json:"homeTeamId"`
	AwayTeamID  string `json:"awayTeamId"`
	HomeLogoURL string `json:"homeLogoUrl,omitempty"`
	AwayLogoURL string `json:"awayLogoUrl,omitempty"`
	HomeScore   *int   `json:"homeScore"`
	AwayScore   *int   `json:"awayScore"`
	ResultType  string `json:"resultType,omitempty"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Tournament  string `json:"tournament"`
	Venue       string `json:"venue,omitempty"`
	Status      string `json:"status"`
}

// MatchListResponse represents a list of matches.
type MatchListResponse struct {
	Matches []MatchDTO `json:"matches"`
}

// RankedPlayerDTO represents a player in rankings.
type RankedPlayerDTO struct {
	Rank           int    `json:"rank"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	PhotoURL       string `json:"photoUrl,omitempty"`
	Position       string `json:"position"`
	BirthYear      int    `json:"birthYear"`
	Team           string `json:"team"`
	TeamID         string `json:"teamId"`
	TeamLogoURL    string `json:"teamLogoUrl,omitempty"`
	TeamCity       string `json:"teamCity,omitempty"`
	Games          int    `json:"games"`
	Goals          int    `json:"goals"`
	Assists        int    `json:"assists"`
	Points         int    `json:"points"`
	PlusMinus      int    `json:"plusMinus"`
	PenaltyMinutes int    `json:"penaltyMinutes"`
}

// RankingsResponse represents player rankings.
type RankingsResponse struct {
	Season  string            `json:"season"`
	Players []RankedPlayerDTO `json:"players"`
}

// DomainOption represents a region/domain filter option.
type DomainOption struct {
	Domain string `json:"domain"`
	Label  string `json:"label"`
}

// TournamentOption represents a tournament filter option.
type TournamentOption struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	BirthYears []int  `json:"birthYears,omitempty"`
}

// GroupOption represents a group filter option within a tournament.
type GroupOption struct {
	Name         string `json:"name"`
	TournamentID string `json:"tournamentId"`
}

// RankingsFiltersResponse contains available filter values.
type RankingsFiltersResponse struct {
	BirthYears  []int              `json:"birthYears"`
	Domains     []DomainOption     `json:"domains"`
	Tournaments []TournamentOption `json:"tournaments"`
	Groups      []GroupOption      `json:"groups"`
}

// MatchDetailDTO represents detailed match information.
type MatchDetailDTO struct {
	ID            string            `json:"id"`
	ExternalID    string            `json:"externalId"`
	HomeTeam      MatchTeamDTO      `json:"homeTeam"`
	AwayTeam      MatchTeamDTO      `json:"awayTeam"`
	HomeScore     *int              `json:"homeScore"`
	AwayScore     *int              `json:"awayScore"`
	ScoreByPeriod *ScoreByPeriodDTO `json:"scoreByPeriod,omitempty"`
	ResultType    string            `json:"resultType,omitempty"`
	Date          string            `json:"date"`
	Time          string            `json:"time"`
	Tournament    TournamentInfoDTO `json:"tournament"`
	Venue         string            `json:"venue,omitempty"`
	Status        string            `json:"status"`
	GroupName     string            `json:"groupName,omitempty"`
	BirthYear     *int              `json:"birthYear,omitempty"`
	MatchNumber   *int              `json:"matchNumber,omitempty"`
	Events        []MatchEventDTO   `json:"events"`
	HomeLineup    []LineupPlayerDTO `json:"homeLineup"`
	AwayLineup    []LineupPlayerDTO `json:"awayLineup"`
}

// MatchTeamDTO represents a team in match detail.
type MatchTeamDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city,omitempty"`
	LogoURL string `json:"logoUrl,omitempty"`
}

// ScoreByPeriodDTO represents score breakdown by period.
type ScoreByPeriodDTO struct {
	HomeP1 *int `json:"homeP1,omitempty"`
	AwayP1 *int `json:"awayP1,omitempty"`
	HomeP2 *int `json:"homeP2,omitempty"`
	AwayP2 *int `json:"awayP2,omitempty"`
	HomeP3 *int `json:"homeP3,omitempty"`
	AwayP3 *int `json:"awayP3,omitempty"`
	HomeOT *int `json:"homeOt,omitempty"`
	AwayOT *int `json:"awayOt,omitempty"`
}

// TournamentInfoDTO represents tournament info in match detail.
type TournamentInfoDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MatchEventDTO represents a match event (goal, penalty).
type MatchEventDTO struct {
	Type        string `json:"type"`
	Period      *int   `json:"period,omitempty"`
	Time        string `json:"time,omitempty"`
	IsHome      bool   `json:"isHome"`
	TeamName    string `json:"teamName,omitempty"`
	TeamLogoURL string `json:"teamLogoUrl,omitempty"`
	PlayerID    string `json:"playerId,omitempty"`
	PlayerName  string `json:"playerName,omitempty"`
	PlayerPhoto string `json:"playerPhoto,omitempty"`
	Assist1ID   string `json:"assist1Id,omitempty"`
	Assist1Name string `json:"assist1Name,omitempty"`
	Assist2ID   string `json:"assist2Id,omitempty"`
	Assist2Name string `json:"assist2Name,omitempty"`
	GoalType    string `json:"goalType,omitempty"`
	PenaltyMins *int   `json:"penaltyMins,omitempty"`
	PenaltyText string `json:"penaltyText,omitempty"`
}

// LineupPlayerDTO represents a player in match lineup.
type LineupPlayerDTO struct {
	PlayerID       string `json:"playerId"`
	PlayerName     string `json:"playerName"`
	PlayerPhoto    string `json:"playerPhoto,omitempty"`
	JerseyNumber   *int   `json:"jerseyNumber,omitempty"`
	Position       string `json:"position,omitempty"`
	Goals          int    `json:"goals"`
	Assists        int    `json:"assists"`
	Points         int    `json:"points"`
	PenaltyMinutes int    `json:"penaltyMinutes"`
	PlusMinus      int    `json:"plusMinus"`
	Saves          *int   `json:"saves,omitempty"`
	GoalsAgainst   *int   `json:"goalsAgainst,omitempty"`
}
