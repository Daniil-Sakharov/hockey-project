package charts

import (
	"fmt"
	"strings"
)

// LineDataset набор данных для линейного графика
type LineDataset struct {
	Label  string
	Values []int
	Color  string
}

// LineChartOptions настройки линейного графика
type LineChartOptions struct {
	Size       Size
	Padding    Padding
	ShowGrid   bool
	ShowPoints bool
	ShowArea   bool
	Smooth     bool
}

// DefaultLineOptions возвращает настройки по умолчанию
func DefaultLineOptions() LineChartOptions {
	return LineChartOptions{
		Size:       DefaultSize(),
		Padding:    Padding{Top: 20, Right: 20, Bottom: 35, Left: 40},
		ShowGrid:   true,
		ShowPoints: true,
		ShowArea:   true,
		Smooth:     false,
	}
}

// GenerateLineChart генерирует SVG линейный график
func GenerateLineChart(labels []string, datasets []LineDataset, opts *LineChartOptions) string {
	if len(labels) == 0 || len(datasets) == 0 {
		return generateEmptyChart("Нет данных")
	}

	if opts == nil {
		defaultOpts := DefaultLineOptions()
		opts = &defaultOpts
	}

	width := opts.Size.Width
	height := opts.Size.Height
	pad := opts.Padding

	chartWidth := width - pad.Left - pad.Right
	chartHeight := height - pad.Top - pad.Bottom

	maxVal := 0
	for _, ds := range datasets {
		if m := maxInt(ds.Values); m > maxVal {
			maxVal = m
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white" rx="8"/>`, width, height))

	if opts.ShowGrid {
		sb.WriteString(generateLineGrid(pad, chartWidth, chartHeight, maxVal))
	}

	for _, ds := range datasets {
		if len(ds.Values) == 0 {
			continue
		}

		points := calculateLinePoints(ds.Values, pad, chartWidth, chartHeight, maxVal)
		color := ds.Color
		if color == "" {
			color = Colors.Primary
		}

		if opts.ShowArea {
			areaPath := buildAreaPath(points, pad.Top+chartHeight)
			sb.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" fill-opacity="0.1"/>`, areaPath, color))
		}

		linePath := buildLinePath(points)
		sb.WriteString(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>`,
			linePath, color))

		if opts.ShowPoints {
			for _, p := range points {
				sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="4" fill="white" stroke="%s" stroke-width="2"/>`,
					p.X, p.Y, color))
			}
		}
	}

	pointCount := len(labels)
	if pointCount > 0 {
		step := float64(chartWidth) / float64(pointCount-1)
		if pointCount == 1 {
			step = 0
		}
		for i, label := range labels {
			x := float64(pad.Left) + float64(i)*step
			y := float64(height) - float64(pad.Bottom)/3
			sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.0f" text-anchor="middle" font-size="9" fill="%s">%s</text>`,
				x, y, Colors.Gray, escapeText(label)))
		}
	}

	if len(datasets) > 1 {
		legendY := height - 8
		legendX := pad.Left
		for i, ds := range datasets {
			color := ds.Color
			if color == "" {
				color = ChartPalette[i%len(ChartPalette)]
			}
			sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2"/>`,
				legendX, legendY, legendX+15, legendY, color))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="9" fill="%s">%s</text>`,
				legendX+20, legendY+3, Colors.Gray, escapeText(ds.Label)))
			legendX += 70
		}
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func calculateLinePoints(values []int, pad Padding, chartWidth, chartHeight, maxVal int) []Point {
	points := make([]Point, len(values))
	pointCount := len(values)

	step := float64(chartWidth) / float64(pointCount-1)
	if pointCount == 1 {
		step = 0
	}

	for i, val := range values {
		x := float64(pad.Left) + float64(i)*step
		y := float64(pad.Top) + float64(chartHeight) - (float64(val)/float64(maxVal))*float64(chartHeight)
		points[i] = Point{X: x, Y: y}
	}

	return points
}

func buildLinePath(points []Point) string {
	if len(points) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("M %.1f %.1f", points[0].X, points[0].Y))

	for i := 1; i < len(points); i++ {
		sb.WriteString(fmt.Sprintf(" L %.1f %.1f", points[i].X, points[i].Y))
	}

	return sb.String()
}

func buildAreaPath(points []Point, baseY int) string {
	if len(points) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("M %.1f %d", points[0].X, baseY))

	for _, p := range points {
		sb.WriteString(fmt.Sprintf(" L %.1f %.1f", p.X, p.Y))
	}

	sb.WriteString(fmt.Sprintf(" L %.1f %d Z", points[len(points)-1].X, baseY))

	return sb.String()
}

func generateLineGrid(pad Padding, chartWidth, chartHeight, maxVal int) string {
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
