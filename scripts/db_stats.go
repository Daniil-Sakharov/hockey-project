package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	var tournaments, teams, players, playerTeams, statsRecords, playersWithStats int
	db.Get(&tournaments, "SELECT COUNT(*) FROM tournaments")
	db.Get(&teams, "SELECT COUNT(*) FROM teams")
	db.Get(&players, "SELECT COUNT(*) FROM players")
	db.Get(&playerTeams, "SELECT COUNT(*) FROM player_teams")
	db.Get(&statsRecords, "SELECT COUNT(*) FROM player_statistics")
	db.Get(&playersWithStats, "SELECT COUNT(DISTINCT player_id) FROM player_statistics")

	// Процент покрытия статистики
	coveragePercent := 0.0
	if players > 0 {
		coveragePercent = float64(playersWithStats) * 100.0 / float64(players)
	}

	fmt.Println()
	fmt.Println("📊 СТАТИСТИКА БД:")
	fmt.Println("================================================================================")
	fmt.Printf("  🏆 Турниры:                    %d\n", tournaments)
	fmt.Printf("  🏒 Команды:                    %d\n", teams)
	fmt.Printf("  👤 Игроки:                     %d\n", players)
	fmt.Printf("  🔗 Связи игрок-турнир:         %d\n", playerTeams)
	fmt.Println()
	fmt.Println("📈 СТАТИСТИКА ИГРОКОВ:")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("  📊 Записей статистики:         %d\n", statsRecords)
	fmt.Printf("  ✅ Игроков со статистикой:     %d (%.2f%%)\n", playersWithStats, coveragePercent)
	fmt.Printf("  ❌ Игроков БЕЗ статистики:     %d (%.2f%%)\n", players-playersWithStats, 100.0-coveragePercent)
	fmt.Println("================================================================================")
	fmt.Println()
	
	if coveragePercent < 50.0 {
		fmt.Println("⚠️  Низкое покрытие статистики! Рекомендуется запустить: task stats:parse")
		fmt.Println()
	}
}
