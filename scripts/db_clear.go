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
		log.Fatalf("❌ Не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	fmt.Println("🗑️  Очистка БД...")
	fmt.Println()

	tables := []string{"player_teams", "players", "teams", "tournaments"}

	for _, table := range tables {
		_, err := db.Exec("TRUNCATE " + table + " CASCADE")
		if err != nil {
			log.Fatalf("❌ Ошибка при очистке %s: %v", table, err)
		}
		fmt.Printf("  ✅ %s очищена\n", table)
	}

	fmt.Println()
	fmt.Println("✅ БД полностью очищена и готова к парсингу!")
}
