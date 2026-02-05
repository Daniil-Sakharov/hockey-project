package services

import (
	"context"
	"fmt"
	"sort"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/api/dto"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
)

var domainLabels = map[string]string{
	"https://pfo.fhr.ru":     "ПФО",
	"https://cfo.fhr.ru":     "ЦФО",
	"https://spb.fhr.ru":     "СПб",
	"https://junior.fhr.ru":  "Юниор",
	"https://dfo.fhr.ru":     "ДФО",
	"https://komi.fhr.ru":    "Коми",
	"https://kuzbass.fhr.ru": "Кузбасс",
	"https://len.fhr.ru":     "Ленобласть",
	"https://nsk.fhr.ru":     "Новосибирск",
	"https://sam.fhr.ru":     "Самара",
	"https://sfo.fhr.ru":     "СФО",
	"https://szfo.fhr.ru":    "СЗФО",
	"https://ufo.fhr.ru":     "УрФО",
	"https://vrn.fhr.ru":     "Воронеж",
	"https://yfo.fhr.ru":     "ЮФО",
}

// GetRankingsFilters returns available filter values for the current season.
func (s *ExploreMatchesService) GetRankingsFilters(ctx context.Context) (*dto.RankingsFiltersResponse, error) {
	var season string
	if err := s.db.GetContext(ctx, &season, "SELECT season FROM tournaments ORDER BY season DESC LIMIT 1"); err != nil {
		return nil, fmt.Errorf("failed to get current season: %w", err)
	}

	// Birth years
	var birthYears []int
	err := s.db.SelectContext(ctx, &birthYears, `
		SELECT DISTINCT ps.birth_year
		FROM player_statistics ps
		JOIN tournaments t ON ps.tournament_id = t.id
		WHERE t.season = $1 AND ps.group_name != 'Общая статистика' AND ps.birth_year > 0
		ORDER BY ps.birth_year DESC
	`, season)
	if err != nil {
		logger.Warn(ctx, "failed to get birth years for rankings filters: "+err.Error())
	}

	// Domains
	type domainRow struct {
		Domain string `db:"domain"`
	}
	var domainRows []domainRow
	err = s.db.SelectContext(ctx, &domainRows, `
		SELECT DISTINCT domain FROM tournaments WHERE season = $1 AND domain != '' ORDER BY domain
	`, season)
	if err != nil {
		logger.Warn(ctx, "failed to get domains for rankings filters: "+err.Error())
	}
	domains := make([]dto.DomainOption, 0, len(domainRows))
	for _, d := range domainRows {
		label := domainLabels[d.Domain]
		if label == "" {
			label = d.Domain
		}
		domains = append(domains, dto.DomainOption{Domain: d.Domain, Label: label})
	}
	sort.Slice(domains, func(i, j int) bool { return domains[i].Label < domains[j].Label })

	// Tournaments with their available birth years
	type tournamentRow struct {
		ID     string `db:"id"`
		Name   string `db:"name"`
		Domain string `db:"domain"`
	}
	var tournamentRows []tournamentRow
	err = s.db.SelectContext(ctx, &tournamentRows, `
		SELECT id, name, domain FROM tournaments WHERE season = $1 ORDER BY name
	`, season)
	if err != nil {
		logger.Warn(ctx, "failed to get tournaments for rankings filters: "+err.Error())
	}

	// Get birth years per tournament
	type tournamentBirthYear struct {
		TournamentID string `db:"tournament_id"`
		BirthYear    int    `db:"birth_year"`
	}
	var tby []tournamentBirthYear
	err = s.db.SelectContext(ctx, &tby, `
		SELECT DISTINCT ps.tournament_id, ps.birth_year
		FROM player_statistics ps
		JOIN tournaments t ON ps.tournament_id = t.id
		WHERE t.season = $1 AND ps.group_name != 'Общая статистика' AND ps.birth_year > 0
		ORDER BY ps.tournament_id, ps.birth_year DESC
	`, season)
	if err != nil {
		logger.Warn(ctx, "failed to get tournament birth years: "+err.Error())
	}
	tbyMap := make(map[string][]int)
	for _, v := range tby {
		tbyMap[v.TournamentID] = append(tbyMap[v.TournamentID], v.BirthYear)
	}

	tournaments := make([]dto.TournamentOption, len(tournamentRows))
	for i, t := range tournamentRows {
		tournaments[i] = dto.TournamentOption{
			ID: t.ID, Name: titleCase(t.Name), Domain: t.Domain,
			BirthYears: tbyMap[t.ID],
		}
	}

	// Groups per tournament
	type groupRow struct {
		TournamentID string `db:"tournament_id"`
		GroupName    string `db:"group_name"`
	}
	var groupRows []groupRow
	err = s.db.SelectContext(ctx, &groupRows, `
		SELECT DISTINCT ps.tournament_id, ps.group_name
		FROM player_statistics ps
		JOIN tournaments t ON ps.tournament_id = t.id
		WHERE t.season = $1 AND ps.group_name != 'Общая статистика' AND ps.group_name != ''
		ORDER BY ps.tournament_id, ps.group_name
	`, season)
	if err != nil {
		logger.Warn(ctx, "failed to get groups for rankings filters: "+err.Error())
	}
	groups := make([]dto.GroupOption, len(groupRows))
	for i, g := range groupRows {
		groups[i] = dto.GroupOption{Name: g.GroupName, TournamentID: g.TournamentID}
	}

	return &dto.RankingsFiltersResponse{
		BirthYears:  birthYears,
		Domains:     domains,
		Tournaments: tournaments,
		Groups:      groups,
	}, nil
}
