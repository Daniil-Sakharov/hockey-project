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

	// 1. Проверяем дубликаты по (ID + Domain)
	fmt.Println("🔍 Проверка дубликатов по (ID, Domain):")
	rows, err := db.Query(`
		SELECT id, domain, COUNT(*) as count
		FROM tournaments 
		GROUP BY id, domain 
		HAVING COUNT(*) > 1
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dupCount := 0
	for rows.Next() {
		var id, domain string
		var count int
		rows.Scan(&id, &domain, &count)
		fmt.Printf("  ❌ ID: %s, Domain: %s, Count: %d\n", id, domain, count)
		dupCount++
	}

	if dupCount == 0 {
		fmt.Println("  ✅ Нет дубликатов по (ID, Domain)")
	}

	// 2. Проверяем одинаковые ID в разных доменах
	fmt.Println("\n🔍 Турниры с одинаковым ID в разных доменах:")
	rows2, err := db.Query(`
		SELECT id, array_agg(DISTINCT domain) as domains, COUNT(DISTINCT domain) as domain_count
		FROM tournaments
		GROUP BY id
		HAVING COUNT(DISTINCT domain) > 1
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	crossDomainCount := 0
	for rows2.Next() {
		var id string
		var domains interface{}
		var domainCount int
		rows2.Scan(&id, &domains, &domainCount)
		fmt.Printf("  ⚠️  ID: %s присутствует в %d доменах: %v\n", id, domainCount, domains)
		crossDomainCount++
	}

	if crossDomainCount == 0 {
		fmt.Println("  ✅ Нет турниров с одинаковым ID в разных доменах")
	}

	// 3. Общая статистика
	var totalTournaments, uniqueIDs, uniqueDomains int
	db.QueryRow("SELECT COUNT(*) FROM tournaments").Scan(&totalTournaments)
	db.QueryRow("SELECT COUNT(DISTINCT id) FROM tournaments").Scan(&uniqueIDs)
	db.QueryRow("SELECT COUNT(DISTINCT domain) FROM tournaments").Scan(&uniqueDomains)

	fmt.Println("\n📊 Общая статистика:")
	fmt.Printf("  Всего записей: %d\n", totalTournaments)
	fmt.Printf("  Уникальных ID: %d\n", uniqueIDs)
	fmt.Printf("  Доменов: %d\n", uniqueDomains)

	if totalTournaments == uniqueIDs {
		fmt.Println("  ✅ Все ID уникальны!")
	} else {
		fmt.Printf("  ⚠️  Разница: %d дубликатов\n", totalTournaments-uniqueIDs)
	}
}
