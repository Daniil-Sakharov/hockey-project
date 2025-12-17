package profile

import (
	"context"
	"strings"

	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	domainPlayerStats "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
)

// getRecentTournaments получает последние турниры с обработкой группы "Общая"
func (s *Service) getRecentTournaments(ctx context.Context, playerID string, limit int) ([]domainPlayer.TournamentStats, error) {
	tournaments, err := s.statsRepo.GetRecentTournaments(ctx, playerID, limit)
	if err != nil {
		return nil, err
	}

	// Группируем турниры по tournament_id для определения уникальных групп
	tournamentGroups := buildTournamentGroupsMap(tournaments)

	// Конвертируем в доменную модель с обработкой группы "Общая"
	result := make([]domainPlayer.TournamentStats, 0, len(tournaments))
	for _, t := range tournaments {
		groupName := processGroupName(t.GroupName, tournamentGroups[t.TournamentID])

		// Пропускаем "Общую" если есть другие группы
		if shouldSkipTournament(t.GroupName, tournamentGroups[t.TournamentID]) {
			continue
		}

		result = append(result, domainPlayer.TournamentStats{
			TournamentName: t.TournamentName,
			GroupName:      groupName,
			Season:         t.Season,
			IsChampionship: t.IsChampionship,
			Games:          t.Games,
			Goals:          t.Goals,
			Assists:        t.Assists,
			Points:         t.Points,
			PlusMinus:      t.PlusMinus,
			Penalties:      t.PenaltyMinutes,
			HatTricks:      t.HatTricks,
			WinningGoals:   t.GameWinningGoals,
		})
	}

	return result, nil
}

// buildTournamentGroupsMap группирует турниры по tournament_id для определения уникальных групп
func buildTournamentGroupsMap(tournaments []*domainPlayerStats.TournamentStat) map[string][]string {
	tournamentGroups := make(map[string][]string)
	for _, t := range tournaments {
		tournamentGroups[t.TournamentID] = append(tournamentGroups[t.TournamentID], t.GroupName)
	}
	return tournamentGroups
}

// processGroupName обрабатывает название группы согласно правилам:
// - Если есть другие группы кроме "Общая" - показываем только другие
// - Если только "Общая" - не показываем название группы вообще
func processGroupName(groupName string, allGroups []string) string {
	hasOtherGroups := hasNonGeneralGroups(allGroups)

	// Если только "Общая" - не показываем название группы
	if !hasOtherGroups && isGeneralGroup(groupName) {
		return ""
	}

	return groupName
}

// shouldSkipTournament проверяет нужно ли пропустить турнир
// Возвращает true если это "Общая" группа и есть другие группы
func shouldSkipTournament(groupName string, allGroups []string) bool {
	hasOtherGroups := hasNonGeneralGroups(allGroups)
	return hasOtherGroups && isGeneralGroup(groupName)
}

// hasNonGeneralGroups проверяет есть ли группы кроме "Общая"
func hasNonGeneralGroups(groups []string) bool {
	for _, g := range groups {
		if !isGeneralGroup(g) {
			return true
		}
	}
	return false
}

// isGeneralGroup проверяет является ли группа "Общей"
func isGeneralGroup(groupName string) bool {
	return strings.Contains(strings.ToLower(groupName), "общ")
}
