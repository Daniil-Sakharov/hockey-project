package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// MatchRow represents a match from DB.
type MatchRow struct {
	ID          string     `db:"id"`
	HomeTeam    string     `db:"home_team"`
	AwayTeam    string     `db:"away_team"`
	HomeTeamID  string     `db:"home_team_id"`
	AwayTeamID  string     `db:"away_team_id"`
	HomeLogoURL string     `db:"home_logo_url"`
	AwayLogoURL string     `db:"away_logo_url"`
	HomeScore   *int       `db:"home_score"`
	AwayScore   *int       `db:"away_score"`
	ResultType  string     `db:"result_type"`
	ScheduledAt *time.Time `db:"scheduled_at"`
	Tournament  string     `db:"tournament_name"`
	Venue       string     `db:"venue"`
	Status      string     `db:"status"`
}

// RankedPlayerRow represents a ranked player from DB.
type RankedPlayerRow struct {
	ID             string    `db:"id"`
	Name           string    `db:"name"`
	PhotoURL       string    `db:"photo_url"`
	Position       string    `db:"position"`
	BirthDate      time.Time `db:"birth_date"`
	Team           string    `db:"team_name"`
	TeamID         string    `db:"team_id"`
	Games          int       `db:"games"`
	Goals          int       `db:"goals"`
	Assists        int       `db:"assists"`
	Points         int       `db:"points"`
	PlusMinus      int       `db:"plus_minus"`
	PenaltyMinutes int       `db:"penalty_minutes"`
}

// RankingsResult holds rankings with season info.
type RankingsResult struct {
	Season  string
	Players []RankedPlayerRow
}

// ExploreMatchesService provides match and ranking data.
type ExploreMatchesService struct {
	db *sqlx.DB
}

// NewExploreMatchesService creates a new explore matches service.
func NewExploreMatchesService(db *sqlx.DB) *ExploreMatchesService {
	return &ExploreMatchesService{db: db}
}

// GetRecentResults returns recently finished matches.
func (s *ExploreMatchesService) GetRecentResults(ctx context.Context, tournament string, limit int) ([]MatchRow, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.getMatches(ctx, "finished", tournament, limit, "m.scheduled_at DESC NULLS LAST")
}

// GetUpcomingMatches returns upcoming scheduled matches.
func (s *ExploreMatchesService) GetUpcomingMatches(ctx context.Context, tournament string, limit int) ([]MatchRow, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.getMatches(ctx, "scheduled", tournament, limit, "m.scheduled_at ASC NULLS LAST")
}

// GetTournamentMatches returns matches for a specific tournament with optional filters.
func (s *ExploreMatchesService) GetTournamentMatches(ctx context.Context, tournamentID string, birthYear int, groupName string, limit int) ([]MatchRow, error) {
	if limit <= 0 {
		limit = 30
	}

	where := []string{"m.tournament_id = $1", "m.home_team_id IS NOT NULL", "m.away_team_id IS NOT NULL"}
	args := []interface{}{tournamentID}
	argN := 2

	if birthYear > 0 {
		where = append(where, fmt.Sprintf("m.birth_year = $%d", argN))
		args = append(args, birthYear)
		argN++
	}
	if groupName != "" {
		where = append(where, fmt.Sprintf("m.group_name = $%d", argN))
		args = append(args, groupName)
		argN++
	}

	query := fmt.Sprintf(`
		SELECT m.id, COALESCE(ht.name, '') as home_team, COALESCE(at.name, '') as away_team,
			COALESCE(m.home_team_id, '') as home_team_id, COALESCE(m.away_team_id, '') as away_team_id,
			COALESCE(ht.logo_url, '') as home_logo_url, COALESCE(at.logo_url, '') as away_logo_url,
			m.home_score, m.away_score, COALESCE(m.result_type, '') as result_type, m.scheduled_at,
			COALESCE(t.name, '') as tournament_name, COALESCE(m.venue, '') as venue,
			COALESCE(m.status, 'scheduled') as status
		FROM matches m
		LEFT JOIN teams ht ON m.home_team_id = ht.id
		LEFT JOIN teams at ON m.away_team_id = at.id
		LEFT JOIN tournaments t ON m.tournament_id = t.id
		WHERE %s
		ORDER BY m.scheduled_at DESC NULLS LAST
		LIMIT $%d
	`, strings.Join(where, " AND "), argN)
	args = append(args, limit)

	var rows []MatchRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get tournament matches: %w", err)
	}
	for i := range rows {
		rows[i].HomeTeam = titleCase(rows[i].HomeTeam)
		rows[i].AwayTeam = titleCase(rows[i].AwayTeam)
		rows[i].Tournament = titleCase(rows[i].Tournament)
	}
	return rows, nil
}

func (s *ExploreMatchesService) getMatches(ctx context.Context, status, tournament string, limit int, orderBy string) ([]MatchRow, error) {
	where := []string{"m.status = $1", "m.home_team_id IS NOT NULL", "m.away_team_id IS NOT NULL"}
	args := []interface{}{status}
	argN := 2

	if tournament != "" {
		where = append(where, fmt.Sprintf("t.name ILIKE $%d", argN))
		args = append(args, "%"+tournament+"%")
		argN++
	}

	query := fmt.Sprintf(`
		SELECT m.id, COALESCE(ht.name, '') as home_team, COALESCE(at.name, '') as away_team,
			COALESCE(m.home_team_id, '') as home_team_id, COALESCE(m.away_team_id, '') as away_team_id,
			COALESCE(ht.logo_url, '') as home_logo_url, COALESCE(at.logo_url, '') as away_logo_url,
			m.home_score, m.away_score, COALESCE(m.result_type, '') as result_type, m.scheduled_at,
			COALESCE(t.name, '') as tournament_name, COALESCE(m.venue, '') as venue,
			COALESCE(m.status, 'scheduled') as status
		FROM matches m
		LEFT JOIN teams ht ON m.home_team_id = ht.id
		LEFT JOIN teams at ON m.away_team_id = at.id
		LEFT JOIN tournaments t ON m.tournament_id = t.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d
	`, strings.Join(where, " AND "), orderBy, argN)
	args = append(args, limit)

	var rows []MatchRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}

	for i := range rows {
		rows[i].HomeTeam = titleCase(rows[i].HomeTeam)
		rows[i].AwayTeam = titleCase(rows[i].AwayTeam)
		rows[i].Tournament = titleCase(rows[i].Tournament)
	}
	return rows, nil
}

// RankingsFilter holds optional filter parameters for rankings.
type RankingsFilter struct {
	BirthYear    int
	Domain       string
	TournamentID string
	GroupName    string
}

// GetRankings returns players ranked by a stat field for the current season.
func (s *ExploreMatchesService) GetRankings(ctx context.Context, sortBy string, limit int, filter RankingsFilter) (*RankingsResult, error) {
	if limit <= 0 {
		limit = 20
	}

	var season string
	if err := s.db.GetContext(ctx, &season, "SELECT season FROM tournaments ORDER BY season DESC LIMIT 1"); err != nil {
		return nil, fmt.Errorf("failed to get current season: %w", err)
	}

	allowedSorts := map[string]string{
		"points": "total_points", "goals": "total_goals", "assists": "total_assists",
		"plusMinus": "total_plus_minus", "penaltyMinutes": "total_penalty",
	}
	sortCol, ok := allowedSorts[sortBy]
	if !ok {
		sortCol = "total_points"
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.name, COALESCE(p.photo_url, '') as photo_url, COALESCE(p.position, '') as position, p.birth_date,
			COALESCE((SELECT t.name FROM player_teams pt JOIN teams t ON pt.team_id = t.id WHERE pt.player_id = p.id AND pt.is_active = true LIMIT 1), '') as team_name,
			COALESCE((SELECT pt.team_id FROM player_teams pt WHERE pt.player_id = p.id AND pt.is_active = true LIMIT 1), '') as team_id,
			COALESCE(SUM(ps.games), 0)::int as games,
			COALESCE(SUM(ps.goals), 0)::int as total_goals,
			COALESCE(SUM(ps.assists), 0)::int as total_assists,
			COALESCE(SUM(ps.points), 0)::int as total_points,
			COALESCE(SUM(ps.plus_minus), 0)::int as total_plus_minus,
			COALESCE(SUM(ps.penalty_minutes), 0)::int as total_penalty
		FROM players p
		JOIN player_statistics ps ON p.id = ps.player_id AND ps.group_name != 'Общая статистика'
		JOIN tournaments tr ON ps.tournament_id = tr.id AND tr.season = $1
			AND ($3 = 0 OR ps.birth_year = $3)
			AND ($4 = '' OR tr.domain = $4)
			AND ($5 = '' OR ps.tournament_id = $5)
			AND ($6 = '' OR ps.group_name = $6)
		GROUP BY p.id, p.name, p.position, p.birth_date
		HAVING SUM(ps.games) > 0
		ORDER BY %s DESC
		LIMIT $2
	`, sortCol)

	type rankedRow struct {
		ID             string    `db:"id"`
		Name           string    `db:"name"`
		PhotoURL       string    `db:"photo_url"`
		Position       string    `db:"position"`
		BirthDate      time.Time `db:"birth_date"`
		Team           string    `db:"team_name"`
		TeamID         string    `db:"team_id"`
		Games          int       `db:"games"`
		Goals          int       `db:"total_goals"`
		Assists        int       `db:"total_assists"`
		Points         int       `db:"total_points"`
		PlusMinus      int       `db:"total_plus_minus"`
		PenaltyMinutes int       `db:"total_penalty"`
	}

	var rows []rankedRow
	if err := s.db.SelectContext(ctx, &rows, query, season, limit, filter.BirthYear, filter.Domain, filter.TournamentID, filter.GroupName); err != nil {
		return nil, fmt.Errorf("failed to get rankings: %w", err)
	}

	players := make([]RankedPlayerRow, len(rows))
	for i, r := range rows {
		players[i] = RankedPlayerRow{
			ID: r.ID, Name: r.Name, PhotoURL: r.PhotoURL,
			Position: mapPositionToAPI(r.Position), BirthDate: r.BirthDate,
			Team: titleCase(r.Team), TeamID: r.TeamID,
			Games: r.Games, Goals: r.Goals, Assists: r.Assists, Points: r.Points,
			PlusMinus: r.PlusMinus, PenaltyMinutes: r.PenaltyMinutes,
		}
	}
	return &RankingsResult{Season: season, Players: players}, nil
}

// MatchDetailRow represents detailed match data from DB.
type MatchDetailRow struct {
	ID             string     `db:"id"`
	ExternalID     string     `db:"external_id"`
	HomeTeamID     string     `db:"home_team_id"`
	HomeTeamName   string     `db:"home_team_name"`
	HomeTeamCity   string     `db:"home_team_city"`
	HomeLogoURL    string     `db:"home_logo_url"`
	AwayTeamID     string     `db:"away_team_id"`
	AwayTeamName   string     `db:"away_team_name"`
	AwayTeamCity   string     `db:"away_team_city"`
	AwayLogoURL    string     `db:"away_logo_url"`
	HomeScore      *int       `db:"home_score"`
	AwayScore      *int       `db:"away_score"`
	HomeScoreP1    *int       `db:"home_score_p1"`
	AwayScoreP1    *int       `db:"away_score_p1"`
	HomeScoreP2    *int       `db:"home_score_p2"`
	AwayScoreP2    *int       `db:"away_score_p2"`
	HomeScoreP3    *int       `db:"home_score_p3"`
	AwayScoreP3    *int       `db:"away_score_p3"`
	HomeScoreOT    *int       `db:"home_score_ot"`
	AwayScoreOT    *int       `db:"away_score_ot"`
	ResultType     string     `db:"result_type"`
	ScheduledAt    *time.Time `db:"scheduled_at"`
	TournamentID   string     `db:"tournament_id"`
	TournamentName string     `db:"tournament_name"`
	Venue          string     `db:"venue"`
	Status         string     `db:"status"`
	GroupName      string     `db:"group_name"`
	BirthYear      *int       `db:"birth_year"`
	MatchNumber    *int       `db:"match_number"`
}

// MatchEventRow represents a match event from DB.
type MatchEventRow struct {
	ID            string  `db:"id"`
	EventType     string  `db:"event_type"`
	Period        *int    `db:"period"`
	TimeMinutes   *int    `db:"time_minutes"`
	TimeSeconds   *int    `db:"time_seconds"`
	IsHome        *bool   `db:"is_home"`
	TeamID        string  `db:"team_id"`
	TeamName      string  `db:"team_name"`
	TeamLogoURL   string  `db:"team_logo_url"`
	ScorerID      string  `db:"scorer_id"`
	ScorerName    string  `db:"scorer_name"`
	ScorerPhoto   string  `db:"scorer_photo"`
	Assist1ID     string  `db:"assist1_id"`
	Assist1Name   string  `db:"assist1_name"`
	Assist2ID     string  `db:"assist2_id"`
	Assist2Name   string  `db:"assist2_name"`
	GoalType      string  `db:"goal_type"`
	PenaltyMins   *int    `db:"penalty_minutes"`
	PenaltyReason string  `db:"penalty_reason"`
	PenaltyPlayerID   string `db:"penalty_player_id"`
	PenaltyPlayerName string `db:"penalty_player_name"`
	PenaltyPlayerPhoto string `db:"penalty_player_photo"`
}

// LineupRow represents a player in lineup from DB.
type LineupRow struct {
	PlayerID       string `db:"player_id"`
	PlayerName     string `db:"player_name"`
	PlayerPhoto    string `db:"player_photo"`
	JerseyNumber   *int   `db:"jersey_number"`
	Position       string `db:"position"`
	Goals          int    `db:"goals"`
	Assists        int    `db:"assists"`
	PenaltyMinutes int    `db:"penalty_minutes"`
	PlusMinus      int    `db:"plus_minus"`
	Saves          *int   `db:"saves"`
	GoalsAgainst   *int   `db:"goals_against"`
	IsHome         bool   `db:"is_home"`
}

// MatchDetailResult holds match detail with events and lineups.
type MatchDetailResult struct {
	Match   MatchDetailRow
	Events  []MatchEventRow
	Lineups []LineupRow
}

// GetMatchDetail returns detailed match information.
func (s *ExploreMatchesService) GetMatchDetail(ctx context.Context, id string) (*MatchDetailResult, error) {
	// Get match info
	matchQuery := `
		SELECT m.id, m.external_id,
			COALESCE(m.home_team_id, '') as home_team_id,
			COALESCE(ht.name, '') as home_team_name,
			COALESCE(ht.city, '') as home_team_city,
			COALESCE(ht.logo_url, '') as home_logo_url,
			COALESCE(m.away_team_id, '') as away_team_id,
			COALESCE(at.name, '') as away_team_name,
			COALESCE(at.city, '') as away_team_city,
			COALESCE(at.logo_url, '') as away_logo_url,
			m.home_score, m.away_score,
			m.home_score_p1, m.away_score_p1,
			m.home_score_p2, m.away_score_p2,
			m.home_score_p3, m.away_score_p3,
			m.home_score_ot, m.away_score_ot,
			COALESCE(m.result_type, '') as result_type,
			m.scheduled_at,
			COALESCE(m.tournament_id, '') as tournament_id,
			COALESCE(t.name, '') as tournament_name,
			COALESCE(m.venue, '') as venue,
			COALESCE(m.status, 'scheduled') as status,
			COALESCE(m.group_name, '') as group_name,
			m.birth_year, m.match_number
		FROM matches m
		LEFT JOIN teams ht ON m.home_team_id = ht.id
		LEFT JOIN teams at ON m.away_team_id = at.id
		LEFT JOIN tournaments t ON m.tournament_id = t.id
		WHERE m.id = $1
	`
	var match MatchDetailRow
	if err := s.db.GetContext(ctx, &match, matchQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// Get events
	eventsQuery := `
		SELECT me.id, me.event_type, me.period, me.time_minutes, me.time_seconds,
			me.is_home,
			COALESCE(me.team_id, '') as team_id,
			COALESCE(t.name, '') as team_name,
			COALESCE(t.logo_url, '') as team_logo_url,
			COALESCE(me.scorer_player_id, '') as scorer_id,
			COALESCE(sp.name, '') as scorer_name,
			COALESCE(sp.photo_url, '') as scorer_photo,
			COALESCE(me.assist1_player_id, '') as assist1_id,
			COALESCE(a1.name, '') as assist1_name,
			COALESCE(me.assist2_player_id, '') as assist2_id,
			COALESCE(a2.name, '') as assist2_name,
			COALESCE(me.goal_type, '') as goal_type,
			me.penalty_minutes,
			COALESCE(me.penalty_reason, '') as penalty_reason,
			COALESCE(me.penalty_player_id, '') as penalty_player_id,
			COALESCE(pp.name, '') as penalty_player_name,
			COALESCE(pp.photo_url, '') as penalty_player_photo
		FROM match_events me
		LEFT JOIN teams t ON me.team_id = t.id
		LEFT JOIN players sp ON me.scorer_player_id = sp.id
		LEFT JOIN players a1 ON me.assist1_player_id = a1.id
		LEFT JOIN players a2 ON me.assist2_player_id = a2.id
		LEFT JOIN players pp ON me.penalty_player_id = pp.id
		WHERE me.match_id = $1
		ORDER BY me.period ASC NULLS LAST, me.time_minutes ASC NULLS LAST, me.time_seconds ASC
	`
	var events []MatchEventRow
	if err := s.db.SelectContext(ctx, &events, eventsQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	// Get lineups
	lineupsQuery := `
		SELECT ml.player_id,
			COALESCE(p.name, '') as player_name,
			COALESCE(p.photo_url, '') as player_photo,
			ml.jersey_number,
			COALESCE(ml.position, '') as position,
			COALESCE(ml.goals, 0) as goals,
			COALESCE(ml.assists, 0) as assists,
			COALESCE(ml.penalty_minutes, 0) as penalty_minutes,
			COALESCE(ml.plus_minus, 0) as plus_minus,
			ml.saves,
			ml.goals_against,
			(ml.team_id = $2) as is_home
		FROM match_lineups ml
		JOIN players p ON ml.player_id = p.id
		WHERE ml.match_id = $1
		ORDER BY ml.position ASC, ml.goals DESC, ml.assists DESC
	`
	var lineups []LineupRow
	if err := s.db.SelectContext(ctx, &lineups, lineupsQuery, id, match.HomeTeamID); err != nil {
		return nil, fmt.Errorf("failed to get lineups: %w", err)
	}

	return &MatchDetailResult{
		Match:   match,
		Events:  events,
		Lineups: lineups,
	}, nil
}
