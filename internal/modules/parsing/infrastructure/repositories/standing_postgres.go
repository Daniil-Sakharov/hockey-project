package repositories

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/jmoiron/sqlx"
)

type StandingPostgres struct {
	db *sqlx.DB
}

func NewStandingPostgres(db *sqlx.DB) *StandingPostgres {
	return &StandingPostgres{db: db}
}

func (r *StandingPostgres) Create(ctx context.Context, s *entities.TeamStanding) error {
	query := `
		INSERT INTO team_standings (id, tournament_id, team_id, position, points,
			games, wins, wins_ot, wins_so, losses_so, losses_ot, losses, draws,
			goals_for, goals_against, goal_difference, group_name, birth_year,
			source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), NOW())
		ON CONFLICT (tournament_id, team_id, COALESCE(birth_year, 0), COALESCE(group_name, ''), source) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.TournamentID, s.TeamID, s.Position, s.Points,
		s.Games, s.Wins, s.WinsOT, s.WinsSO, s.LossesSO, s.LossesOT, s.Losses, s.Draws,
		s.GoalsFor, s.GoalsAgainst, s.GoalDifference, s.GroupName, s.BirthYear, s.Source,
	)
	if err != nil {
		return fmt.Errorf("create standing: %w", err)
	}
	return nil
}

func (r *StandingPostgres) CreateBatch(ctx context.Context, standings []*entities.TeamStanding) error {
	if len(standings) == 0 {
		return nil
	}

	query := `
		INSERT INTO team_standings (id, tournament_id, team_id, position, points,
			games, wins, wins_ot, wins_so, losses_so, losses_ot, losses, draws,
			goals_for, goals_against, goal_difference, group_name, birth_year,
			source, created_at, updated_at)
		VALUES (:id, :tournament_id, :team_id, :position, :points,
			:games, :wins, :wins_ot, :wins_so, :losses_so, :losses_ot, :losses, :draws,
			:goals_for, :goals_against, :goal_difference, :group_name, :birth_year,
			:source, NOW(), NOW())
		ON CONFLICT (tournament_id, team_id, COALESCE(birth_year, 0), COALESCE(group_name, ''), source) DO NOTHING`

	_, err := r.db.NamedExecContext(ctx, query, standings)
	if err != nil {
		return fmt.Errorf("create standings batch: %w", err)
	}
	return nil
}

func (r *StandingPostgres) Upsert(ctx context.Context, s *entities.TeamStanding) error {
	query := `
		INSERT INTO team_standings (id, tournament_id, team_id, position, points,
			games, wins, wins_ot, wins_so, losses_so, losses_ot, losses, draws,
			goals_for, goals_against, goal_difference, group_name, birth_year,
			source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), NOW())
		ON CONFLICT (tournament_id, team_id, COALESCE(birth_year, 0), COALESCE(group_name, ''), source) DO UPDATE SET
			position = EXCLUDED.position,
			points = EXCLUDED.points,
			games = EXCLUDED.games,
			wins = EXCLUDED.wins,
			wins_ot = EXCLUDED.wins_ot,
			wins_so = EXCLUDED.wins_so,
			losses_so = EXCLUDED.losses_so,
			losses_ot = EXCLUDED.losses_ot,
			losses = EXCLUDED.losses,
			draws = EXCLUDED.draws,
			goals_for = EXCLUDED.goals_for,
			goals_against = EXCLUDED.goals_against,
			goal_difference = EXCLUDED.goal_difference,
			updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query,
		s.ID, s.TournamentID, s.TeamID, s.Position, s.Points,
		s.Games, s.Wins, s.WinsOT, s.WinsSO, s.LossesSO, s.LossesOT, s.Losses, s.Draws,
		s.GoalsFor, s.GoalsAgainst, s.GoalDifference, s.GroupName, s.BirthYear, s.Source,
	)
	if err != nil {
		return fmt.Errorf("upsert standing: %w", err)
	}
	return nil
}

func (r *StandingPostgres) UpsertBatch(ctx context.Context, standings []*entities.TeamStanding) error {
	for _, s := range standings {
		if err := r.Upsert(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (r *StandingPostgres) GetByTournament(ctx context.Context, tournamentID string) ([]*entities.TeamStanding, error) {
	var standings []*entities.TeamStanding
	err := r.db.SelectContext(ctx, &standings,
		`SELECT * FROM team_standings WHERE tournament_id = $1 ORDER BY position`, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("get standings by tournament: %w", err)
	}
	return standings, nil
}

func (r *StandingPostgres) GetByTeam(ctx context.Context, teamID string) ([]*entities.TeamStanding, error) {
	var standings []*entities.TeamStanding
	err := r.db.SelectContext(ctx, &standings,
		`SELECT * FROM team_standings WHERE team_id = $1 ORDER BY created_at DESC`, teamID)
	if err != nil {
		return nil, fmt.Errorf("get standings by team: %w", err)
	}
	return standings, nil
}

func (r *StandingPostgres) GetByTournamentAndGroup(ctx context.Context, tournamentID, groupName string) ([]*entities.TeamStanding, error) {
	var standings []*entities.TeamStanding
	err := r.db.SelectContext(ctx, &standings,
		`SELECT * FROM team_standings WHERE tournament_id = $1 AND group_name = $2 ORDER BY position`, tournamentID, groupName)
	if err != nil {
		return nil, fmt.Errorf("get standings by tournament and group: %w", err)
	}
	return standings, nil
}

func (r *StandingPostgres) DeleteByTournament(ctx context.Context, tournamentID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM team_standings WHERE tournament_id = $1`, tournamentID)
	if err != nil {
		return fmt.Errorf("delete standings by tournament: %w", err)
	}
	return nil
}
