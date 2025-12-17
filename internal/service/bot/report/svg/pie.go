package svg

import (
	"fmt"
	"strings"
)

// PieChartOptions настройки круговой диаграммы
type PieChartOptions struct {
	Size        Size
	Colors      []string
	ShowLegend  bool
	DonutMode   bool
	InnerRadius float64
}

// DefaultPieOptions возвращает настройки по умолчанию
func DefaultPieOptions() PieChartOptions {
	return PieChartOptions{
		Size:        Size{Width: 300, Height: 180},
		Colors:      ChartPalette,
		ShowLegend:  true,
		DonutMode:   true,
		InnerRadius: 0.5,
	}
}

// GeneratePieChart генерирует SVG круговую диаграмму
func GeneratePieChart(labels []string, values []int, opts *PieChartOptions) string {
	if len(labels) == 0 || len(values) == 0 {
		return generateEmptyChart("Нет данных")
	}

	total := sumInt(values)
	if total == 0 {
		return generateEmptyChart("Нет данных")
	}

	if opts == nil {
		defaultOpts := DefaultPieOptions()
		opts = &defaultOpts
	}

	width := opts.Size.Width
	height := opts.Size.Height

	// Центр и радиус
	legendWidth := 0
	if opts.ShowLegend {
		legendWidth = 100
	}

	chartAreaWidth := width - legendWidth
	cx := float64(chartAreaWidth) / 2
	cy := float64(height) / 2
	radius := min(cx, cy) - 10

	innerRadius := 0.0
	if opts.DonutMode {
		innerRadius = radius * opts.InnerRadius
	}

	var sb strings.Builder

	// SVG начало
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))

	// Фон
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white" rx="8"/>`, width, height))

	// Секторы
	startAngle := 0.0
	for i, val := range values {
		if val == 0 {
			continue
		}

		percentage := float64(val) / float64(total)
		sweepAngle := percentage * 360

		color := opts.Colors[i%len(opts.Colors)]

		if opts.DonutMode {
			sb.WriteString(generateDonutSegment(cx, cy, radius, innerRadius, startAngle, startAngle+sweepAngle, color))
		} else {
			path := describeArc(cx, cy, radius, startAngle, startAngle+sweepAngle)
			sb.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="white" stroke-width="2"/>`, path, color))
		}

		startAngle += sweepAngle
	}

	// Центральный текст для donut
	if opts.DonutMode {
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" dominant-baseline="middle" font-size="20" font-weight="700" fill="%s">%d</text>`,
			cx, cy-5, Colors.PrimaryDark, total))
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-size="10" fill="%s">всего</text>`,
			cx, cy+12, Colors.Gray))
	}

	// Легенда
	if opts.ShowLegend {
		legendX := chartAreaWidth + 10
		legendY := 20
		for i, label := range labels {
			if i >= len(values) {
				break
			}
			color := opts.Colors[i%len(opts.Colors)]
			percentage := float64(values[i]) / float64(total) * 100

			// Квадратик цвета
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="12" height="12" fill="%s" rx="2"/>`,
				legendX, legendY+i*22, color))

			// Текст
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="10" fill="%s">%s</text>`,
				legendX+18, legendY+i*22+10, Colors.PrimaryDark, escapeText(label)))

			// Процент
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="9" fill="%s">%.0f%%</text>`,
				legendX+18, legendY+i*22+22, Colors.Gray, percentage))
		}
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

// generateDonutSegment генерирует сегмент donut диаграммы
func generateDonutSegment(cx, cy, outerR, innerR, startAngle, endAngle float64, color string) string {
	if endAngle-startAngle >= 360 {
		// Полный круг
		return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="%s" stroke-width="%.1f"/>`,
			cx, cy, (outerR+innerR)/2, color, outerR-innerR)
	}

	outerStart := polarToCartesian(cx, cy, outerR, startAngle)
	outerEnd := polarToCartesian(cx, cy, outerR, endAngle)
	innerStart := polarToCartesian(cx, cy, innerR, startAngle)
	innerEnd := polarToCartesian(cx, cy, innerR, endAngle)

	largeArc := 0
	if endAngle-startAngle > 180 {
		largeArc = 1
	}

	path := fmt.Sprintf("M %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 0 %.2f %.2f Z",
		outerStart.X, outerStart.Y,
		outerR, outerR, largeArc, outerEnd.X, outerEnd.Y,
		innerEnd.X, innerEnd.Y,
		innerR, innerR, largeArc, innerStart.X, innerStart.Y,
	)

	return fmt.Sprintf(`<path d="%s" fill="%s" stroke="white" stroke-width="2"/>`, path, color)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
