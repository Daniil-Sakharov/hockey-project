package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type TournamentStats struct {
	ID               string
	Name             string
	URL              string
	TotalPlayers     int
	PlayersWithStats int
	StatsRecords     int
	CoveragePercent  float64
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
	fmt.Println("🏆 АНАЛИЗ ПОКРЫТИЯ СТАТИСТИКИ ПО ТУРНИРАМ")
	fmt.Println("=" + strings.Repeat("=", 79))
	fmt.Println()

	// Получаем список турниров с анализом
	tournaments, err := analyzeTournaments(ctx, db)
	if err != nil {
		log.Fatalf("Ошибка анализа: %v", err)
	}

	// Выводим результаты
	printTournamentStats(tournaments)
}

func analyzeTournaments(ctx context.Context, db *sql.DB) ([]TournamentStats, error) {
	query := `
		SELECT 
			t.id,
			t.name,
			t.url,
			COUNT(DISTINCT pt.player_id) as total_players,
			COUNT(DISTINCT ps.player_id) as players_with_stats,
			COUNT(ps.id) as stats_records
		FROM tournaments t
		LEFT JOIN player_teams pt ON t.id = pt.tournament_id
		LEFT JOIN player_statistics ps ON t.id = ps.tournament_id
		GROUP BY t.id, t.name, t.url
		HAVING COUNT(DISTINCT pt.player_id) > 0
		ORDER BY t.name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournaments := []TournamentStats{}
	for rows.Next() {
		var ts TournamentStats
		var url sql.NullString

		err := rows.Scan(
			&ts.ID,
			&ts.Name,
			&url,
			&ts.TotalPlayers,
			&ts.PlayersWithStats,
			&ts.StatsRecords,
		)
		if err != nil {
			return nil, err
		}

		if url.Valid {
			ts.URL = url.String
		}

		if ts.TotalPlayers > 0 {
			ts.CoveragePercent = float64(ts.PlayersWithStats) * 100.0 / float64(ts.TotalPlayers)
		}

		tournaments = append(tournaments, ts)
	}

	return tournaments, nil
}

func printTournamentStats(tournaments []TournamentStats) {
	if len(tournaments) == 0 {
		fmt.Println("❌ Турниры не найдены в БД")
		return
	}

	// Группируем турниры по покрытию
	perfect := []TournamentStats{} // 100%
	good := []TournamentStats{}    // 80-99%
	medium := []TournamentStats{}  // 50-79%
	low := []TournamentStats{}     // 1-49%
	empty := []TournamentStats{}   // 0%

	for _, ts := range tournaments {
		switch {
		case ts.CoveragePercent == 100.0:
			perfect = append(perfect, ts)
		case ts.CoveragePercent >= 80.0:
			good = append(good, ts)
		case ts.CoveragePercent >= 50.0:
			medium = append(medium, ts)
		case ts.CoveragePercent > 0:
			low = append(low, ts)
		default:
			empty = append(empty, ts)
		}
	}

	// Общая статистика
	fmt.Printf("📊 Всего турниров: %d\n", len(tournaments))
	fmt.Println()

	// Отличное покрытие
	if len(perfect) > 0 {
		fmt.Printf("✅ ОТЛИЧНОЕ ПОКРЫТИЕ (100%%): %d турниров\n", len(perfect))
		fmt.Println(strings.Repeat("-", 80))
		for i, ts := range perfect {
			fmt.Printf("[%d] %s\n", i+1, ts.Name)
			fmt.Printf("    Игроков: %d | Записей статистики: %d\n", ts.TotalPlayers, ts.StatsRecords)
			if ts.URL != "" {
				fmt.Printf("    URL: %s/stats/\n", ts.URL)
			}
			fmt.Println()
		}
	}

	// Хорошее покрытие
	if len(good) > 0 {
		fmt.Printf("✅ ХОРОШЕЕ ПОКРЫТИЕ (80-99%%): %d турниров\n", len(good))
		fmt.Println(strings.Repeat("-", 80))
		for i, ts := range good {
			fmt.Printf("[%d] %s\n", i+1, ts.Name)
			fmt.Printf("    Игроков: %d | Со статистикой: %d (%.1f%%) | Записей: %d\n",
				ts.TotalPlayers, ts.PlayersWithStats, ts.CoveragePercent, ts.StatsRecords)
			if ts.URL != "" {
				fmt.Printf("    URL: %s/stats/\n", ts.URL)
			}
			fmt.Println()
		}
	}

	// Среднее покрытие
	if len(medium) > 0 {
		fmt.Printf("⚠️  СРЕДНЕЕ ПОКРЫТИЕ (50-79%%): %d турниров\n", len(medium))
		fmt.Println(strings.Repeat("-", 80))
		for i, ts := range medium {
			fmt.Printf("[%d] %s\n", i+1, ts.Name)
			fmt.Printf("    Игроков: %d | Со статистикой: %d (%.1f%%) | Записей: %d\n",
				ts.TotalPlayers, ts.PlayersWithStats, ts.CoveragePercent, ts.StatsRecords)
			if ts.URL != "" {
				fmt.Printf("    URL: %s/stats/\n", ts.URL)
			}
			fmt.Println()
		}
	}

	// Низкое покрытие
	if len(low) > 0 {
		fmt.Printf("🔴 НИЗКОЕ ПОКРЫТИЕ (1-49%%): %d турниров\n", len(low))
		fmt.Println(strings.Repeat("-", 80))
		for i, ts := range low {
			fmt.Printf("[%d] %s\n", i+1, ts.Name)
			fmt.Printf("    Игроков: %d | Со статистикой: %d (%.1f%%) | Записей: %d\n",
				ts.TotalPlayers, ts.PlayersWithStats, ts.CoveragePercent, ts.StatsRecords)
			if ts.URL != "" {
				fmt.Printf("    URL: %s/stats/\n", ts.URL)
			}
			fmt.Println()
		}
	}

	// Без статистики
	if len(empty) > 0 {
		fmt.Printf("❌ БЕЗ СТАТИСТИКИ (0%%): %d турниров\n", len(empty))
		fmt.Println(strings.Repeat("-", 80))
		for i, ts := range empty {
			fmt.Printf("[%d] %s\n", i+1, ts.Name)
			fmt.Printf("    Игроков: %d | Записей статистики: 0\n", ts.TotalPlayers)
			if ts.URL != "" {
				fmt.Printf("    URL: %s/stats/\n", ts.URL)
			}
			fmt.Println()
		}
	}

	// Итоги
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📈 ИТОГОВАЯ СТАТИСТИКА:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("  ✅ Отличное покрытие (100%%):    %d турниров\n", len(perfect))
	fmt.Printf("  ✅ Хорошее покрытие (80-99%%):   %d турниров\n", len(good))
	fmt.Printf("  ⚠️  Среднее покрытие (50-79%%):  %d турниров\n", len(medium))
	fmt.Printf("  🔴 Низкое покрытие (1-49%%):     %d турниров\n", len(low))
	fmt.Printf("  ❌ Без статистики (0%%):         %d турниров\n", len(empty))
	fmt.Println()

	if len(empty) > 0 || len(low) > 0 {
		fmt.Println("💡 Рекомендуется запустить парсер статистики:")
		fmt.Println("   task stats:parse")
	} else {
		fmt.Println("✅ Все турниры имеют хорошее покрытие статистики!")
	}
	fmt.Println(strings.Repeat("=", 80))
}
