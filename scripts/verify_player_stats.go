package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type PlayerInfo struct {
	ID          string
	Name        string
	BirthYear   int
	HasStats    bool
	Tournaments []TournamentInfo
}

type TournamentInfo struct {
	ID   string
	Name string
	URL  string
}

type StatsReport struct {
	TotalPlayers        int
	PlayersWithStats    int
	PlayersWithoutStats int
	CoveragePercent     float64
	PlayersByYear       map[int]*YearStats
}

type YearStats struct {
	Year                int
	Total               int
	WithStats           int
	WithoutStats        int
	WithoutStatsPlayers []PlayerInfo
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println("🔍 ПРОВЕРКА ПОКРЫТИЯ СТАТИСТИКИ ИГРОКОВ")
	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println()

	// 1. Получаем всех игроков
	players, err := getAllPlayers(ctx, db)
	if err != nil {
		log.Fatalf("Ошибка получения игроков: %v", err)
	}

	// 2. Проверяем наличие статистики
	err = checkPlayersStats(ctx, db, players)
	if err != nil {
		log.Fatalf("Ошибка проверки статистики: %v", err)
	}

	// 3. Загружаем турниры для игроков без статистики
	err = loadPlayerTournaments(ctx, db, players)
	if err != nil {
		log.Fatalf("Ошибка загрузки турниров: %v", err)
	}

	// 4. Формируем отчет
	report := generateReport(players)

	// 5. Выводим результаты
	printReport(report)
}

func getAllPlayers(ctx context.Context, db *sql.DB) ([]PlayerInfo, error) {
	query := `
		SELECT id, name, EXTRACT(YEAR FROM birth_date)::INT as birth_year
		FROM players
		ORDER BY birth_date DESC, name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []PlayerInfo{}
	for rows.Next() {
		var p PlayerInfo
		var birthYear sql.NullInt64

		if err := rows.Scan(&p.ID, &p.Name, &birthYear); err != nil {
			return nil, err
		}

		if birthYear.Valid {
			p.BirthYear = int(birthYear.Int64)
		}

		players = append(players, p)
	}

	return players, nil
}

func checkPlayersStats(ctx context.Context, db *sql.DB, players []PlayerInfo) error {
	// Получаем список игроков с статистикой
	query := `
		SELECT DISTINCT player_id
		FROM player_statistics
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	playersWithStats := make(map[string]bool)
	for rows.Next() {
		var playerID string
		if err := rows.Scan(&playerID); err != nil {
			return err
		}
		playersWithStats[playerID] = true
	}

	// Помечаем игроков
	for i := range players {
		players[i].HasStats = playersWithStats[players[i].ID]
	}

	return nil
}

func loadPlayerTournaments(ctx context.Context, db *sql.DB, players []PlayerInfo) error {
	// Для игроков без статистики загружаем их турниры
	for i := range players {
		if players[i].HasStats {
			continue
		}

		query := `
			SELECT DISTINCT t.id, t.name, t.url
			FROM player_teams pt
			JOIN tournaments t ON pt.tournament_id = t.id
			WHERE pt.player_id = $1
			ORDER BY t.name
		`

		rows, err := db.QueryContext(ctx, query, players[i].ID)
		if err != nil {
			return err
		}

		tournaments := []TournamentInfo{}
		for rows.Next() {
			var t TournamentInfo
			var url sql.NullString

			if err := rows.Scan(&t.ID, &t.Name, &url); err != nil {
				rows.Close()
				return err
			}

			if url.Valid {
				t.URL = url.String
			}

			tournaments = append(tournaments, t)
		}
		rows.Close()

		players[i].Tournaments = tournaments
	}

	return nil
}

func generateReport(players []PlayerInfo) StatsReport {
	report := StatsReport{
		TotalPlayers:  len(players),
		PlayersByYear: make(map[int]*YearStats),
	}

	for _, p := range players {
		if p.HasStats {
			report.PlayersWithStats++
		} else {
			report.PlayersWithoutStats++
		}

		// Группируем по годам
		year := p.BirthYear
		if year == 0 {
			year = 9999 // для игроков без указанного года
		}

		if report.PlayersByYear[year] == nil {
			report.PlayersByYear[year] = &YearStats{
				Year:                year,
				WithoutStatsPlayers: []PlayerInfo{},
			}
		}

		ys := report.PlayersByYear[year]
		ys.Total++

		if p.HasStats {
			ys.WithStats++
		} else {
			ys.WithoutStats++
			ys.WithoutStatsPlayers = append(ys.WithoutStatsPlayers, p)
		}
	}

	if report.TotalPlayers > 0 {
		report.CoveragePercent = float64(report.PlayersWithStats) * 100.0 / float64(report.TotalPlayers)
	}

	return report
}

func printReport(report StatsReport) {
	// Общая статистика
	fmt.Println("📊 ОБЩАЯ СТАТИСТИКА:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  Всего игроков:           %d\n", report.TotalPlayers)
	fmt.Printf("  Со статистикой:          %d (%.2f%%)\n", report.PlayersWithStats, report.CoveragePercent)
	fmt.Printf("  БЕЗ статистики:          %d (%.2f%%)\n", report.PlayersWithoutStats, 100.0-report.CoveragePercent)
	fmt.Println()

	// Сортируем года
	years := make([]int, 0, len(report.PlayersByYear))
	for year := range report.PlayersByYear {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	// Показываем только проблемные года (где есть игроки без статистики)
	fmt.Println("📅 ИГРОКИ БЕЗ СТАТИСТИКИ ПО ГОДАМ РОЖДЕНИЯ:")
	fmt.Println(strings.Repeat("=", 80))

	hasProblems := false
	for _, year := range years {
		ys := report.PlayersByYear[year]

		if ys.WithoutStats == 0 {
			continue // Пропускаем года где все ОК
		}

		hasProblems = true

		yearLabel := fmt.Sprintf("%d", year)
		if year == 9999 {
			yearLabel = "Не указан"
		}

		fmt.Printf("\n🔴 ГОД РОЖДЕНИЯ: %s\n", yearLabel)
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("  Всего игроков:    %d\n", ys.Total)
		fmt.Printf("  Со статистикой:   %d (%.2f%%)\n", ys.WithStats, float64(ys.WithStats)*100.0/float64(ys.Total))
		fmt.Printf("  БЕЗ статистики:   %d (%.2f%%)\n", ys.WithoutStats, float64(ys.WithoutStats)*100.0/float64(ys.Total))
		fmt.Println()

		// Показываем первых 10 игроков без статистики
		limit := 10
		if len(ys.WithoutStatsPlayers) < limit {
			limit = len(ys.WithoutStatsPlayers)
		}

		if limit > 0 {
			fmt.Printf("  📋 Примеры игроков БЕЗ статистики (показано %d из %d):\n", limit, len(ys.WithoutStatsPlayers))
			fmt.Println(strings.Repeat("-", 80))

			for i := 0; i < limit; i++ {
				p := ys.WithoutStatsPlayers[i]
				fmt.Printf("  [%d] %s (ID: %s)\n", i+1, p.Name, p.ID)

				if len(p.Tournaments) > 0 {
					fmt.Printf("      Турниры: %d\n", len(p.Tournaments))
					for j, t := range p.Tournaments {
						if j >= 3 {
							fmt.Printf("      ... и еще %d турниров\n", len(p.Tournaments)-3)
							break
						}
						fmt.Printf("        - %s\n", t.Name)
						if t.URL != "" {
							fmt.Printf("          %s\n", t.URL)
						}
					}
				} else {
					fmt.Printf("      Турниры: не найдены\n")
				}
				fmt.Println()
			}

			if len(ys.WithoutStatsPlayers) > limit {
				fmt.Printf("  ... и еще %d игроков\n", len(ys.WithoutStatsPlayers)-limit)
				fmt.Println()
			}
		}
	}

	if !hasProblems {
		fmt.Println("✅ У всех игроков есть статистика!")
		fmt.Println()
	}

	// Итоговые рекомендации
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("💡 РЕКОМЕНДАЦИИ:")
	fmt.Println(strings.Repeat("-", 80))

	if report.CoveragePercent >= 90.0 {
		fmt.Println("  ✅ Отличное покрытие! Более 90% игроков имеют статистику.")
	} else if report.CoveragePercent >= 50.0 {
		fmt.Println("  ⚠️  Среднее покрытие. Рекомендуется запустить парсер статистики.")
	} else {
		fmt.Println("  🔴 Низкое покрытие! Необходимо срочно запустить парсер статистики.")
	}

	fmt.Println()
	fmt.Println("  Для парсинга статистики запустите:")
	fmt.Println("    task stats:parse")
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
}
