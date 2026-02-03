package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type MatchLineupPostgres struct {
	db *sqlx.DB
}

func NewMatchLineupPostgres(db *sqlx.DB) *MatchLineupPostgres {
	return &MatchLineupPostgres{db: db}
}

func (r *MatchLineupPostgres) Create(ctx context.Context, l *entities.MatchLineup) error {
	query := `
		INSERT INTO match_lineups (id, match_id, player_id, team_id, jersey_number, position, captain_role,
			goals, assists, penalty_minutes, plus_minus, saves, goals_against, time_on_ice,
			source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())
		ON CONFLICT (match_id, player_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.MatchID, l.PlayerID, l.TeamID, l.JerseyNumber, l.Position, l.CaptainRole,
		l.Goals, l.Assists, l.PenaltyMinutes, l.PlusMinus, l.Saves, l.GoalsAgainst, l.TimeOnIce,
		l.Source,
	)
	if err != nil {
		return fmt.Errorf("create match lineup: %w", err)
	}
	return nil
}

func (r *MatchLineupPostgres) CreateBatch(ctx context.Context, lineups []*entities.MatchLineup) error {
	if len(lineups) == 0 {
		return nil
	}

	query := `
		INSERT INTO match_lineups (id, match_id, player_id, team_id, jersey_number, position, captain_role,
			goals, assists, penalty_minutes, plus_minus, saves, goals_against, time_on_ice,
			source, created_at)
		VALUES (:id, :match_id, :player_id, :team_id, :jersey_number, :position, :captain_role,
			:goals, :assists, :penalty_minutes, :plus_minus, :saves, :goals_against, :time_on_ice,
			:source, NOW())
		ON CONFLICT (match_id, player_id) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, query, lineups)
	if err != nil {
		return fmt.Errorf("create match lineups batch: %w", err)
	}
	return nil
}

func (r *MatchLineupPostgres) Upsert(ctx context.Context, l *entities.MatchLineup) error {
	query := `
		INSERT INTO match_lineups (id, match_id, player_id, team_id, jersey_number, position, captain_role,
			goals, assists, penalty_minutes, plus_minus, saves, goals_against, time_on_ice,
			source, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())
		ON CONFLICT (match_id, player_id) DO UPDATE SET
			captain_role = EXCLUDED.captain_role,
			goals = EXCLUDED.goals,
			assists = EXCLUDED.assists,
			penalty_minutes = EXCLUDED.penalty_minutes,
			plus_minus = EXCLUDED.plus_minus,
			saves = EXCLUDED.saves,
			goals_against = EXCLUDED.goals_against,
			time_on_ice = EXCLUDED.time_on_ice`

	_, err := r.db.ExecContext(ctx, query,
		l.ID, l.MatchID, l.PlayerID, l.TeamID, l.JerseyNumber, l.Position, l.CaptainRole,
		l.Goals, l.Assists, l.PenaltyMinutes, l.PlusMinus, l.Saves, l.GoalsAgainst, l.TimeOnIce,
		l.Source,
	)
	if err != nil {
		return fmt.Errorf("upsert match lineup: %w", err)
	}
	return nil
}

func (r *MatchLineupPostgres) GetByMatchID(ctx context.Context, matchID string) ([]*entities.MatchLineup, error) {
	var lineups []*entities.MatchLineup
	err := r.db.SelectContext(ctx, &lineups,
		`SELECT * FROM match_lineups WHERE match_id = $1 ORDER BY team_id, position, jersey_number`, matchID)
	if err != nil {
		return nil, fmt.Errorf("get lineups by match: %w", err)
	}
	return lineups, nil
}

func (r *MatchLineupPostgres) GetByPlayerID(ctx context.Context, playerID string) ([]*entities.MatchLineup, error) {
	var lineups []*entities.MatchLineup
	err := r.db.SelectContext(ctx, &lineups,
		`SELECT * FROM match_lineups WHERE player_id = $1 ORDER BY created_at DESC`, playerID)
	if err != nil {
		return nil, fmt.Errorf("get lineups by player: %w", err)
	}
	return lineups, nil
}

func (r *MatchLineupPostgres) GetByMatchAndTeam(ctx context.Context, matchID, teamID string) ([]*entities.MatchLineup, error) {
	var lineups []*entities.MatchLineup
	err := r.db.SelectContext(ctx, &lineups,
		`SELECT * FROM match_lineups WHERE match_id = $1 AND team_id = $2 ORDER BY position, jersey_number`, matchID, teamID)
	if err != nil {
		return nil, fmt.Errorf("get lineups by match and team: %w", err)
	}
	return lineups, nil
}

func (r *MatchLineupPostgres) GetByMatchAndJersey(ctx context.Context, matchID string, jerseyNumber int) (*entities.MatchLineup, error) {
	var lineup entities.MatchLineup
	err := r.db.GetContext(ctx, &lineup,
		`SELECT * FROM match_lineups WHERE match_id = $1 AND jersey_number = $2 LIMIT 1`, matchID, jerseyNumber)
	if err != nil {
		return nil, fmt.Errorf("get lineup by match and jersey: %w", err)
	}
	return &lineup, nil
}

func (r *MatchLineupPostgres) DeleteByMatchID(ctx context.Context, matchID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM match_lineups WHERE match_id = $1`, matchID)
	if err != nil {
		return fmt.Errorf("delete lineups by match: %w", err)
	}
	return nil
}
