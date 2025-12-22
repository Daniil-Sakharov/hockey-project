package charts

import (
	"fmt"
	"strings"
)

// BarChartOptions настройки столбчатой диаграммы
type BarChartOptions struct {
	Size       Size
	Padding    Padding
	BarColor   string
	ShowValues bool
	ShowGrid   bool
}

// DefaultBarOptions возвращает настройки по умолчанию
func DefaultBarOptions() BarChartOptions {
	return BarChartOptions{
		Size:       DefaultSize(),
		Padding:    DefaultPadding(),
		BarColor:   Colors.Primary,
		ShowValues: true,
		ShowGrid:   true,
	}
}

// GenerateBarChart генерирует SVG столбчатую диаграмму
func GenerateBarChart(labels []string, values []int, opts *BarChartOptions) string {
	if len(labels) == 0 || len(values) == 0 {
		return generateEmptyChart("Нет данных")
	}

	if opts == nil {
		defaultOpts := DefaultBarOptions()
		opts = &defaultOpts
	}

	width := opts.Size.Width
	height := opts.Size.Height
	pad := opts.Padding

	chartWidth := width - pad.Left - pad.Right
	chartHeight := height - pad.Top - pad.Bottom

	maxVal := maxInt(values)
	if maxVal == 0 {
		maxVal = 1
	}

	barCount := len(values)
	barWidth := float64(chartWidth) / float64(barCount) * 0.7
	barGap := float64(chartWidth) / float64(barCount) * 0.3

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white" rx="8"/>`, width, height))

	if opts.ShowGrid {
		sb.WriteString(generateBarGrid(pad, chartWidth, chartHeight, maxVal))
	}

	for i, val := range values {
		barHeight := float64(val) / float64(maxVal) * float64(chartHeight)
		x := float64(pad.Left) + float64(i)*(barWidth+barGap) + barGap/2
		y := float64(pad.Top) + float64(chartHeight) - barHeight

		sb.WriteString(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" rx="4">`,
			x, y, barWidth, barHeight, opts.BarColor))
		sb.WriteString(fmt.Sprintf(`<animate attributeName="height" from="0" to="%.1f" dur="0.5s" fill="freeze"/>`, barHeight))
		sb.WriteString(fmt.Sprintf(`<animate attributeName="y" from="%.1f" to="%.1f" dur="0.5s" fill="freeze"/>`,
			float64(pad.Top)+float64(chartHeight), y))
		sb.WriteString(`</rect>`)

		if opts.ShowValues && val > 0 {
			textY := y - 5
			if textY < float64(pad.Top)+10 {
				textY = y + 15
			}
			sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" font-size="11" font-weight="600" fill="%s">%d</text>`,
				x+barWidth/2, textY, Colors.PrimaryDark, val))
		}

		if i < len(labels) {
			labelY := float64(height) - float64(pad.Bottom)/3
			sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.0f" text-anchor="middle" font-size="10" fill="%s">%s</text>`,
				x+barWidth/2, labelY, Colors.Gray, escapeText(labels[i])))
		}
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func generateBarGrid(pad Padding, chartWidth, chartHeight, maxVal int) string {
	var sb strings.Builder

	gridLines := 4
	for i := 0; i <= gridLines; i++ {
		y := float64(pad.Top) + float64(i)*float64(chartHeight)/float64(gridLines)
		val := maxVal - i*maxVal/gridLines

		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%.1f" x2="%d" y2="%.1f" stroke="#e5e7eb" stroke-width="1"/>`,
			pad.Left, y, pad.Left+chartWidth, y))
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%.1f" text-anchor="end" font-size="9" fill="%s">%d</text>`,
			pad.Left-5, y+3, Colors.Gray, val))
	}

	return sb.String()
}

func generateEmptyChart(message string) string {
	return fmt.Sprintf(`<svg width="300" height="180" viewBox="0 0 300 180" xmlns="http://www.w3.org/2000/svg">
		<rect width="300" height="180" fill="white" rx="8"/>
		<text x="150" y="90" text-anchor="middle" font-size="14" fill="%s">%s</text>
	</svg>`, Colors.Gray, escapeText(message))
}
