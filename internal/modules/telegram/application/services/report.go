package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/Daniil-Sakharov/HockeyProject/internal/modules/shared/charts"
)

// ReportService сервис генерации HTML отчётов
type ReportService struct {
	dataCollector *DataCollector
	template      *template.Template
}

// ReportRepository интерфейс для получения данных отчёта
type ReportRepository interface {
	GetFullReport(ctx context.Context, playerID string) (*FullPlayerReport, error)
}

// NewReportService создает новый сервис отчётов
func NewReportService(repo ReportRepository) *ReportService {
	dataCollector := NewDataCollector(repo)

	funcMap := template.FuncMap{
		"formatFloat": func(f float64) string {
			return fmt.Sprintf("%.2f", f)
		},
		"plusMinusFormat": func(v int) string {
			if v > 0 {
				return fmt.Sprintf("+%d", v)
			}
			return fmt.Sprintf("%d", v)
		},
	}

	tmpl := template.Must(template.New("report").Funcs(funcMap).Parse(reportHTMLTemplate))

	return &ReportService{
		dataCollector: dataCollector,
		template:      tmpl,
	}
}

// GenerateReport генерирует HTML отчёт для игрока
func (s *ReportService) GenerateReport(ctx context.Context, playerID string) ([]byte, string, error) {
	report, err := s.dataCollector.CollectFullReport(ctx, playerID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to collect report data: %w", err)
	}

	svgCharts := s.generateCharts(report)

	data := struct {
		Report *FullPlayerReport
		Charts SVGCharts
	}{
		Report: report,
		Charts: svgCharts,
	}

	var buf bytes.Buffer
	if err := s.template.Execute(&buf, data); err != nil {
		return nil, "", fmt.Errorf("failed to execute template: %w", err)
	}

	filename := fmt.Sprintf("%s_report.html", transliterate(report.Player.Name))
	return buf.Bytes(), filename, nil
}

// SVGCharts содержит сгенерированные SVG графики
type SVGCharts struct {
	GoalsTypePie template.HTML
	PeriodBar    template.HTML
	ProgressLine template.HTML
	ProfileRadar template.HTML
}

func (s *ReportService) generateCharts(report *FullPlayerReport) SVGCharts {
	result := SVGCharts{}

	if !report.HasStats {
		return result
	}

	// Круговая диаграмма - типы голов
	if report.HasDetailedStats {
		labels := []string{"В равных", "В большинстве", "В меньшинстве"}
		values := []int{
			report.GoalsByType.EvenStrength,
			report.GoalsByType.PowerPlay,
			report.GoalsByType.ShortHanded,
		}
		result.GoalsTypePie = template.HTML(charts.GeneratePieChart(labels, values, nil)) //nolint:gosec
	}

	// Столбчатая диаграмма - голы по периодам
	periodLabels := []string{"1 период", "2 период", "3 период", "OT"}
	periodValues := []int{
		report.GoalsByPeriod.Period1,
		report.GoalsByPeriod.Period2,
		report.GoalsByPeriod.Period3,
		report.GoalsByPeriod.Overtime,
	}
	result.PeriodBar = template.HTML(charts.GenerateBarChart(periodLabels, periodValues, nil)) //nolint:gosec

	// Линейный график - прогресс по сезонам
	if report.HasMultipleSeasons && len(report.SeasonStats) > 1 {
		var seasonLabels []string
		var goalsValues, pointsValues []int

		for _, ss := range report.SeasonStats {
			seasonLabels = append(seasonLabels, ss.Season)
			goalsValues = append(goalsValues, ss.Goals)
			pointsValues = append(pointsValues, ss.Points)
		}

		datasets := []charts.LineDataset{
			{Label: "Очки", Values: pointsValues, Color: charts.Colors.Primary},
			{Label: "Голы", Values: goalsValues, Color: charts.Colors.Accent},
		}
		result.ProgressLine = template.HTML(charts.GenerateLineChart(seasonLabels, datasets, nil)) //nolint:gosec
	}

	// Radar диаграмма - профиль игрока
	radarLabels := []string{"Голы", "Пасы", "+/-", "Хет-трики", "Поб. голы"}
	radarValues := []float64{
		float64(report.TotalStats.TotalGoals),
		float64(report.TotalStats.TotalAssists),
		float64(abs(report.TotalStats.TotalPlusMinus)),
		float64(report.TotalStats.TotalHatTricks) * 10,
		float64(report.TotalStats.TotalWinningGoals),
	}
	result.ProfileRadar = template.HTML(charts.GenerateRadarChart(radarLabels, radarValues, nil)) //nolint:gosec

	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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
