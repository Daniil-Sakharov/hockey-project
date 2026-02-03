package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type MatchPostgres struct {
	db *sqlx.DB
}

func NewMatchPostgres(db *sqlx.DB) *MatchPostgres {
	return &MatchPostgres{db: db}
}

func (r *MatchPostgres) Create(ctx context.Context, m *entities.Match) error {
	query := `
		INSERT INTO matches (id, external_id, tournament_id, home_team_id, away_team_id,
			home_score, away_score, home_score_p1, away_score_p1, home_score_p2, away_score_p2,
			home_score_p3, away_score_p3, home_score_ot, away_score_ot, match_number,
			scheduled_at, status, result_type, venue, group_name, birth_year, video_url,
			source, domain, details_parsed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, NOW(), NOW())
		ON CONFLICT (source, external_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.ExternalID, m.TournamentID, m.HomeTeamID, m.AwayTeamID,
		m.HomeScore, m.AwayScore, m.HomeScoreP1, m.AwayScoreP1, m.HomeScoreP2, m.AwayScoreP2,
		m.HomeScoreP3, m.AwayScoreP3, m.HomeScoreOT, m.AwayScoreOT, m.MatchNumber,
		m.ScheduledAt, m.Status, m.ResultType, m.Venue, m.GroupName, m.BirthYear, m.VideoURL,
		m.Source, m.Domain, m.DetailsParsed,
	)
	if err != nil {
		return fmt.Errorf("create match: %w", err)
	}
	return nil
}

func (r *MatchPostgres) CreateBatch(ctx context.Context, matches []*entities.Match) error {
	if len(matches) == 0 {
		return nil
	}

	query := `
		INSERT INTO matches (id, external_id, tournament_id, home_team_id, away_team_id,
			home_score, away_score, scheduled_at, status, group_name, birth_year,
			source, domain, created_at, updated_at)
		VALUES (:id, :external_id, :tournament_id, :home_team_id, :away_team_id,
			:home_score, :away_score, :scheduled_at, :status, :group_name, :birth_year,
			:source, :domain, NOW(), NOW())
		ON CONFLICT (source, external_id) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, query, matches)
	if err != nil {
		return fmt.Errorf("create matches batch: %w", err)
	}
	return nil
}

func (r *MatchPostgres) Update(ctx context.Context, m *entities.Match) error {
	query := `
		UPDATE matches SET
			home_score = $1, away_score = $2, home_score_p1 = $3, away_score_p1 = $4,
			home_score_p2 = $5, away_score_p2 = $6, home_score_p3 = $7, away_score_p3 = $8,
			home_score_ot = $9, away_score_ot = $10, status = $11, result_type = $12,
			video_url = $13, details_parsed = $14, updated_at = NOW()
		WHERE id = $15`

	_, err := r.db.ExecContext(ctx, query,
		m.HomeScore, m.AwayScore, m.HomeScoreP1, m.AwayScoreP1,
		m.HomeScoreP2, m.AwayScoreP2, m.HomeScoreP3, m.AwayScoreP3,
		m.HomeScoreOT, m.AwayScoreOT, m.Status, m.ResultType,
		m.VideoURL, m.DetailsParsed, m.ID,
	)
	if err != nil {
		return fmt.Errorf("update match: %w", err)
	}
	return nil
}

func (r *MatchPostgres) Upsert(ctx context.Context, m *entities.Match) error {
	query := `
		INSERT INTO matches (id, external_id, tournament_id, home_team_id, away_team_id,
			home_score, away_score, home_score_p1, away_score_p1, home_score_p2, away_score_p2,
			home_score_p3, away_score_p3, home_score_ot, away_score_ot, match_number,
			scheduled_at, status, result_type, venue, group_name, birth_year, video_url,
			source, domain, details_parsed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, NOW(), NOW())
		ON CONFLICT (source, external_id) DO UPDATE SET
			home_team_id = COALESCE(EXCLUDED.home_team_id, matches.home_team_id),
			away_team_id = COALESCE(EXCLUDED.away_team_id, matches.away_team_id),
			home_score = COALESCE(EXCLUDED.home_score, matches.home_score),
			away_score = COALESCE(EXCLUDED.away_score, matches.away_score),
			match_number = COALESCE(EXCLUDED.match_number, matches.match_number),
			scheduled_at = COALESCE(EXCLUDED.scheduled_at, matches.scheduled_at),
			status = EXCLUDED.status,
			result_type = COALESCE(EXCLUDED.result_type, matches.result_type),
			venue = COALESCE(EXCLUDED.venue, matches.venue),
			group_name = COALESCE(EXCLUDED.group_name, matches.group_name),
			birth_year = COALESCE(EXCLUDED.birth_year, matches.birth_year),
			video_url = COALESCE(EXCLUDED.video_url, matches.video_url),
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.ExternalID, m.TournamentID, m.HomeTeamID, m.AwayTeamID,
		m.HomeScore, m.AwayScore, m.HomeScoreP1, m.AwayScoreP1, m.HomeScoreP2, m.AwayScoreP2,
		m.HomeScoreP3, m.AwayScoreP3, m.HomeScoreOT, m.AwayScoreOT, m.MatchNumber,
		m.ScheduledAt, m.Status, m.ResultType, m.Venue, m.GroupName, m.BirthYear, m.VideoURL,
		m.Source, m.Domain, m.DetailsParsed,
	)
	if err != nil {
		return fmt.Errorf("upsert match: %w", err)
	}
	return nil
}

func (r *MatchPostgres) GetByID(ctx context.Context, id string) (*entities.Match, error) {
	var m entities.Match
	err := r.db.GetContext(ctx, &m, `SELECT * FROM matches WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get match by id: %w", err)
	}
	return &m, nil
}

func (r *MatchPostgres) GetByExternalID(ctx context.Context, externalID, source string) (*entities.Match, error) {
	var m entities.Match
	err := r.db.GetContext(ctx, &m, `SELECT * FROM matches WHERE external_id = $1 AND source = $2`, externalID, source)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get match by external id: %w", err)
	}
	return &m, nil
}

func (r *MatchPostgres) GetByTournament(ctx context.Context, tournamentID string) ([]*entities.Match, error) {
	var matches []*entities.Match
	err := r.db.SelectContext(ctx, &matches,
		`SELECT * FROM matches WHERE tournament_id = $1 ORDER BY scheduled_at`, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("get matches by tournament: %w", err)
	}
	return matches, nil
}

func (r *MatchPostgres) GetByTeam(ctx context.Context, teamID string) ([]*entities.Match, error) {
	var matches []*entities.Match
	err := r.db.SelectContext(ctx, &matches,
		`SELECT * FROM matches WHERE home_team_id = $1 OR away_team_id = $1 ORDER BY scheduled_at`, teamID)
	if err != nil {
		return nil, fmt.Errorf("get matches by team: %w", err)
	}
	return matches, nil
}

func (r *MatchPostgres) GetUnparsedFinished(ctx context.Context, source string, limit int) ([]*entities.Match, error) {
	var matches []*entities.Match
	query := `SELECT * FROM matches WHERE source = $1 AND status = 'finished' AND details_parsed = false
		ORDER BY scheduled_at DESC`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	err := r.db.SelectContext(ctx, &matches, query, source)
	if err != nil {
		return nil, fmt.Errorf("get unparsed finished matches: %w", err)
	}
	return matches, nil
}

func (r *MatchPostgres) MarkDetailsParsed(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE matches SET details_parsed = true, updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("mark details parsed: %w", err)
	}
	return nil
}
