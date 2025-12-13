package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL не установлен")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Количество турниров
	var tournamentsCount int
	err = db.QueryRow("SELECT COUNT(*) FROM tournaments").Scan(&tournamentsCount)
	if err != nil {
		log.Fatal(err)
	}

	// Количество турниров по сезонам
	rows, err := db.Query(`
		SELECT season, COUNT(*) as count 
		FROM tournaments 
		GROUP BY season 
		ORDER BY season DESC
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("📊 СТАТИСТИКА ТУРНИРОВ:")
	fmt.Printf("Всего турниров: %d\n\n", tournamentsCount)
	fmt.Println("По сезонам:")
	for rows.Next() {
		var season string
		var count int
		if err := rows.Scan(&season, &count); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %d турниров\n", season, count)
	}

	// Примеры турниров с датами
	fmt.Println("\n📅 ПРИМЕРЫ ТУРНИРОВ С ДАТАМИ:")
	rows2, err := db.Query(`
		SELECT id, name, season, start_date, end_date, is_ended 
		FROM tournaments 
		ORDER BY start_date DESC NULLS LAST
		LIMIT 10
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var id, name, season string
		var startDate, endDate sql.NullTime
		var isEnded bool
		if err := rows2.Scan(&id, &name, &season, &startDate, &endDate, &isEnded); err != nil {
			log.Fatal(err)
		}

		startStr := "NULL"
		if startDate.Valid {
			startStr = startDate.Time.Format("02.01.2006")
		}
		endStr := "NULL"
		if endDate.Valid {
			endStr = endDate.Time.Format("02.01.2006")
		}

		endedStr := "❌"
		if isEnded {
			endedStr = "✅"
		}

		fmt.Printf("  [%s] %s | %s | %s - %s | Закончен: %s\n", 
			id, name, season, startStr, endStr, endedStr)
	}

	// Количество игроков
	var playersCount int
	err = db.QueryRow("SELECT COUNT(*) FROM players").Scan(&playersCount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n👥 Всего игроков: %d\n", playersCount)
}
