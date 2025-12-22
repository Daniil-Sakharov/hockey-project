package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/domain/entities"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/scheduler/domain"
)

func sourceCondition(source string) string {
	switch source {
	case "junior":
		return "AND id NOT LIKE 'spb:%'"
	case "fhspb":
		return "AND id LIKE 'spb:%'"
	default:
		return ""
	}
}

// GetByPriorityForPlayers возвращает турниры для парсинга игроков по приоритету
func (r *TournamentPostgres) GetByPriorityForPlayers(ctx context.Context, priority domain.Priority, source string) ([]*entities.Tournament, error) {
	query := buildPriorityQuery("last_players_parsed_at", priority, source)

	var tournaments []*entities.Tournament
	if err := r.db.SelectContext(ctx, &tournaments, query); err != nil {
		return nil, fmt.Errorf("get tournaments by priority: %w", err)
	}
	return tournaments, nil
}

// GetByPriorityForStats возвращает турниры для парсинга статистики по приоритету
func (r *TournamentPostgres) GetByPriorityForStats(ctx context.Context, priority domain.Priority, source string) ([]*entities.Tournament, error) {
	query := buildPriorityQuery("last_stats_parsed_at", priority, source)

	var tournaments []*entities.Tournament
	if err := r.db.SelectContext(ctx, &tournaments, query); err != nil {
		return nil, fmt.Errorf("get tournaments by priority: %w", err)
	}
	return tournaments, nil
}

// UpdateLastPlayersParsed обновляет время последнего парсинга игроков
func (r *TournamentPostgres) UpdateLastPlayersParsed(ctx context.Context, tournamentID string) error {
	query := `UPDATE tournaments SET last_players_parsed_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tournamentID)
	return err
}

// UpdateLastStatsParsed обновляет время последнего парсинга статистики
func (r *TournamentPostgres) UpdateLastStatsParsed(ctx context.Context, tournamentID string) error {
	query := `UPDATE tournaments SET last_stats_parsed_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tournamentID)
	return err
}

func buildPriorityQuery(parsedAtField string, priority domain.Priority, source string) string {
	var condition string
	var interval string

	switch priority {
	case domain.PriorityActive:
		condition = "(is_ended = false OR end_date IS NULL)"
		interval = "4 hours"
	case domain.PriorityRecent:
		condition = "is_ended = true AND end_date > NOW() - INTERVAL '1 month'"
		interval = "1 day"
	case domain.PriorityMedium:
		condition = "is_ended = true AND end_date BETWEEN NOW() - INTERVAL '6 months' AND NOW() - INTERVAL '1 month'"
		interval = "14 days"
	case domain.PriorityOld:
		condition = "is_ended = true AND end_date BETWEEN NOW() - INTERVAL '1 year' AND NOW() - INTERVAL '6 months'"
		interval = "30 days"
	case domain.PriorityArchive:
		condition = "is_ended = true AND end_date < NOW() - INTERVAL '1 year'"
		interval = "180 days"
	}

	return fmt.Sprintf(`
		SELECT * FROM tournaments 
		WHERE %s 
		AND (%s IS NULL OR %s < NOW() - INTERVAL '%s')
		%s
		ORDER BY %s NULLS FIRST
	`, condition, parsedAtField, parsedAtField, interval, sourceCondition(source), parsedAtField)
}

// GetAllForParsing возвращает все турниры которые нужно парсить (для всех приоритетов)
func (r *TournamentPostgres) GetAllForParsing(ctx context.Context, forStats bool, source string) ([]*entities.Tournament, error) {
	parsedAtField := "last_players_parsed_at"
	if forStats {
		parsedAtField = "last_stats_parsed_at"
	}

	query := fmt.Sprintf(`
		SELECT * FROM tournaments 
		WHERE (
			-- ACTIVE: каждые 4 часа
			((is_ended = false OR end_date IS NULL) AND (%s IS NULL OR %s < NOW() - INTERVAL '4 hours'))
			OR
			-- RECENT: раз в день
			(is_ended = true AND end_date > NOW() - INTERVAL '1 month' AND (%s IS NULL OR %s < NOW() - INTERVAL '1 day'))
			OR
			-- MEDIUM: раз в 2 недели
			(is_ended = true AND end_date BETWEEN NOW() - INTERVAL '6 months' AND NOW() - INTERVAL '1 month' AND (%s IS NULL OR %s < NOW() - INTERVAL '14 days'))
			OR
			-- OLD: раз в месяц
			(is_ended = true AND end_date BETWEEN NOW() - INTERVAL '1 year' AND NOW() - INTERVAL '6 months' AND (%s IS NULL OR %s < NOW() - INTERVAL '30 days'))
			OR
			-- ARCHIVE: раз в полгода
			(is_ended = true AND end_date < NOW() - INTERVAL '1 year' AND (%s IS NULL OR %s < NOW() - INTERVAL '180 days'))
		)
		%s
		ORDER BY 
			CASE 
				WHEN is_ended = false OR end_date IS NULL THEN 1
				WHEN end_date > NOW() - INTERVAL '1 month' THEN 2
				WHEN end_date > NOW() - INTERVAL '6 months' THEN 3
				WHEN end_date > NOW() - INTERVAL '1 year' THEN 4
				ELSE 5
			END,
			%s NULLS FIRST
	`, parsedAtField, parsedAtField, parsedAtField, parsedAtField, parsedAtField, parsedAtField,
		parsedAtField, parsedAtField, parsedAtField, parsedAtField, sourceCondition(source), parsedAtField)

	var tournaments []*entities.Tournament
	if err := r.db.SelectContext(ctx, &tournaments, query); err != nil {
		return nil, fmt.Errorf("get all tournaments for parsing: %w", err)
	}
	return tournaments, nil
}

// CountBySource возвращает количество турниров по источнику
func (r *TournamentPostgres) CountBySource(ctx context.Context, source string) (int, error) {
	var query string
	switch source {
	case "junior":
		query = `SELECT COUNT(*) FROM tournaments WHERE id NOT LIKE 'spb:%'`
	default:
		query = `SELECT COUNT(*) FROM tournaments WHERE id LIKE 'spb:%'`
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query); err != nil {
		return 0, err
	}
	return count, nil
}

// HasData проверяет есть ли данные в БД (для bootstrap mode)
func (r *TournamentPostgres) HasData(ctx context.Context) (bool, error) {
	var count int
	if err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM tournaments LIMIT 1`); err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetNeverParsed возвращает турниры которые ни разу не парсились
func (r *TournamentPostgres) GetNeverParsed(ctx context.Context, forStats bool, source string) ([]*entities.Tournament, error) {
	parsedAtField := "last_players_parsed_at"
	if forStats {
		parsedAtField = "last_stats_parsed_at"
	}

	query := fmt.Sprintf(`SELECT * FROM tournaments WHERE %s IS NULL %s`, parsedAtField, sourceCondition(source))

	var tournaments []*entities.Tournament
	if err := r.db.SelectContext(ctx, &tournaments, query); err != nil {
		return nil, err
	}
	return tournaments, nil
}

// GetParsedBefore возвращает турниры спарсенные до указанного времени
func (r *TournamentPostgres) GetParsedBefore(ctx context.Context, before time.Time, forStats bool, source string) ([]*entities.Tournament, error) {
	parsedAtField := "last_players_parsed_at"
	if forStats {
		parsedAtField = "last_stats_parsed_at"
	}

	query := fmt.Sprintf(`SELECT * FROM tournaments WHERE %s < $1 %s ORDER BY %s`, parsedAtField, sourceCondition(source), parsedAtField)

	var tournaments []*entities.Tournament
	if err := r.db.SelectContext(ctx, &tournaments, query, before); err != nil {
		return nil, err
	}
	return tournaments, nil
}
