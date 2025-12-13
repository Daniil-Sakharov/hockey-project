package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	godotenv.Load()
	
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Проверяем конкретные ID турниров из check_all.md
	ids := []string{
		"16756891", "16756910", "16756905", // cfo 2025/2026
		"16724630", "16724632", "16724635", // szfo 2025/2026
		"18626069", "22840769", "16735091", "16735092", "16735093", // pfo 2025/2026
	}
	
	fmt.Println("\n📊 ПРОВЕРКА ТУРНИРОВ ИЗ CHECK_ALL.MD (сезон 2025/2026):")
	fmt.Println("================================================================================")
	
	for _, id := range ids {
		var name, domain, season sql.NullString
		err := db.QueryRow("SELECT name, domain, season FROM tournaments WHERE id = $1", id).Scan(&name, &domain, &season)
		if err == sql.ErrNoRows {
			fmt.Printf("  ❌ ID %s: НЕ НАЙДЕН\n", id)
		} else if err != nil {
			fmt.Printf("  ⚠️  ID %s: ошибка: %v\n", id, err)
		} else {
			fmt.Printf("  ✅ ID %s: %s [%s] %s\n", id, name.String, season.String, domain.String)
		}
	}

	// Турниры по доменам
	fmt.Println("\n📊 ТУРНИРЫ ПО ДОМЕНАМ:")
	fmt.Println("================================================================================")
	rows, _ := db.Query("SELECT domain, COUNT(*) as count FROM tournaments GROUP BY domain ORDER BY count DESC")
	defer rows.Close()
	for rows.Next() {
		var domain string
		var count int
		rows.Scan(&domain, &count)
		fmt.Printf("  %-40s %d\n", domain, count)
	}

	// Проверим турниры cfo.fhr.ru сезона 2025/2026
	fmt.Println("\n📊 ТУРНИРЫ cfo.fhr.ru с сезоном 2025/2026:")
	fmt.Println("================================================================================")
	rows5, _ := db.Query("SELECT id, name, season FROM tournaments WHERE domain = 'https://cfo.fhr.ru' AND season = '2025/2026'")
	defer rows5.Close()
	count5 := 0
	for rows5.Next() {
		var id, name, season string
		rows5.Scan(&id, &name, &season)
		fmt.Printf("  [%s] %s | ID: %s\n", season, name, id)
		count5++
	}
	if count5 == 0 {
		fmt.Println("  ❌ НЕТ турниров!")
	}

	// Проверим конкретные турниры для cfo.fhr.ru
	fmt.Println("\n📊 ВСЕ ТУРНИРЫ cfo.fhr.ru:")
	fmt.Println("================================================================================")
	rows2, _ := db.Query("SELECT id, name, season FROM tournaments WHERE domain = 'https://cfo.fhr.ru' ORDER BY season DESC, name")
	defer rows2.Close()
	for rows2.Next() {
		var id, name, season string
		rows2.Scan(&id, &name, &season)
		fmt.Printf("  [%s] %s | ID: %s\n", season, name, id)
	}

	// Проверим турниры для pfo.fhr.ru
	fmt.Println("\n📊 ВСЕ ТУРНИРЫ pfo.fhr.ru:")
	fmt.Println("================================================================================")
	rows3, _ := db.Query("SELECT id, name, season FROM tournaments WHERE domain = 'https://pfo.fhr.ru' ORDER BY season DESC, name")
	defer rows3.Close()
	for rows3.Next() {
		var id, name, season string
		rows3.Scan(&id, &name, &season)
		fmt.Printf("  [%s] %s | ID: %s\n", season, name, id)
	}

	// Уникальные сезоны
	fmt.Println("\n📊 УНИКАЛЬНЫЕ СЕЗОНЫ В БД:")
	fmt.Println("================================================================================")
	rows6, _ := db.Query("SELECT DISTINCT season, COUNT(*) as count FROM tournaments GROUP BY season ORDER BY season DESC")
	defer rows6.Close()
	for rows6.Next() {
		var season sql.NullString
		var count int
		rows6.Scan(&season, &count)
		if season.Valid && season.String != "" {
			fmt.Printf("  [%s]: %d турниров\n", season.String, count)
		} else {
			fmt.Printf("  [ПУСТО/NULL]: %d турниров\n", count)
		}
	}

	// Проверка конкретного турнира 5153197
	fmt.Println("\n📊 ТУРНИР 5153197 (Первенство ППФО 2021/2022):")
	fmt.Println("================================================================================")
	
	// Общее количество команд
	var totalTeams int
	db.QueryRow("SELECT COUNT(*) FROM teams").Scan(&totalTeams)
	fmt.Printf("  Всего команд в БД: %d\n", totalTeams)
	
	// Связи player_teams для турнира 5153197
	var playerTeamsCount int
	db.QueryRow("SELECT COUNT(*) FROM player_teams WHERE tournament_id = '5153197'").Scan(&playerTeamsCount)
	fmt.Printf("  Связей player_teams для турнира: %d\n", playerTeamsCount)
	
	// Уникальные команды через player_teams
	var uniqueTeamsInTournament int
	db.QueryRow("SELECT COUNT(DISTINCT team_id) FROM player_teams WHERE tournament_id = '5153197'").Scan(&uniqueTeamsInTournament)
	fmt.Printf("  Уникальных команд в турнире (через player_teams): %d\n", uniqueTeamsInTournament)
	
	// Связи игрок-турнир
	var playerTournamentCount int
	db.QueryRow("SELECT COUNT(*) FROM player_tournaments WHERE tournament_id = '5153197'").Scan(&playerTournamentCount)
	fmt.Printf("  Связей player_tournaments: %d\n", playerTournamentCount)
	
	// Уникальные игроки в турнире
	var uniquePlayers int
	db.QueryRow("SELECT COUNT(DISTINCT player_id) FROM player_tournaments WHERE tournament_id = '5153197'").Scan(&uniquePlayers)
	fmt.Printf("  Уникальных игроков в турнире: %d\n", uniquePlayers)
}
