package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type MatchTeamStatsPostgres struct {
	db *sqlx.DB
}

func NewMatchTeamStatsPostgres(db *sqlx.DB) *MatchTeamStatsPostgres {
	return &MatchTeamStatsPostgres{db: db}
}

func (r *MatchTeamStatsPostgres) Create(ctx context.Context, s *entities.MatchTeamStats) error {
	query := `
		INSERT INTO match_team_stats (id, match_id, team_id, shots_p1, shots_p2, shots_p3, shots_ot, shots_total, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (match_id, team_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.MatchID, s.TeamID, s.ShotsP1, s.ShotsP2, s.ShotsP3, s.ShotsOT, s.ShotsTotal, s.Source,
	)
	if err != nil {
		return fmt.Errorf("create match team stats: %w", err)
	}
	return nil
}

func (r *MatchTeamStatsPostgres) CreateBatch(ctx context.Context, stats []*entities.MatchTeamStats) error {
	if len(stats) == 0 {
		return nil
	}

	query := `
		INSERT INTO match_team_stats (id, match_id, team_id, shots_p1, shots_p2, shots_p3, shots_ot, shots_total, source, created_at)
		VALUES (:id, :match_id, :team_id, :shots_p1, :shots_p2, :shots_p3, :shots_ot, :shots_total, :source, NOW())
		ON CONFLICT (match_id, team_id) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, query, stats)
	if err != nil {
		return fmt.Errorf("create match team stats batch: %w", err)
	}
	return nil
}

func (r *MatchTeamStatsPostgres) Upsert(ctx context.Context, s *entities.MatchTeamStats) error {
	query := `
		INSERT INTO match_team_stats (id, match_id, team_id, shots_p1, shots_p2, shots_p3, shots_ot, shots_total, source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (match_id, team_id) DO UPDATE SET
			shots_p1 = EXCLUDED.shots_p1,
			shots_p2 = EXCLUDED.shots_p2,
			shots_p3 = EXCLUDED.shots_p3,
			shots_ot = EXCLUDED.shots_ot,
			shots_total = EXCLUDED.shots_total`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.MatchID, s.TeamID, s.ShotsP1, s.ShotsP2, s.ShotsP3, s.ShotsOT, s.ShotsTotal, s.Source,
	)
	if err != nil {
		return fmt.Errorf("upsert match team stats: %w", err)
	}
	return nil
}

func (r *MatchTeamStatsPostgres) UpsertBatch(ctx context.Context, stats []*entities.MatchTeamStats) error {
	for _, s := range stats {
		if err := r.Upsert(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (r *MatchTeamStatsPostgres) GetByMatch(ctx context.Context, matchID string) ([]*entities.MatchTeamStats, error) {
	var stats []*entities.MatchTeamStats
	err := r.db.SelectContext(ctx, &stats,
		`SELECT * FROM match_team_stats WHERE match_id = $1`, matchID)
	if err != nil {
		return nil, fmt.Errorf("get match team stats by match: %w", err)
	}
	return stats, nil
}

func (r *MatchTeamStatsPostgres) GetByMatchAndTeam(ctx context.Context, matchID, teamID string) (*entities.MatchTeamStats, error) {
	var s entities.MatchTeamStats
	err := r.db.GetContext(ctx, &s,
		`SELECT * FROM match_team_stats WHERE match_id = $1 AND team_id = $2`, matchID, teamID)
	if err != nil {
		return nil, fmt.Errorf("get match team stats by match and team: %w", err)
	}
	return &s, nil
}

func (r *MatchTeamStatsPostgres) DeleteByMatch(ctx context.Context, matchID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM match_team_stats WHERE match_id = $1`, matchID)
	if err != nil {
		return fmt.Errorf("delete match team stats by match: %w", err)
	}
	return nil
}
