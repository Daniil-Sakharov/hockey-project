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

// PlayerSearchRow represents a player search result from DB.
type PlayerSearchRow struct {
	ID           string `db:"id"`
	Name         string `db:"name"`
	Position     string `db:"position"`
	BirthDate    time.Time `db:"birth_date"`
	Team         string `db:"team_name"`
	TeamID       string `db:"team_id"`
	JerseyNumber int    `db:"jersey_number"`
	PhotoURL     string `db:"photo_url"`
	Games        int    `db:"games"`
	Goals        int    `db:"goals"`
	Assists      int    `db:"assists"`
	Points       int    `db:"points"`
	PlusMinus    int    `db:"plus_minus"`
	PenaltyMins  int    `db:"penalty_minutes"`
}

// PlayerProfileRow represents a full player profile from DB.
type PlayerProfileRow struct {
	PlayerSearchRow
	Height     *int   `db:"height"`
	Weight     *int   `db:"weight"`
	Handedness string `db:"handedness"`
	BirthPlace string `db:"birth_place"`
	PhotoURL   string `db:"photo_url"`
}

// ExplorePlayersService provides player/team explore data.
type ExplorePlayersService struct {
	db *sqlx.DB
}

// NewExplorePlayersService creates a new explore players service.
func NewExplorePlayersService(db *sqlx.DB) *ExplorePlayersService {
	return &ExplorePlayersService{db: db}
}

// SearchPlayers searches players by name, position, birth year.
func (s *ExplorePlayersService) SearchPlayers(ctx context.Context, q, position, season string, birthYear, limit, offset int) ([]PlayerSearchRow, int, error) {
	if limit <= 0 {
		limit = 20
	}

	where := []string{"1=1"}
	args := []interface{}{}
	argN := 1

	if q != "" {
		where = append(where, fmt.Sprintf("p.name ILIKE $%d", argN))
		args = append(args, "%"+q+"%")
		argN++
	}
	if position != "" {
		dbPos := mapPositionToDB(position)
		where = append(where, fmt.Sprintf("p.position = $%d", argN))
		args = append(args, dbPos)
		argN++
	}
	if birthYear > 0 {
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM p.birth_date) = $%d", argN))
		args = append(args, birthYear)
		argN++
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT p.id) FROM players p WHERE %s", whereClause)
	var total int
	if err := s.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count players: %w", err)
	}

	// Build season filter for stats subquery
	seasonFilter := ""
	if season != "" {
		seasonFilter = fmt.Sprintf("AND tt.season = $%d", argN)
		args = append(args, season)
		argN++
	}

	// Fetch players with latest team + season-aggregated stats
	query := fmt.Sprintf(`
		SELECT p.id, p.name, COALESCE(p.position, '') as position, p.birth_date,
			COALESCE(t.name, '') as team_name, COALESCE(t.id, '') as team_id,
			COALESCE(pt.jersey_number, 0) as jersey_number,
			COALESCE(p.photo_url, '') as photo_url,
			COALESCE(stats.games, 0) as games, COALESCE(stats.goals, 0) as goals,
			COALESCE(stats.assists, 0) as assists, COALESCE(stats.points, 0) as points,
			COALESCE(stats.plus_minus, 0) as plus_minus, COALESCE(stats.penalty_minutes, 0) as penalty_minutes
		FROM players p
		LEFT JOIN LATERAL (
			SELECT team_id, jersey_number FROM player_teams
			WHERE player_id = p.id ORDER BY is_active DESC, tournament_id DESC LIMIT 1
		) pt ON true
		LEFT JOIN teams t ON pt.team_id = t.id
		LEFT JOIN LATERAL (
			SELECT
				SUM(ps.games)::int as games, SUM(ps.goals)::int as goals,
				SUM(ps.assists)::int as assists, SUM(ps.points)::int as points,
				SUM(ps.plus_minus)::int as plus_minus, SUM(ps.penalty_minutes)::int as penalty_minutes
			FROM player_statistics ps
			JOIN tournaments tt ON ps.tournament_id = tt.id
			WHERE ps.player_id = p.id AND ps.group_name = 'Общая статистика' %s
		) stats ON true
		WHERE %s
		ORDER BY stats.points DESC NULLS LAST, p.id
		LIMIT $%d OFFSET $%d
	`, seasonFilter, whereClause, argN, argN+1)
	args = append(args, limit, offset)

	var rows []PlayerSearchRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to search players: %w", err)
	}

	for i := range rows {
		rows[i].Position = mapPositionToAPI(rows[i].Position)
		rows[i].Team = titleCase(rows[i].Team)
	}
	return rows, total, nil
}

// GetPlayerProfile returns a full player profile.
func (s *ExplorePlayersService) GetPlayerProfile(ctx context.Context, id, season string) (*PlayerProfileRow, error) {
	args := []interface{}{id}
	seasonFilter := ""
	if season != "" {
		seasonFilter = "JOIN tournaments tt ON ps.tournament_id = tt.id AND tt.season = $2"
		args = append(args, season)
	}

	query := fmt.Sprintf(`
		WITH player_team AS (
			SELECT player_id, team_id, jersey_number
			FROM player_teams
			WHERE player_id = $1
			ORDER BY is_active DESC, tournament_id DESC
			LIMIT 1
		)
		SELECT p.id, p.name, COALESCE(p.position, '') as position, p.birth_date,
			COALESCE(t.name, '') as team_name, COALESCE(t.id, '') as team_id,
			COALESCE(pt.jersey_number, 0) as jersey_number,
			p.height, p.weight, COALESCE(p.handedness, '') as handedness,
			COALESCE(p.birth_place, '') as birth_place, COALESCE(p.photo_url, '') as photo_url,
			COALESCE(stats.games, 0) as games, COALESCE(stats.goals, 0) as goals,
			COALESCE(stats.assists, 0) as assists, COALESCE(stats.points, 0) as points,
			COALESCE(stats.plus_minus, 0) as plus_minus, COALESCE(stats.penalty_minutes, 0) as penalty_minutes
		FROM players p
		LEFT JOIN player_team pt ON p.id = pt.player_id
		LEFT JOIN teams t ON pt.team_id = t.id
		LEFT JOIN LATERAL (
			SELECT
				SUM(ps.games)::int as games, SUM(ps.goals)::int as goals,
				SUM(ps.assists)::int as assists, SUM(ps.points)::int as points,
				SUM(ps.plus_minus)::int as plus_minus, SUM(ps.penalty_minutes)::int as penalty_minutes
			FROM player_statistics ps
			%s
			WHERE ps.player_id = p.id AND ps.group_name = 'Общая статистика'
		) stats ON true
		WHERE p.id = $1
		LIMIT 1
	`, seasonFilter)
	var row PlayerProfileRow
	if err := s.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	row.Position = mapPositionToAPI(row.Position)
	row.Handedness = mapHandednessToAPI(row.Handedness)
	row.Team = titleCase(row.Team)
	return &row, nil
}

// PlayerStatRow represents a detailed stat entry for a player.
type PlayerStatRow struct {
	Season         string `db:"season"`
	TournamentID   string `db:"tournament_id"`
	TournamentName string `db:"tournament_name"`
	GroupName      string `db:"group_name"`
	BirthYear      int    `db:"birth_year"`
	Games          int    `db:"games"`
	Goals          int    `db:"goals"`
	Assists        int    `db:"assists"`
	Points         int    `db:"points"`
	PlusMinus      int    `db:"plus_minus"`
	PenaltyMinutes int    `db:"penalty_minutes"`
}

// GetPlayerStats returns detailed stats for a player across all seasons/tournaments/groups.
func (s *ExplorePlayersService) GetPlayerStats(ctx context.Context, id string) ([]PlayerStatRow, error) {
	query := `
		SELECT t.season, ps.tournament_id, t.name as tournament_name,
			ps.group_name, COALESCE(ps.birth_year, 0) as birth_year,
			COALESCE(ps.games, 0) as games, COALESCE(ps.goals, 0) as goals,
			COALESCE(ps.assists, 0) as assists, COALESCE(ps.points, 0) as points,
			COALESCE(ps.plus_minus, 0) as plus_minus, COALESCE(ps.penalty_minutes, 0) as penalty_minutes
		FROM player_statistics ps
		JOIN tournaments t ON ps.tournament_id = t.id
		WHERE ps.player_id = $1
		ORDER BY t.season DESC, t.name, ps.group_name
	`
	var rows []PlayerStatRow
	if err := s.db.SelectContext(ctx, &rows, query, id); err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}
	for i := range rows {
		rows[i].TournamentName = titleCase(rows[i].TournamentName)
	}
	return rows, nil
}

// TeamProfileData holds team profile data.
type TeamProfileData struct {
	ID          string
	Name        string
	City        string
	LogoURL     string
	Tournaments []string
	Roster      []PlayerSearchRow
	Stats       TeamStats
	Matches     []MatchRow
}

// TeamStats holds aggregated team stats.
type TeamStats struct {
	Wins         int
	Losses       int
	Draws        int
	GoalsFor     int
	GoalsAgainst int
}

// GetTeamProfile returns a team profile with roster and matches.
func (s *ExplorePlayersService) GetTeamProfile(ctx context.Context, id string) (*TeamProfileData, error) {
	// Team info
	type teamRow struct {
		ID      string `db:"id"`
		Name    string `db:"name"`
		City    string `db:"city"`
		LogoURL string `db:"logo_url"`
	}
	var team teamRow
	teamQuery := `SELECT id, name, COALESCE(city, '') as city, COALESCE(logo_url, '') as logo_url FROM teams WHERE id = $1`
	if err := s.db.GetContext(ctx, &team, teamQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Tournaments
	var tournaments []string
	tourQuery := `SELECT DISTINCT t.name FROM tournaments t JOIN teams te ON te.tournament_id = t.id WHERE te.id = $1`
	if err := s.db.SelectContext(ctx, &tournaments, tourQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get team tournaments: %w", err)
	}

	// Roster
	rosterQuery := `
		SELECT DISTINCT ON (p.id) p.id, p.name, COALESCE(p.position, '') as position, p.birth_date,
			COALESCE(t.name, '') as team_name, COALESCE(t.id, '') as team_id,
			COALESCE(pt.jersey_number, 0) as jersey_number,
			COALESCE(p.photo_url, '') as photo_url,
			COALESCE(ps.games, 0) as games, COALESCE(ps.goals, 0) as goals,
			COALESCE(ps.assists, 0) as assists, COALESCE(ps.points, 0) as points,
			COALESCE(ps.plus_minus, 0) as plus_minus, COALESCE(ps.penalty_minutes, 0) as penalty_minutes
		FROM players p
		JOIN player_teams pt ON p.id = pt.player_id AND pt.team_id = $1
		LEFT JOIN teams t ON pt.team_id = t.id
		LEFT JOIN player_statistics ps ON p.id = ps.player_id AND ps.team_id = $1 AND ps.group_name = 'Общая статистика'
		ORDER BY p.id, ps.points DESC NULLS LAST
	`
	var roster []PlayerSearchRow
	if err := s.db.SelectContext(ctx, &roster, rosterQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get roster: %w", err)
	}
	for i := range roster {
		roster[i].Position = mapPositionToAPI(roster[i].Position)
		roster[i].Team = titleCase(roster[i].Team)
	}

	// Standings aggregated
	var stats TeamStats
	statsQuery := `
		SELECT COALESCE(SUM(wins), 0) as wins, COALESCE(SUM(losses), 0) as losses,
			COALESCE(SUM(draws), 0) as draws, COALESCE(SUM(goals_for), 0) as goals_for,
			COALESCE(SUM(goals_against), 0) as goals_against
		FROM team_standings WHERE team_id = $1
	`
	if err := s.db.GetContext(ctx, &stats, statsQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get team stats: %w", err)
	}

	// Recent matches
	matchesQuery := `
		SELECT m.id, COALESCE(ht.name, '') as home_team, COALESCE(at.name, '') as away_team,
			COALESCE(m.home_team_id, '') as home_team_id, COALESCE(m.away_team_id, '') as away_team_id,
			COALESCE(ht.logo_url, '') as home_logo_url, COALESCE(at.logo_url, '') as away_logo_url,
			m.home_score, m.away_score, COALESCE(m.result_type, '') as result_type, m.scheduled_at,
			COALESCE(tr.name, '') as tournament_name, COALESCE(m.venue, '') as venue,
			COALESCE(m.status, 'scheduled') as status
		FROM matches m
		LEFT JOIN teams ht ON m.home_team_id = ht.id
		LEFT JOIN teams at ON m.away_team_id = at.id
		LEFT JOIN tournaments tr ON m.tournament_id = tr.id
		WHERE (m.home_team_id = $1 OR m.away_team_id = $1) AND m.status = 'finished'
		ORDER BY m.scheduled_at DESC NULLS LAST
		LIMIT 5
	`
	var matches []MatchRow
	if err := s.db.SelectContext(ctx, &matches, matchesQuery, id); err != nil {
		return nil, fmt.Errorf("failed to get team matches: %w", err)
	}
	for i := range matches {
		matches[i].HomeTeam = titleCase(matches[i].HomeTeam)
		matches[i].AwayTeam = titleCase(matches[i].AwayTeam)
	}

	return &TeamProfileData{
		ID:          team.ID,
		Name:        titleCase(team.Name),
		City:        team.City,
		LogoURL:     team.LogoURL,
		Tournaments: tournaments,
		Roster:      roster,
		Stats:       stats,
		Matches:     matches,
	}, nil
}

// mapPositionToDB maps API position to DB value.
func mapPositionToDB(pos string) string {
	switch strings.ToLower(pos) {
	case "forward":
		return "Нападающий"
	case "defender":
		return "Защитник"
	case "goalie":
		return "Вратарь"
	default:
		return pos
	}
}

// mapPositionToAPI maps DB position to API value.
func mapPositionToAPI(pos string) string {
	switch pos {
	case "Нападающий":
		return "forward"
	case "Защитник":
		return "defender"
	case "Вратарь":
		return "goalie"
	default:
		return pos
	}
}

// mapHandednessToAPI maps DB handedness to API value.
func mapHandednessToAPI(h string) string {
	switch h {
	case "Левый":
		return "left"
	case "Правый":
		return "right"
	default:
		return h
	}
}
