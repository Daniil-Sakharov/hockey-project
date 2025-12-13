package player

import (
	"context"
	"fmt"
	"strings"

	"github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
)

// Search ищет игроков по фильтрам (старый метод, используется SearchWithTeam)
func (r *repository) Search(ctx context.Context, filters player.SearchFilters) ([]*player.Player, error) {
	results, _, err := r.SearchWithTeam(ctx, filters)
	if err != nil {
		return nil, err
	}

	players := make([]*player.Player, len(results))
	for i, result := range results {
		players[i] = result.Player
	}

	return players, nil
}

// SearchWithTeam ищет игроков с информацией о команде
func (r *repository) SearchWithTeam(ctx context.Context, filters player.SearchFilters) ([]*player.PlayerWithTeam, int, error) {
	// Базовый запрос с JOIN на последнюю активную команду
	baseQuery := `
		WITH latest_teams AS (
			SELECT DISTINCT ON (pt.player_id)
				pt.player_id,
				t.name as team_name,
				t.city as team_city
			FROM player_teams pt
			JOIN teams t ON pt.team_id = t.id
			WHERE pt.is_active = true
			ORDER BY pt.player_id, pt.started_at DESC NULLS LAST
		)
		SELECT 
			p.id, p.profile_url, p.name, p.birth_date, p.position,
			p.height, p.weight, p.handedness,
			p.registry_id, p.school, p.rank, p.data_season,
			p.source, p.created_at, p.updated_at,
			COALESCE(lt.team_name, '') as team_name,
			COALESCE(lt.team_city, '') as team_city
		FROM players p
		LEFT JOIN latest_teams lt ON p.id = lt.player_id
	`

	// Строим WHERE условия
	conditions := []string{}
	args := []interface{}{}
	argCounter := 1

	// Фильтр по имени (ILIKE для поиска по части строки)
	if filters.FirstName != "" {
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE $%d", argCounter))
		args = append(args, "%"+filters.FirstName+"%")
		argCounter++
	}

	// Фильтр по фамилии
	if filters.LastName != "" {
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE $%d", argCounter))
		args = append(args, filters.LastName+"%") // Фамилия обычно в начале
		argCounter++
	}

	// Фильтр по году рождения
	if filters.BirthYear != nil {
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM p.birth_date) = $%d", argCounter))
		args = append(args, *filters.BirthYear)
		argCounter++
	}

	// Фильтр по позиции
	if filters.Position != "" {
		conditions = append(conditions, fmt.Sprintf("p.position = $%d", argCounter))
		args = append(args, filters.Position)
		argCounter++
	}

	// Фильтр по росту
	if filters.MinHeight != nil && filters.MaxHeight != nil {
		conditions = append(conditions, fmt.Sprintf("p.height BETWEEN $%d AND $%d", argCounter, argCounter+1))
		args = append(args, *filters.MinHeight, *filters.MaxHeight)
		argCounter += 2
	}

	// Фильтр по весу
	if filters.MinWeight != nil && filters.MaxWeight != nil {
		conditions = append(conditions, fmt.Sprintf("p.weight BETWEEN $%d AND $%d", argCounter, argCounter+1))
		args = append(args, *filters.MinWeight, *filters.MaxWeight)
		argCounter += 2
	}

	// Добавляем WHERE если есть условия
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// Запрос для подсчета общего количества
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM players p
		%s
	`, whereClause)

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count players: %w", err)
	}

	// Основной запрос с сортировкой и пагинацией
	query := fmt.Sprintf(`
		%s
		%s
		ORDER BY p.name ASC
		LIMIT $%d OFFSET $%d
	`, baseQuery, whereClause, argCounter, argCounter+1)

	args = append(args, filters.Limit, filters.Offset)

	// Выполняем запрос
	type Row struct {
		player.Player
		TeamName string `db:"team_name"`
		TeamCity string `db:"team_city"`
	}

	var rows []Row
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to search players: %w", err)
	}

	// Конвертируем результаты
	results := make([]*player.PlayerWithTeam, len(rows))

	for i, row := range rows {
		p := row.Player
		results[i] = &player.PlayerWithTeam{
			Player:   &p,
			TeamName: row.TeamName,
			TeamCity: row.TeamCity,
		}
	}

	return results, totalCount, nil
}
