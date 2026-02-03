package main

import (
	"context"
	"log"

	juniorCalendar "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/application/orchestrators/junior/calendar"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/repositories"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior"
	jrCalendarParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/calendar"
	jrGameParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/game"
	jrStandingsParser "github.com/Daniil-Sakharov/HockeyProject/internal/modules/parsing/infrastructure/sources/junior/standings"
	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/di"
	"github.com/Daniil-Sakharov/HockeyProject/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	if err := logger.Init("debug", false, nil); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	container := di.NewContainer()
	defer func() { _ = container.Close() }()

	db, err := container.DB(ctx)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	// Repositories
	matchRepo, _ := container.MatchRepository(ctx)
	matchEventRepo, _ := container.MatchEventRepository(ctx)
	matchLineupRepo, _ := container.MatchLineupRepository(ctx)
	standingRepo, _ := container.StandingRepository(ctx)
	tournamentRepo := repositories.NewTournamentPostgres(db)
	teamRepo, _ := container.ParsingTeamRepository(ctx)
	playerRepo, _ := container.ParsingPlayerRepository(ctx)

	// Parsers
	juniorClient := junior.NewClient()
	calendarParser := jrCalendarParser.NewParser(juniorClient)
	gameParser := jrGameParser.NewParser(juniorClient)
	standingsParser := jrStandingsParser.NewParser(juniorClient)

	// Config
	config := &testConfig{}

	// Orchestrator
	orch := juniorCalendar.NewOrchestrator(
		juniorClient, // HTTP клиент для AJAX-запросов
		calendarParser,
		gameParser,
		standingsParser,
		nil, // profileParser
		matchRepo,
		matchEventRepo,
		matchLineupRepo,
		standingRepo,
		tournamentRepo,
		teamRepo,
		playerRepo,
		config,
	)

	// Получаем только 1 турнир для теста
	tournaments, err := tournamentRepo.GetBySource(ctx, "junior")
	if err != nil {
		log.Fatalf("Failed to get tournaments: %v", err)
	}

	log.Printf("📊 Found %d tournaments, testing with first one", len(tournaments))

	if len(tournaments) == 0 {
		log.Fatal("No tournaments found")
	}

	// Берём турнир с pfo.fhr.ru для теста (там точно есть матчи)
	var testTournament *struct {
		ID, Name, URL, Domain string
	}
	for _, t := range tournaments {
		if t.Domain == "https://pfo.fhr.ru" {
			testTournament = &struct{ ID, Name, URL, Domain string }{t.ID, t.Name, t.URL, t.Domain}
			break
		}
	}
	if testTournament == nil {
		// Fallback to first tournament
		t := tournaments[0]
		testTournament = &struct{ ID, Name, URL, Domain string }{t.ID, t.Name, t.URL, t.Domain}
	}

	log.Printf("🏆 Testing with: %s", testTournament.Name)
	log.Printf("   URL: %s%s", testTournament.Domain, testTournament.URL)
	log.Println("")

	// Удаляем старые данные для этого турнира
	db.Exec("DELETE FROM match_lineups WHERE match_id IN (SELECT id FROM matches WHERE tournament_id = $1)", testTournament.ID)
	db.Exec("DELETE FROM match_events WHERE match_id IN (SELECT id FROM matches WHERE tournament_id = $1)", testTournament.ID)
	db.Exec("DELETE FROM matches WHERE tournament_id = $1", testTournament.ID)
	db.Exec("DELETE FROM team_standings WHERE tournament_id = $1", testTournament.ID)

	// Запускаем обработку турнира напрямую
	fullURL := testTournament.Domain + testTournament.URL

	// 1. Парсим таблицу
	log.Println("📊 Step 1: Parsing standings...")
	standings, err := standingsParser.Parse(fullURL)
	if err != nil {
		log.Printf("   ❌ Error: %v", err)
	} else {
		log.Printf("   ✅ Found %d teams in standings", len(standings))
	}

	// 2. Парсим календарь
	log.Println("")
	log.Println("🗓️ Step 2: Parsing calendar...")
	matches, err := calendarParser.Parse(fullURL)
	if err != nil {
		log.Printf("   ❌ Error: %v", err)
	} else {
		finished := 0
		for _, m := range matches {
			if m.Status == "finished" {
				finished++
			}
		}
		log.Printf("   ✅ Found %d matches (%d finished)", len(matches), finished)
	}

	// 3. Запускаем полный orchestrator для этого турнира
	log.Println("")
	log.Println("🚀 Step 3: Running full orchestrator...")

	// Создаём временный orchestrator, который обработает только 1 турнир
	if err := orch.Run(ctx); err != nil {
		log.Printf("   ❌ Error: %v", err)
	}

	// 4. Проверяем результаты
	log.Println("")
	log.Println("📊 RESULTS:")
	log.Println("============================================")

	type counter struct {
		Count int `db:"count"`
	}
	var c counter

	db.Get(&c, "SELECT COUNT(*) as count FROM matches WHERE source = 'junior'")
	log.Printf("   Matches: %d", c.Count)

	db.Get(&c, "SELECT COUNT(*) as count FROM matches WHERE source = 'junior' AND details_parsed = true")
	log.Printf("   Matches with details: %d", c.Count)

	db.Get(&c, "SELECT COUNT(*) as count FROM match_events WHERE source = 'junior' AND event_type = 'goal'")
	log.Printf("   Goals: %d", c.Count)

	db.Get(&c, "SELECT COUNT(*) as count FROM match_events WHERE source = 'junior' AND event_type = 'penalty'")
	log.Printf("   Penalties: %d", c.Count)

	db.Get(&c, "SELECT COUNT(*) as count FROM match_lineups WHERE source = 'junior'")
	log.Printf("   Lineups: %d", c.Count)

	db.Get(&c, "SELECT COUNT(*) as count FROM team_standings WHERE source = 'junior'")
	log.Printf("   Standings: %d", c.Count)

	log.Println("============================================")

	// Показываем пример данных
	log.Println("")
	log.Println("📋 Sample data:")

	type goal struct {
		Period  int    `db:"period"`
		Minutes int    `db:"time_minutes"`
		Seconds int    `db:"time_seconds"`
	}
	var goals []goal
	db.Select(&goals, "SELECT period, time_minutes, time_seconds FROM match_events WHERE source = 'junior' AND event_type = 'goal' LIMIT 5")
	if len(goals) > 0 {
		log.Println("   Goals:")
		for _, g := range goals {
			log.Printf("      Period %d, %d:%02d", g.Period, g.Minutes, g.Seconds)
		}
	}

	type penalty struct {
		Minutes int    `db:"penalty_minutes"`
		Reason  string `db:"penalty_reason"`
	}
	var penalties []penalty
	db.Select(&penalties, "SELECT penalty_minutes, penalty_reason FROM match_events WHERE source = 'junior' AND event_type = 'penalty' LIMIT 5")
	if len(penalties) > 0 {
		log.Println("   Penalties:")
		for _, p := range penalties {
			log.Printf("      %d min - %s", p.Minutes, p.Reason)
		}
	}

	type lineup struct {
		Position string `db:"position"`
		Count    int    `db:"count"`
	}
	var lineups []lineup
	db.Select(&lineups, "SELECT position, COUNT(*) as count FROM match_lineups WHERE source = 'junior' GROUP BY position")
	if len(lineups) > 0 {
		log.Println("   Lineups by position:")
		for _, l := range lineups {
			log.Printf("      %s: %d players", l.Position, l.Count)
		}
	}

	log.Println("")
	log.Println("✅ Test completed!")
}

type testConfig struct{}

func (c *testConfig) RequestDelay() int      { return 200 }
func (c *testConfig) TournamentWorkers() int { return 1 }
func (c *testConfig) GameWorkers() int       { return 2 }
func (c *testConfig) ParseProtocol() bool    { return true }
func (c *testConfig) ParseLineups() bool     { return true }
func (c *testConfig) SkipExisting() bool     { return false }
func (c *testConfig) MaxTournaments() int    { return 1 }
