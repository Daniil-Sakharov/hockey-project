package report

import (
	"context"
	"fmt"
	"sort"

	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	domainPlayerStats "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	domainPlayerTeam "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	domainTeam "github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	domainTournament "github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
)

// DataCollector собирает все данные для отчета
type DataCollector struct {
	playerRepo     domainPlayer.Repository
	statsRepo      domainPlayerStats.Repository
	playerTeamRepo domainPlayerTeam.Repository
	teamRepo       domainTeam.Repository
	tournamentRepo domainTournament.Repository
}

// NewDataCollector создает новый DataCollector
func NewDataCollector(
	playerRepo domainPlayer.Repository,
	statsRepo domainPlayerStats.Repository,
	playerTeamRepo domainPlayerTeam.Repository,
	teamRepo domainTeam.Repository,
	tournamentRepo domainTournament.Repository,
) *DataCollector {
	return &DataCollector{
		playerRepo:     playerRepo,
		statsRepo:      statsRepo,
		playerTeamRepo: playerTeamRepo,
		teamRepo:       teamRepo,
		tournamentRepo: tournamentRepo,
	}
}

// CollectFullReport собирает все данные для полного отчета игрока
func (dc *DataCollector) CollectFullReport(ctx context.Context, playerID string) (*FullPlayerReport, error) {
	// 1. Получаем базовую информацию об игроке
	player, err := dc.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	// 2. Получаем команду и регион
	team, region := dc.getPlayerTeamAndRegion(ctx, playerID)

	// 3. Получаем ВСЮ статистику игрока
	allStats, err := dc.statsRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player stats: %w", err)
	}

	// 4. Формируем отчет
	report := &FullPlayerReport{
		Player: PlayerInfo{
			ID:        player.ID,
			Name:      player.Name,
			BirthYear: player.BirthDate.Year(),
			Position:  player.Position,
			Height:    player.Height,
			Weight:    player.Weight,
			Team:      team,
			Region:    region,
		},
		HasStats: len(allStats) > 0,
	}

	if len(allStats) == 0 {
		return report, nil
	}

	// 5. Агрегируем статистику
	report.TotalStats = dc.aggregateTotalStats(allStats)
	report.GoalsByType = dc.calculateGoalsBreakdown(report.TotalStats)
	report.GoalsByPeriod = PeriodGoals{
		Period1:  report.TotalStats.GoalsPeriod1,
		Period2:  report.TotalStats.GoalsPeriod2,
		Period3:  report.TotalStats.GoalsPeriod3,
		Overtime: report.TotalStats.GoalsOvertime,
	}

	// 6. Статистика по сезонам
	report.SeasonStats = dc.aggregateBySeasons(ctx, allStats)
	report.HasMultipleSeasons = len(report.SeasonStats) > 1

	// 7. Полная история турниров
	report.Tournaments = dc.buildTournamentsList(ctx, allStats)

	// 8. Флаг детальной статистики
	report.HasDetailedStats = report.TotalStats.GoalsEvenStrength > 0 ||
		report.TotalStats.GoalsPowerPlay > 0 ||
		report.TotalStats.GoalsShortHanded > 0

	return report, nil
}

// getPlayerTeamAndRegion получает текущую команду и регион игрока
func (dc *DataCollector) getPlayerTeamAndRegion(ctx context.Context, playerID string) (team, region string) {
	team = "Неизвестно"
	region = "Неизвестно"

	// Получаем связи игрок-команда
	links, err := dc.playerTeamRepo.GetByPlayer(ctx, playerID)
	if err != nil || len(links) == 0 {
		return
	}

	// Ищем активную команду или последнюю
	var activeLink *domainPlayerTeam.PlayerTeam
	for _, link := range links {
		if link.IsActive {
			activeLink = link
			break
		}
	}
	if activeLink == nil && len(links) > 0 {
		activeLink = links[0]
	}

	if activeLink == nil {
		return
	}

	// Получаем информацию о команде
	teamInfo, err := dc.teamRepo.GetByID(ctx, activeLink.TeamID)
	if err != nil {
		return
	}

	team = teamInfo.Name
	if teamInfo.City != "" {
		region = teamInfo.City
	}

	return
}

// aggregateTotalStats агрегирует статистику за всё время
func (dc *DataCollector) aggregateTotalStats(stats []*domainPlayerStats.PlayerStatistic) TotalStatistics {
	total := TotalStatistics{}
	tournamentsSet := make(map[string]bool)

	for _, s := range stats {
		tournamentsSet[s.TournamentID] = true
		total.TotalGames += s.Games
		total.TotalGoals += s.Goals
		total.TotalAssists += s.Assists
		total.TotalPoints += s.Points
		total.TotalPlusMinus += s.PlusMinus
		total.TotalPenalties += s.PenaltyMinutes

		total.GoalsEvenStrength += s.GoalsEvenStrength
		total.GoalsPowerPlay += s.GoalsPowerPlay
		total.GoalsShortHanded += s.GoalsShortHanded
		total.GoalsPeriod1 += s.GoalsPeriod1
		total.GoalsPeriod2 += s.GoalsPeriod2
		total.GoalsPeriod3 += s.GoalsPeriod3
		total.GoalsOvertime += s.GoalsOvertime

		total.TotalHatTricks += s.HatTricks
		total.TotalWinningGoals += s.GameWinningGoals
	}

	total.TotalTournaments = len(tournamentsSet)

	// Средние показатели
	if total.TotalGames > 0 {
		total.GoalsPerGame = float64(total.TotalGoals) / float64(total.TotalGames)
		total.AssistsPerGame = float64(total.TotalAssists) / float64(total.TotalGames)
		total.PointsPerGame = float64(total.TotalPoints) / float64(total.TotalGames)
		total.PenaltiesPerGame = float64(total.TotalPenalties) / float64(total.TotalGames)
	}

	return total
}

// calculateGoalsBreakdown рассчитывает распределение голов по типу
func (dc *DataCollector) calculateGoalsBreakdown(total TotalStatistics) GoalsBreakdown {
	breakdown := GoalsBreakdown{
		EvenStrength: total.GoalsEvenStrength,
		PowerPlay:    total.GoalsPowerPlay,
		ShortHanded:  total.GoalsShortHanded,
	}

	totalTypedGoals := breakdown.EvenStrength + breakdown.PowerPlay + breakdown.ShortHanded
	if totalTypedGoals > 0 {
		breakdown.EvenStrengthPct = float64(breakdown.EvenStrength) / float64(totalTypedGoals) * 100
		breakdown.PowerPlayPct = float64(breakdown.PowerPlay) / float64(totalTypedGoals) * 100
		breakdown.ShortHandedPct = float64(breakdown.ShortHanded) / float64(totalTypedGoals) * 100
	}

	return breakdown
}

// aggregateBySeasons группирует статистику по сезонам
func (dc *DataCollector) aggregateBySeasons(ctx context.Context, stats []*domainPlayerStats.PlayerStatistic) []SeasonSummary {
	seasonMap := make(map[string]*SeasonSummary)

	for _, s := range stats {
		// Получаем сезон турнира
		tournament, err := dc.tournamentRepo.GetByID(ctx, s.TournamentID)
		if err != nil {
			continue
		}

		season := tournament.Season
		if season == "" {
			continue
		}

		if _, exists := seasonMap[season]; !exists {
			seasonMap[season] = &SeasonSummary{Season: season}
		}

		seasonMap[season].Games += s.Games
		seasonMap[season].Goals += s.Goals
		seasonMap[season].Assists += s.Assists
		seasonMap[season].Points += s.Points
	}

	// Конвертируем в slice и сортируем
	result := make([]SeasonSummary, 0, len(seasonMap))
	for _, ss := range seasonMap {
		result = append(result, *ss)
	}

	// Сортировка по сезону (от старого к новому для графика)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Season < result[j].Season
	})

	return result
}

// buildTournamentsList строит полный список турниров
func (dc *DataCollector) buildTournamentsList(ctx context.Context, stats []*domainPlayerStats.PlayerStatistic) []TournamentFullStats {
	// Группируем по турнирам (может быть несколько записей для одного турнира в разных группах)
	tournamentMap := make(map[string]*TournamentFullStats)

	for _, s := range stats {
		key := s.TournamentID + "_" + s.GroupName

		if _, exists := tournamentMap[key]; !exists {
			// Получаем информацию о турнире
			tournament, err := dc.tournamentRepo.GetByID(ctx, s.TournamentID)
			tournamentName := "Неизвестный турнир"
			season := ""
			if err == nil && tournament != nil {
				tournamentName = tournament.Name
				season = tournament.Season
			}

			// Получаем название команды
			teamName := "Неизвестно"
			if team, err := dc.teamRepo.GetByID(ctx, s.TeamID); err == nil && team != nil {
				teamName = team.Name
			}

			tournamentMap[key] = &TournamentFullStats{
				TournamentID:   s.TournamentID,
				TournamentName: tournamentName,
				GroupName:      s.GroupName,
				Season:         season,
				TeamName:       teamName,
			}
		}

		t := tournamentMap[key]
		t.Games += s.Games
		t.Goals += s.Goals
		t.Assists += s.Assists
		t.Points += s.Points
		t.Plus += s.Plus
		t.Minus += s.Minus
		t.PlusMinus += s.PlusMinus
		t.PenaltyMinutes += s.PenaltyMinutes
		t.GoalsEvenStrength += s.GoalsEvenStrength
		t.GoalsPowerPlay += s.GoalsPowerPlay
		t.GoalsShortHanded += s.GoalsShortHanded
		t.GoalsPeriod1 += s.GoalsPeriod1
		t.GoalsPeriod2 += s.GoalsPeriod2
		t.GoalsPeriod3 += s.GoalsPeriod3
		t.GoalsOvertime += s.GoalsOvertime
		t.HatTricks += s.HatTricks
		t.GameWinningGoals += s.GameWinningGoals
	}

	// Конвертируем в slice
	result := make([]TournamentFullStats, 0, len(tournamentMap))
	for _, t := range tournamentMap {
		// Рассчитываем средние
		if t.Games > 0 {
			t.GoalsPerGame = float64(t.Goals) / float64(t.Games)
			t.PointsPerGame = float64(t.Points) / float64(t.Games)
			t.PenaltiesPerGame = float64(t.PenaltyMinutes) / float64(t.Games)
		}
		result = append(result, *t)
	}

	// Сортируем по сезону (от нового к старому)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Season != result[j].Season {
			return result[i].Season > result[j].Season
		}
		return result[i].TournamentName < result[j].TournamentName
	})

	return result
}

// PrepareChartData подготавливает данные для JavaScript графиков
func (dc *DataCollector) PrepareChartData(report *FullPlayerReport) ChartData {
	data := ChartData{}

	// Круговая диаграмма голов по типу
	if report.HasDetailedStats {
		data.GoalTypeLabels = []string{"В равных", "В большинстве", "В меньшинстве"}
		data.GoalTypeValues = []int{
			report.GoalsByType.EvenStrength,
			report.GoalsByType.PowerPlay,
			report.GoalsByType.ShortHanded,
		}
		data.GoalTypeColors = []string{"#4a90d9", "#7bb8e8", "#2c5aa0"}
	}

	// Столбчатая диаграмма голов по периодам
	data.PeriodLabels = []string{"1-й период", "2-й период", "3-й период", "Овертайм"}
	data.PeriodValues = []int{
		report.GoalsByPeriod.Period1,
		report.GoalsByPeriod.Period2,
		report.GoalsByPeriod.Period3,
		report.GoalsByPeriod.Overtime,
	}

	// Линейный график прогресса по сезонам
	for _, s := range report.SeasonStats {
		data.SeasonLabels = append(data.SeasonLabels, s.Season)
		data.SeasonGoals = append(data.SeasonGoals, s.Goals)
		data.SeasonPoints = append(data.SeasonPoints, s.Points)
	}

	// Radar chart - нормализуем значения для наглядности
	data.RadarLabels = []string{"Голы", "Пасы", "+/-", "Хет-трики", "Поб. голы"}

	// Нормализация для radar chart (приводим к шкале 0-100)
	maxGoals := float64(report.TotalStats.TotalGoals)
	maxAssists := float64(report.TotalStats.TotalAssists)
	maxPlusMinus := float64(abs(report.TotalStats.TotalPlusMinus))
	maxHatTricks := float64(report.TotalStats.TotalHatTricks)
	maxWinningGoals := float64(report.TotalStats.TotalWinningGoals)

	// Простая нормализация - показываем реальные значения
	data.RadarValues = []float64{
		maxGoals,
		maxAssists,
		maxPlusMinus,
		maxHatTricks * 10, // Умножаем для наглядности
		maxWinningGoals,
	}

	return data
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
