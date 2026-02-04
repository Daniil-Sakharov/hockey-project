package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// ExploreOverview represents platform-wide stats.
type ExploreOverview struct {
	Players     int64 `db:"players"`
	Teams       int64 `db:"teams"`
	Tournaments int64 `db:"tournaments"`
	Matches     int64 `db:"matches"`
}

// TournamentItem represents a tournament with counts.
type TournamentItem struct {
	ID                 string  `db:"id"`
	Name               string  `db:"name"`
	Domain             string  `db:"domain"`
	Season             string  `db:"season"`
	Source             string  `db:"source"`
	IsEnded            bool    `db:"is_ended"`
	BirthYearGroupsRaw *string `db:"birth_year_groups"`
	TeamsCount         int     `db:"teams_count"`
	MatchesCount       int     `db:"matches_count"`
}

// StandingRow represents a team standing row from DB.
type StandingRow struct {
	Position     int    `db:"position"`
	Team         string `db:"team_name"`
	TeamID       string `db:"team_id"`
	LogoURL      string `db:"logo_url"`
	Games        int    `db:"games"`
	Wins         int    `db:"wins"`
	WinsOT       int    `db:"wins_ot"`
	Losses       int    `db:"losses"`
	LossesOT     int    `db:"losses_ot"`
	Draws        int    `db:"draws"`
	GoalsFor     int    `db:"goals_for"`
	GoalsAgainst int    `db:"goals_against"`
	Points       int    `db:"points"`
	GroupName    string `db:"group_name"`
}

// ScorerRow represents a tournament scorer from DB.
type ScorerRow struct {
	PlayerID string `db:"player_id"`
	Name     string `db:"player_name"`
	PhotoURL string `db:"photo_url"`
	Team     string `db:"team_name"`
	TeamID   string `db:"team_id"`
	LogoURL  string `db:"logo_url"`
	Games    int    `db:"games"`
	Goals    int    `db:"goals"`
	Assists  int    `db:"assists"`
	Points   int    `db:"points"`
}

// ExploreService provides explore dashboard data.
type ExploreService struct {
	db *sqlx.DB
}

// NewExploreService creates a new explore service.
func NewExploreService(db *sqlx.DB) *ExploreService {
	return &ExploreService{db: db}
}

// GetOverview returns platform-wide statistics.
func (s *ExploreService) GetOverview(ctx context.Context) (*ExploreOverview, error) {
	var result ExploreOverview
	query := `
		SELECT
			(SELECT COUNT(*) FROM players) as players,
			(SELECT COUNT(*) FROM teams) as teams,
			(SELECT COUNT(*) FROM tournaments) as tournaments,
			(SELECT COUNT(*) FROM matches) as matches
	`
	if err := s.db.GetContext(ctx, &result, query); err != nil {
		return nil, fmt.Errorf("failed to get overview: %w", err)
	}
	return &result, nil
}

// GetTournaments returns list of tournaments with team/match counts.
func (s *ExploreService) GetTournaments(ctx context.Context, source, domain string) ([]TournamentItem, error) {
	query := `
		SELECT t.id, t.name, COALESCE(t.domain, '') as domain,
			COALESCE(t.season, '') as season, COALESCE(t.source, 'junior') as source,
			COALESCE(t.is_ended, false) as is_ended,
			t.birth_year_groups::text as birth_year_groups,
			(SELECT COUNT(*) FROM teams te WHERE te.tournament_id = t.id) as teams_count,
			(SELECT COUNT(*) FROM matches m WHERE m.tournament_id = t.id) as matches_count
		FROM tournaments t
	`
	args := []interface{}{}
	argN := 1
	var conditions []string

	if source != "" {
		conditions = append(conditions, "t.source = $"+strconv.Itoa(argN))
		args = append(args, source)
		argN++
	}
	if domain != "" {
		// Support multiple formats:
		// "ufo" -> "https://ufo.fhr.ru"
		// "ufo.fhr.ru" -> "https://ufo.fhr.ru"
		// "https://ufo.fhr.ru" -> "https://ufo.fhr.ru"
		domainURL := domain
		if !strings.Contains(domain, "://") {
			if strings.HasSuffix(domain, ".fhr.ru") {
				domainURL = "https://" + domain
			} else {
				domainURL = "https://" + domain + ".fhr.ru"
			}
		}
		conditions = append(conditions, "t.domain = $"+strconv.Itoa(argN))
		args = append(args, domainURL)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY t.is_ended ASC, t.name ASC"

	var items []TournamentItem
	if err := s.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get tournaments: %w", err)
	}

	for i := range items {
		items[i].Name = titleCase(items[i].Name)
		items[i].Domain = stripProtocol(items[i].Domain)
	}
	return items, nil
}

// GetTournamentStandings returns standings for a tournament with optional filters.
func (s *ExploreService) GetTournamentStandings(ctx context.Context, tournamentID string, birthYear int, groupName string) ([]StandingRow, error) {
	where := []string{"ts.tournament_id = $1"}
	args := []interface{}{tournamentID}
	argN := 2

	if birthYear > 0 {
		where = append(where, fmt.Sprintf("ts.birth_year = $%d", argN))
		args = append(args, birthYear)
		argN++
	}
	if groupName != "" {
		where = append(where, fmt.Sprintf("ts.group_name = $%d", argN))
		args = append(args, groupName)
	}

	query := fmt.Sprintf(`
		SELECT ts.position, t.name as team_name, ts.team_id,
			COALESCE(t.logo_url, '') as logo_url,
			COALESCE(ts.games, 0) as games, COALESCE(ts.wins, 0) as wins,
			COALESCE(ts.wins_ot, 0) as wins_ot, COALESCE(ts.losses, 0) as losses,
			COALESCE(ts.losses_ot, 0) as losses_ot, COALESCE(ts.draws, 0) as draws,
			COALESCE(ts.goals_for, 0) as goals_for, COALESCE(ts.goals_against, 0) as goals_against,
			COALESCE(ts.points, 0) as points, COALESCE(ts.group_name, '') as group_name
		FROM team_standings ts
		JOIN teams t ON ts.team_id = t.id
		WHERE %s
		ORDER BY ts.group_name, ts.position
	`, strings.Join(where, " AND "))

	var rows []StandingRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get standings: %w", err)
	}

	for i := range rows {
		rows[i].Team = titleCase(rows[i].Team)
	}
	return rows, nil
}

// GetTournamentScorers returns top scorers for a tournament with optional filters.
func (s *ExploreService) GetTournamentScorers(ctx context.Context, tournamentID string, birthYear int, groupName string, limit int) ([]ScorerRow, error) {
	if limit <= 0 {
		limit = 20
	}

	where := []string{"ps.tournament_id = $1"}
	args := []interface{}{tournamentID}
	argN := 2

	if groupName != "" {
		where = append(where, fmt.Sprintf("ps.group_name = $%d", argN))
		args = append(args, groupName)
		argN++
	} else {
		where = append(where, "ps.group_name = 'Общая статистика'")
	}

	if birthYear > 0 {
		where = append(where, fmt.Sprintf("ps.birth_year = $%d", argN))
		args = append(args, birthYear)
		argN++
	}

	query := fmt.Sprintf(`
		SELECT ps.player_id, p.name as player_name, COALESCE(p.photo_url, '') as photo_url,
			t.name as team_name, ps.team_id, COALESCE(t.logo_url, '') as logo_url,
			ps.games, ps.goals, ps.assists, ps.points
		FROM player_statistics ps
		JOIN players p ON ps.player_id = p.id
		JOIN teams t ON ps.team_id = t.id
		WHERE %s
		ORDER BY ps.points DESC, ps.goals DESC
		LIMIT $%d
	`, strings.Join(where, " AND "), argN)
	args = append(args, limit)

	var rows []ScorerRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get scorers: %w", err)
	}

	for i := range rows {
		rows[i].Team = titleCase(rows[i].Team)
	}
	return rows, nil
}

// GroupStats represents team/match counts for a tournament group.
type GroupStats struct {
	TournamentID string `db:"tournament_id"`
	BirthYear    int    `db:"birth_year"`
	GroupName    string `db:"group_name"`
	TeamsCount   int    `db:"teams_count"`
	MatchesCount int    `db:"matches_count"`
}

// GetGroupStats returns team/match counts per (tournament, birth_year, group_name).
func (s *ExploreService) GetGroupStats(ctx context.Context, tournamentIDs []string) ([]GroupStats, error) {
	if len(tournamentIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`
		SELECT gs.tournament_id, gs.birth_year, gs.group_name, gs.teams_count,
			COALESCE(mc.matches_count, 0) as matches_count
		FROM (
			SELECT tournament_id, birth_year, group_name,
				COUNT(DISTINCT team_id) as teams_count
			FROM team_standings
			WHERE tournament_id IN (?) AND birth_year IS NOT NULL AND group_name IS NOT NULL AND group_name != ''
			GROUP BY tournament_id, birth_year, group_name
		) gs
		LEFT JOIN (
			SELECT tournament_id, birth_year, group_name,
				COUNT(*) as matches_count
			FROM matches
			WHERE tournament_id IN (?) AND birth_year IS NOT NULL AND group_name IS NOT NULL AND group_name != ''
			GROUP BY tournament_id, birth_year, group_name
		) mc ON gs.tournament_id = mc.tournament_id
			AND gs.birth_year = mc.birth_year AND gs.group_name = mc.group_name
	`, tournamentIDs, tournamentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to build group stats query: %w", err)
	}
	query = s.db.Rebind(query)

	var rows []GroupStats
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get group stats: %w", err)
	}
	return rows, nil
}

// GetSeasons returns all available seasons sorted desc.
func (s *ExploreService) GetSeasons(ctx context.Context) ([]string, error) {
	var seasons []string
	query := `SELECT DISTINCT season FROM tournaments WHERE season != '' ORDER BY season DESC`
	if err := s.db.SelectContext(ctx, &seasons, query); err != nil {
		return nil, fmt.Errorf("failed to get seasons: %w", err)
	}
	return seasons, nil
}

// abbreviations that should keep their uppercase form.
var knownAbbreviations = map[string]string{
	"пфо":  "ПФО",
	"ппфо": "ППФО",
	"урфо": "УрФО",
	"цфо":  "ЦФО",
	"сзфо": "СЗФО",
	"юфо":  "ЮФО",
	"сфо":  "СФО",
	"дфо":  "ДФО",
	"мо":   "МО",
	"ло":   "ЛО",
	"спб":  "СПб",
	"3х3":  "3х3",
}

// titleCase converts "ТОРПЕДО" → "Торпедо", preserving known abbreviations.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) == 0 {
			continue
		}
		lower := strings.ToLower(w)
		if abbr, ok := knownAbbreviations[lower]; ok {
			words[i] = abbr
			continue
		}
		runes := []rune(lower)
		runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// stripProtocol removes protocol prefix from domain.
func stripProtocol(domain string) string {
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	return domain
}

// TeamRow represents a team in tournament with player count.
type TeamRow struct {
	ID           string `db:"id"`
	Name         string `db:"name"`
	City         string `db:"city"`
	LogoURL      string `db:"logo_url"`
	PlayersCount int    `db:"players_count"`
	GroupName    string `db:"group_name"`
	BirthYear    int    `db:"birth_year"`
}

// GetTournamentTeams returns teams for a tournament with optional filters.
func (s *ExploreService) GetTournamentTeams(ctx context.Context, tournamentID string, birthYear int, groupName string) ([]TeamRow, error) {
	query := `
		SELECT
			t.id, t.name, COALESCE(t.city, '') as city, COALESCE(t.logo_url, '') as logo_url,
			pt.group_name, pt.birth_year,
			COUNT(DISTINCT pt.player_id) as players_count
		FROM teams t
		JOIN player_teams pt ON pt.team_id = t.id
		WHERE pt.tournament_id = $1
	`
	args := []interface{}{tournamentID}
	argN := 2

	if birthYear > 0 {
		query += fmt.Sprintf(" AND pt.birth_year = $%d", argN)
		args = append(args, birthYear)
		argN++
	}
	if groupName != "" {
		query += fmt.Sprintf(" AND pt.group_name = $%d", argN)
		args = append(args, groupName)
	}

	query += " GROUP BY t.id, t.name, t.city, t.logo_url, pt.group_name, pt.birth_year"
	query += " ORDER BY t.name ASC"

	var rows []TeamRow
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get tournament teams: %w", err)
	}

	for i := range rows {
		rows[i].Name = titleCase(rows[i].Name)
		rows[i].City = titleCase(rows[i].City)
	}
	return rows, nil
}
