package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/telegram/application/services"
	"github.com/jmoiron/sqlx"
)

// PlayerSearchRepository реализация поиска игроков
type PlayerSearchRepository struct {
	db *sqlx.DB
}

// NewPlayerSearchRepository создает новый репозиторий
func NewPlayerSearchRepository(db *sqlx.DB) *PlayerSearchRepository {
	return &PlayerSearchRepository{db: db}
}

// SearchWithFilters выполняет поиск с фильтрами
// Группирует игроков по (name, birth_date) для дедупликации Junior/FHSPB
func (r *PlayerSearchRepository) SearchWithFilters(ctx context.Context, f services.SearchFilters) ([]*services.PlayerWithTeam, int, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if f.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("fp.name ILIKE $%d", argNum))
		args = append(args, "%"+f.FirstName+"%")
		argNum++
	}
	if f.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("fp.name ILIKE $%d", argNum))
		args = append(args, f.LastName+"%")
		argNum++
	}
	if f.BirthYear != nil {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM fp.birth_date) = $%d", argNum))
		args = append(args, *f.BirthYear)
		argNum++
	}
	if f.Position != "" {
		conditions = append(conditions, fmt.Sprintf("fp.position = $%d", argNum))
		args = append(args, f.Position)
		argNum++
	}
	if f.MinHeight != nil {
		conditions = append(conditions, fmt.Sprintf("fp.height >= $%d", argNum))
		args = append(args, *f.MinHeight)
		argNum++
	}
	if f.MaxHeight != nil {
		conditions = append(conditions, fmt.Sprintf("fp.height <= $%d", argNum))
		args = append(args, *f.MaxHeight)
		argNum++
	}
	if f.MinWeight != nil {
		conditions = append(conditions, fmt.Sprintf("fp.weight >= $%d", argNum))
		args = append(args, *f.MinWeight)
		argNum++
	}
	if f.MaxWeight != nil {
		conditions = append(conditions, fmt.Sprintf("fp.weight <= $%d", argNum))
		args = append(args, *f.MaxWeight)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// CTE: группируем игроков по (name, birth_date), берём самые свежие данные
	cte := `
		WITH 
		-- Самая свежая команда для каждого игрока
		latest_teams AS (
			SELECT DISTINCT ON (pt.player_id)
				pt.player_id,
				t.name as team_name,
				COALESCE(t.city, t.region, '') as team_city,
				tr.start_date as tournament_date
			FROM player_teams pt
			JOIN teams t ON pt.team_id = t.id
			LEFT JOIN tournaments tr ON pt.tournament_id = tr.id
			ORDER BY pt.player_id, tr.start_date DESC NULLS LAST
		),
		-- Объединённые игроки: группируем по name+birth_date, берём самые свежие данные
		fresh_players AS (
			SELECT DISTINCT ON (p.name, p.birth_date)
				p.id,
				p.name,
				p.birth_date,
				p.position,
				p.height,
				p.weight,
				lt.team_name,
				lt.team_city
			FROM players p
			LEFT JOIN latest_teams lt ON p.id = lt.player_id
			ORDER BY p.name, p.birth_date, lt.tournament_date DESC NULLS LAST
		)`

	// Count query - считаем уникальных игроков
	countQuery := fmt.Sprintf(`%s
		SELECT COUNT(*)
		FROM fresh_players fp
		%s`, cte, whereClause)

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`%s
		SELECT 
			fp.id,
			fp.name, 
			COALESCE(TO_CHAR(fp.birth_date, 'DD.MM.YYYY'), '') as birth_date,
			COALESCE(fp.position, '') as position, 
			COALESCE(fp.height, 0) as height, 
			COALESCE(fp.weight, 0) as weight,
			COALESCE(fp.team_name, '') as team_name, 
			COALESCE(fp.team_city, '') as team_city
		FROM fresh_players fp
		%s
		ORDER BY fp.name ASC
		LIMIT $%d OFFSET $%d
	`, cte, whereClause, argNum, argNum+1)

	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var players []*services.PlayerWithTeam
	for rows.Next() {
		var p services.PlayerWithTeam
		if err := rows.Scan(&p.ID, &p.Name, &p.BirthDate, &p.Position, &p.Height, &p.Weight, &p.TeamName, &p.TeamCity); err != nil {
			return nil, 0, err
		}
		players = append(players, &p)
	}

	return players, totalCount, nil
}
