package report

import (
	"bytes"
	"fmt"
	"html/template"

	"context"

	domainPlayer "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player"
	domainPlayerStats "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_statistics"
	domainPlayerTeam "github.com/Daniil-Sakharov/HockeyProject/internal/domain/player_team"
	domainTeam "github.com/Daniil-Sakharov/HockeyProject/internal/domain/team"
	domainTournament "github.com/Daniil-Sakharov/HockeyProject/internal/domain/tournament"
	"github.com/Daniil-Sakharov/HockeyProject/internal/service/bot/report/svg"
)

// SVGCharts содержит сгенерированные SVG графики
type SVGCharts struct {
	GoalsTypePie template.HTML
	PeriodBar    template.HTML
	ProgressLine template.HTML
	ProfileRadar template.HTML
}

// Service сервис генерации HTML отчетов
type Service struct {
	dataCollector *DataCollector
	template      *template.Template
}

// NewService создает новый сервис отчетов
func NewService(
	playerRepo domainPlayer.Repository,
	statsRepo domainPlayerStats.Repository,
	playerTeamRepo domainPlayerTeam.Repository,
	teamRepo domainTeam.Repository,
	tournamentRepo domainTournament.Repository,
) *Service {
	dc := NewDataCollector(playerRepo, statsRepo, playerTeamRepo, teamRepo, tournamentRepo)

	funcMap := template.FuncMap{
		"formatFloat": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		"formatPct": func(f float64) string {
			return fmt.Sprintf("%.1f%%", f)
		},
		"plusMinusFormat": func(v int) string {
			if v > 0 {
				return fmt.Sprintf("+%d", v)
			}
			return fmt.Sprintf("%d", v)
		},
	}

	tmpl := template.Must(template.New("report").Funcs(funcMap).Parse(HTMLTemplate))

	return &Service{
		dataCollector: dc,
		template:      tmpl,
	}
}

// GeneratePlayerReport генерирует HTML отчет для игрока
func (s *Service) GeneratePlayerReport(ctx context.Context, playerID string) ([]byte, string, error) {
	report, err := s.dataCollector.CollectFullReport(ctx, playerID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to collect report data: %w", err)
	}

	charts := s.generateSVGCharts(report)

	templateData := struct {
		Report *FullPlayerReport
		Charts SVGCharts
	}{
		Report: report,
		Charts: charts,
	}

	var buf bytes.Buffer
	if err := s.template.Execute(&buf, templateData); err != nil {
		return nil, "", fmt.Errorf("failed to execute template: %w", err)
	}

	filename := fmt.Sprintf("%s_report.html", transliterate(report.Player.Name))

	return buf.Bytes(), filename, nil
}

// generateSVGCharts генерирует все SVG графики
func (s *Service) generateSVGCharts(report *FullPlayerReport) SVGCharts {
	charts := SVGCharts{}

	if !report.HasStats {
		return charts
	}

	// Круговая диаграмма - типы голов
	if report.HasDetailedStats {
		labels := []string{"В равных", "В большинстве", "В меньшинстве"}
		values := []int{
			report.GoalsByType.EvenStrength,
			report.GoalsByType.PowerPlay,
			report.GoalsByType.ShortHanded,
		}
		charts.GoalsTypePie = template.HTML(svg.GeneratePieChart(labels, values, nil))
	}

	// Столбчатая диаграмма - голы по периодам
	periodLabels := []string{"1 период", "2 период", "3 период", "OT"}
	periodValues := []int{
		report.GoalsByPeriod.Period1,
		report.GoalsByPeriod.Period2,
		report.GoalsByPeriod.Period3,
		report.GoalsByPeriod.Overtime,
	}
	charts.PeriodBar = template.HTML(svg.GenerateBarChart(periodLabels, periodValues, nil))

	// Линейный график - прогресс по сезонам
	if report.HasMultipleSeasons && len(report.SeasonStats) > 1 {
		var seasonLabels []string
		var goalsValues, pointsValues []int

		for _, s := range report.SeasonStats {
			seasonLabels = append(seasonLabels, s.Season)
			goalsValues = append(goalsValues, s.Goals)
			pointsValues = append(pointsValues, s.Points)
		}

		datasets := []svg.LineDataset{
			{Label: "Очки", Values: pointsValues, Color: svg.Colors.Primary},
			{Label: "Голы", Values: goalsValues, Color: svg.Colors.Accent},
		}
		charts.ProgressLine = template.HTML(svg.GenerateLineChart(seasonLabels, datasets, nil))
	}

	// Radar диаграмма - профиль игрока
	radarLabels := []string{"Голы", "Пасы", "+/-", "Хет-трики", "Поб. голы"}
	radarValues := []float64{
		float64(report.TotalStats.TotalGoals),
		float64(report.TotalStats.TotalAssists),
		float64(abs(report.TotalStats.TotalPlusMinus)),
		float64(report.TotalStats.TotalHatTricks) * 10, // Масштабируем для наглядности
		float64(report.TotalStats.TotalWinningGoals),
	}
	charts.ProfileRadar = template.HTML(svg.GenerateRadarChart(radarLabels, radarValues, nil))

	return charts
}

// transliterate транслитерирует кириллицу в латиницу для имени файла
func transliterate(s string) string {
	translit := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "",
		'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo",
		'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
		'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
		'Ф': "F", 'Х': "H", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "",
		'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
		' ': "_",
	}

	var result []byte
	for _, r := range s {
		if t, ok := translit[r]; ok {
			result = append(result, t...)
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result = append(result, byte(r))
		}
	}
	return string(result)
}
