package player_team

import (
	"context"
	"fmt"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
)

// Upsert создает или обновляет связь игрок-команда-турнир
func (r *repository) Upsert(ctx context.Context, pt *player_team.PlayerTeam) error {
	query := `
		INSERT INTO player_teams (
			player_id, team_id, tournament_id, season, 
			started_at, ended_at, is_active, 
			jersey_number, role, source, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, 
			$5, $6, $7, 
			$8, $9, $10, 
			$11, $12
		)
		ON CONFLICT (player_id, team_id, tournament_id) 
		DO UPDATE SET 
			season = EXCLUDED.season,
			started_at = EXCLUDED.started_at,
			ended_at = EXCLUDED.ended_at,
			is_active = EXCLUDED.is_active,
			jersey_number = EXCLUDED.jersey_number,
			role = EXCLUDED.role,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		pt.PlayerID,
		pt.TeamID,
		pt.TournamentID,
		pt.Season,
		pt.StartedAt,
		pt.EndedAt,
		pt.IsActive,
		pt.JerseyNumber,
		pt.Role,
		pt.Source,
		pt.CreatedAt,
		pt.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert player_team: %w", err)
	}

	return nil
}

// UpsertBatch создает или обновляет связи игрок-команда-турнир батчем
func (r *repository) UpsertBatch(ctx context.Context, pts []*player_team.PlayerTeam) error {
	if len(pts) == 0 {
		return nil
	}

	query := `
		INSERT INTO player_teams (
			player_id, team_id, tournament_id, season, 
			started_at, ended_at, is_active, 
			jersey_number, role, source, 
			created_at, updated_at
		) VALUES `

	// Формируем VALUES для батча
	values := make([]interface{}, 0, len(pts)*12)
	placeholders := ""

	for i, pt := range pts {
		if i > 0 {
			placeholders += ", "
		}
		offset := i * 12
		placeholders += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4,
			offset+5, offset+6, offset+7, offset+8,
			offset+9, offset+10, offset+11, offset+12,
		)

		values = append(values,
			pt.PlayerID,
			pt.TeamID,
			pt.TournamentID,
			pt.Season,
			pt.StartedAt,
			pt.EndedAt,
			pt.IsActive,
			pt.JerseyNumber,
			pt.Role,
			pt.Source,
			pt.CreatedAt,
			pt.UpdatedAt,
		)
	}

	query += placeholders + `
		ON CONFLICT (player_id, team_id, tournament_id) 
		DO UPDATE SET 
			season = EXCLUDED.season,
			started_at = EXCLUDED.started_at,
			ended_at = EXCLUDED.ended_at,
			is_active = EXCLUDED.is_active,
			jersey_number = EXCLUDED.jersey_number,
			role = EXCLUDED.role,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to batch upsert player_teams: %w", err)
	}

	return nil
}
