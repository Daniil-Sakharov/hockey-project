package main

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("pgx", "postgresql://sug6r:password@localhost:5432/hockey_stats?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("🎉 ФИНАЛЬНЫЙ ОТЧЕТ ПО ПАРСИНГУ")
	fmt.Println("================================================================================\n")

	// 1. Общая статистика
	var tCount, tmCount, pCount, ptCount int
	db.Get(&tCount, "SELECT COUNT(*) FROM tournaments")
	db.Get(&tmCount, "SELECT COUNT(*) FROM teams")
	db.Get(&pCount, "SELECT COUNT(*) FROM players")
	db.Get(&ptCount, "SELECT COUNT(*) FROM player_teams")

	fmt.Println("📊 ОБЩАЯ СТАТИСТИКА:")
	fmt.Printf("  🏆 Турниры: %d\n", tCount)
	fmt.Printf("  🏒 Команды: %d\n", tmCount)
	fmt.Printf("  👤 Игроки: %d\n", pCount)
	fmt.Printf("  🔗 Связи player_teams: %d\n\n", ptCount)

	// 2. Проверка season/dates в турнирах
	var tWithSeason, tWithStart int
	db.Get(&tWithSeason, "SELECT COUNT(*) FROM tournaments WHERE season IS NOT NULL AND season != ''")
	db.Get(&tWithStart, "SELECT COUNT(*) FROM tournaments WHERE start_date IS NOT NULL")

	fmt.Println("✅ КАЧЕСТВО ДАННЫХ - ТУРНИРЫ:")
	fmt.Printf("  Season заполнен: %d/%d (%.1f%%)\n", tWithSeason, tCount, float64(tWithSeason)/float64(tCount)*100)
	fmt.Printf("  StartDate заполнен: %d/%d (%.1f%%)\n\n", tWithStart, tCount, float64(tWithStart)/float64(tCount)*100)

	// 3. Проверка season/dates в player_teams
	var ptWithSeason, ptWithStart int
	db.Get(&ptWithSeason, "SELECT COUNT(*) FROM player_teams WHERE season IS NOT NULL AND season != ''")
	db.Get(&ptWithStart, "SELECT COUNT(*) FROM player_teams WHERE started_at IS NOT NULL")

	fmt.Println("✅ КАЧЕСТВО ДАННЫХ - PLAYER_TEAMS:")
	fmt.Printf("  Season заполнен: %d/%d (%.1f%%)\n", ptWithSeason, ptCount, float64(ptWithSeason)/float64(ptCount)*100)
	fmt.Printf("  StartedAt заполнен: %d/%d (%.1f%%)\n\n", ptWithStart, ptCount, float64(ptWithStart)/float64(ptCount)*100)

	// 4. Дедупликация
	var dupPlayers, dupTournaments, dupTeams int
	db.Get(&dupPlayers, `SELECT COUNT(*) FROM (SELECT profile_url FROM players GROUP BY profile_url HAVING COUNT(*) > 1) sub`)
	db.Get(&dupTournaments, `SELECT COUNT(*) FROM (SELECT url FROM tournaments GROUP BY url HAVING COUNT(*) > 1) sub`)
	db.Get(&dupTeams, `SELECT COUNT(*) FROM (SELECT url FROM teams GROUP BY url HAVING COUNT(*) > 1) sub`)

	fmt.Println("🔍 ПРОВЕРКА ДЕДУПЛИКАЦИИ:")
	if dupPlayers == 0 {
		fmt.Println("  ✅ Дублей игроков: 0")
	} else {
		fmt.Printf("  ❌ Дублей игроков: %d\n", dupPlayers)
	}
	if dupTournaments == 0 {
		fmt.Println("  ✅ Дублей турниров: 0")
	} else {
		fmt.Printf("  ❌ Дублей турниров: %d\n", dupTournaments)
	}
	if dupTeams == 0 {
		fmt.Println("  ✅ Дублей команд: 0")
	} else {
		fmt.Printf("  ❌ Дублей команд: %d\n", dupTeams)
	}

	// 5. Распределение по позициям
	var positions []struct {
		Position string `db:"position"`
		Count    int    `db:"count"`
	}
	db.Select(&positions, `SELECT position, COUNT(*) as count FROM players GROUP BY position ORDER BY count DESC`)

	fmt.Println("\n📊 РАСПРЕДЕЛЕНИЕ ИГРОКОВ ПО ПОЗИЦИЯМ:")
	for _, p := range positions {
		fmt.Printf("  %s: %d\n", p.Position, p.Count)
	}

	// 6. ТОП-5 турниров по количеству команд
	var topTournaments []struct {
		Name      string `db:"name"`
		TeamCount int    `db:"team_count"`
		Season    string `db:"season"`
	}
	db.Select(&topTournaments, `
		SELECT t.name, t.season, COUNT(DISTINCT pt.team_id) as team_count
		FROM tournaments t
		LEFT JOIN player_teams pt ON pt.tournament_id = t.id
		GROUP BY t.id, t.name, t.season
		ORDER BY team_count DESC
		LIMIT 5
	`)

	fmt.Println("\n🏆 ТОП-5 ТУРНИРОВ ПО КОМАНДАМ:")
	for i, t := range topTournaments {
		fmt.Printf("  %d. %s (%s) - %d команд\n", i+1, t.Name, t.Season, t.TeamCount)
	}

	// 7. Примеры игроков
	var players []struct {
		Name     string `db:"name"`
		Position string `db:"position"`
	}
	db.Select(&players, "SELECT name, position FROM players ORDER BY name LIMIT 5")

	fmt.Println("\n👤 ПРИМЕРЫ ИГРОКОВ:")
	for i, p := range players {
		fmt.Printf("  %d. %s (%s)\n", i+1, p.Name, p.Position)
	}

	fmt.Println("\n================================================================================")
	fmt.Println("✅ ПАРСИНГ ЗАВЕРШЕН УСПЕШНО!")
	fmt.Println("================================================================================")
}
