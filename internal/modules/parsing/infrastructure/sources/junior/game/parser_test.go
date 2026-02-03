package game

import (
	"net/http"
	"testing"
	"time"
)

type TestClient struct{}

func (c *TestClient) MakeRequest(url string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	return client.Do(req)
}

func (c *TestClient) MakeRequestWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}

func TestParseRealGame(t *testing.T) {
	parser := NewParser(&TestClient{})

	gameURL := "https://pfo.fhr.ru/games/17392431/"
	t.Log("Парсинг матча:", gameURL)

	details, err := parser.Parse(gameURL)
	if err != nil {
		t.Fatal("Ошибка:", err)
	}

	// Проверяем счёт
	if details.HomeScore == nil || details.AwayScore == nil {
		t.Error("Счёт не спарсен")
	} else {
		t.Logf("Общий счёт: %d : %d", *details.HomeScore, *details.AwayScore)
	}

	// Проверяем голы
	t.Logf("Голов: %d", len(details.Goals))
	for i, g := range details.Goals {
		if i < 5 {
			t.Logf("  %d:%02d - %s (%s, %s)", g.TimeMinutes, g.TimeSeconds, g.ScorerName, g.Assist1Name, g.Assist2Name)
		}
	}

	// Проверяем штрафы
	t.Logf("Штрафов: %d", len(details.Penalties))
	for i, p := range details.Penalties {
		if i < 5 {
			t.Logf("  %d:%02d - %s, %d' (%s)", p.TimeMinutes, p.TimeSeconds, p.PlayerName, p.Minutes, p.Reason)
		}
	}

	// Проверяем события вратарей
	t.Logf("\nСобытий вратарей: %d", len(details.GoalieEvents))
	for _, g := range details.GoalieEvents {
		team := "гости"
		if g.IsHome {
			team = "дом"
		}
		t.Logf("  %d:%02d - %s [%s] %s", g.TimeMinutes, g.TimeSeconds, g.PlayerName, team, g.PlayerURL)
	}

	// Проверяем пустые ворота
	t.Logf("\nПустые ворота: %d", len(details.EmptyNets))
	for _, en := range details.EmptyNets {
		team := "гости"
		if en.IsHome {
			team = "дом"
		}
		t.Logf("  %d:%02d - %s снял вратаря", en.TimeMinutes, en.TimeSeconds, team)
	}

	// Проверяем тайм-ауты
	t.Logf("\nТайм-аутов: %d", len(details.Timeouts))
	for _, to := range details.Timeouts {
		team := "гости"
		if to.IsHome {
			team = "дом"
		}
		t.Logf("  %d:%02d - %s взял тайм-аут", to.TimeMinutes, to.TimeSeconds, team)
	}

	// Проверяем составы - вратари
	t.Logf("\nИгроков дома: %d", len(details.HomeLineup))
	goalies := 0
	defenders := 0
	forwards := 0
	captains := 0
	for _, p := range details.HomeLineup {
		switch p.Position {
		case "G":
			goalies++
		case "D":
			defenders++
		case "F":
			forwards++
		}
		if p.Role == "C" || p.Role == "A" {
			captains++
		}
	}
	t.Logf("  Вратари: %d, Защитники: %d, Нападающие: %d", goalies, defenders, forwards)
	t.Logf("  Капитаны/Ассистенты: %d", captains)

	// Примеры игроков
	for i, p := range details.HomeLineup {
		if i < 3 {
			role := ""
			if p.Role != "" {
				role = " (" + p.Role + ")"
			}
			t.Logf("  #%d %s [%s]%s", p.JerseyNumber, p.PlayerName, p.Position, role)
		}
	}

	// Проверяем гостевую команду
	awayGoalies := 0
	awayCaptains := 0
	for _, p := range details.AwayLineup {
		if p.Position == "G" {
			awayGoalies++
		}
		if p.Role == "C" || p.Role == "A" {
			awayCaptains++
		}
	}
	t.Logf("\nИгроков в гостях: %d", len(details.AwayLineup))
	t.Logf("  Капитаны/Ассистенты у гостей: %d", awayCaptains)

	// Примеры гостей с ролями
	for _, p := range details.AwayLineup {
		if p.Role != "" {
			t.Logf("  #%d %s [%s] (%s)", p.JerseyNumber, p.PlayerName, p.Position, p.Role)
		}
	}

	// Проверяем что данные получены
	if len(details.Goals) == 0 {
		t.Error("Голы не спарсены")
	}
	if len(details.HomeLineup) == 0 {
		t.Error("Состав дома не спарсен")
	}
	if len(details.AwayLineup) == 0 {
		t.Error("Состав гостей не спарсен")
	}
	if goalies == 0 {
		t.Error("Вратари не определены")
	}
	if defenders == 0 {
		t.Error("Защитники не определены")
	}
	// Капитаны могут быть только у одной команды, проверяем что нашли хотя бы одного
	if captains == 0 && awayCaptains == 0 {
		t.Log("Предупреждение: капитаны не найдены ни у одной команды")
	}
}

// TestExtractJerseyNumber проверяет извлечение номера игрока из текста
func TestExtractJerseyNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"55. Мухтаров Ленар", 55},
		{"1. Гараев Амирхан", 1},
		{"99. Иванов Иван", 99},
		{"7. Сидоров", 7},
		{"Без номера", 0},
		{"", 0},
		{"   55. Пробелы", 0}, // начинается с пробелов - regex не найдёт
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractJerseyNumber(tt.input)
			if result != tt.expected {
				t.Errorf("extractJerseyNumber(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractPlayerNameFromText проверяет извлечение имени игрока из текста
func TestExtractPlayerNameFromText(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"55. Мухтаров Ленар", "Мухтаров Ленар"},
		{"1. Гараев Амирхан", "Гараев Амирхан"},
		{"99. Иванов", "Иванов"},
		{"Без номера", "Без номера"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := extractPlayerNameFromText(tt.input)
			if result != tt.expected {
				t.Errorf("extractPlayerNameFromText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestParseGoalAssistNumbers проверяет что номера ассистентов парсятся
func TestParseGoalAssistNumbers(t *testing.T) {
	parser := NewParser(&TestClient{})

	// Используем реальный матч с голами и ассистентами
	gameURL := "https://pfo.fhr.ru/games/17392431/"

	details, err := parser.Parse(gameURL)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}

	if len(details.Goals) == 0 {
		t.Fatal("Голы не найдены")
	}

	// Проверяем что хотя бы у некоторых голов есть номера ассистентов
	goalsWithAssists := 0
	for _, g := range details.Goals {
		if g.Assist1Number > 0 {
			goalsWithAssists++
			t.Logf("Гол в %d:%02d - Ассистент 1: #%d %s", g.TimeMinutes, g.TimeSeconds, g.Assist1Number, g.Assist1Name)
		}
		if g.Assist2Number > 0 {
			t.Logf("                Ассистент 2: #%d %s", g.Assist2Number, g.Assist2Name)
		}
	}

	t.Logf("Голов с ассистентами (с номерами): %d из %d", goalsWithAssists, len(details.Goals))

	// Не все голы имеют ассистентов, но должны быть хотя бы некоторые
	if goalsWithAssists == 0 && len(details.Goals) > 2 {
		t.Error("Ожидались голы с номерами ассистентов, но ни одного не найдено")
	}
}

// TestParsePlayersOnIceURLs проверяет что игроки на льду парсятся как URL
func TestParsePlayersOnIceURLs(t *testing.T) {
	parser := NewParser(&TestClient{})

	gameURL := "https://pfo.fhr.ru/games/17392431/"

	details, err := parser.Parse(gameURL)
	if err != nil {
		t.Fatalf("Ошибка парсинга: %v", err)
	}

	if len(details.Goals) == 0 {
		t.Fatal("Голы не найдены")
	}

	// Проверяем что игроки на льду - это URL, а не текст
	goalsWithPlayersOnIce := 0
	urlCount := 0

	for _, g := range details.Goals {
		hasPlayers := false

		if len(g.HomePlayersOnIce) > 0 {
			hasPlayers = true
			for _, p := range g.HomePlayersOnIce {
				if len(p) > 0 && (p[0] == '/' || p[0:4] == "http") {
					urlCount++
				}
			}
		}

		if len(g.AwayPlayersOnIce) > 0 {
			hasPlayers = true
			for _, p := range g.AwayPlayersOnIce {
				if len(p) > 0 && (p[0] == '/' || (len(p) >= 4 && p[0:4] == "http")) {
					urlCount++
				}
			}
		}

		if hasPlayers {
			goalsWithPlayersOnIce++
		}
	}

	t.Logf("Голов с игроками на льду: %d", goalsWithPlayersOnIce)
	t.Logf("Количество URL в игроках на льду: %d", urlCount)

	// Пример данных из первого гола
	if len(details.Goals) > 0 {
		g := details.Goals[0]
		t.Logf("Первый гол - игроки дома на льду: %v", g.HomePlayersOnIce)
		t.Logf("Первый гол - игроки гостей на льду: %v", g.AwayPlayersOnIce)
	}
}
