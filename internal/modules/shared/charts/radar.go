package charts

import (
	"fmt"
	"strings"
)

// RadarChartOptions настройки radar диаграммы
type RadarChartOptions struct {
	Size       Size
	Color      string
	FillColor  string
	ShowGrid   bool
	GridLevels int
}

// DefaultRadarOptions возвращает настройки по умолчанию
func DefaultRadarOptions() RadarChartOptions {
	return RadarChartOptions{
		Size:       Size{Width: 300, Height: 200},
		Color:      Colors.Primary,
		FillColor:  Colors.Primary,
		ShowGrid:   true,
		GridLevels: 4,
	}
}

// GenerateRadarChart генерирует SVG radar диаграмму
func GenerateRadarChart(labels []string, values []float64, opts *RadarChartOptions) string {
	if len(labels) == 0 || len(values) == 0 {
		return generateEmptyChart("Нет данных")
	}

	if len(labels) != len(values) {
		return generateEmptyChart("Ошибка данных")
	}

	if opts == nil {
		defaultOpts := DefaultRadarOptions()
		opts = &defaultOpts
	}

	width := opts.Size.Width
	height := opts.Size.Height

	cx := float64(width) / 2
	cy := float64(height) / 2
	radius := minFloat(cx, cy) - 35

	maxVal := maxFloat(values)
	if maxVal == 0 {
		maxVal = 1
	}

	numPoints := len(labels)
	angleStep := 360.0 / float64(numPoints)

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white" rx="8"/>`, width, height))

	if opts.ShowGrid {
		sb.WriteString(generateRadarGrid(cx, cy, radius, numPoints, opts.GridLevels))
	}

	for i, label := range labels {
		angle := float64(i)*angleStep - 90
		endPoint := polarToCartesian(cx, cy, radius, angle+90)

		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#d1d5db" stroke-width="1"/>`,
			cx, cy, endPoint.X, endPoint.Y))

		labelPoint := polarToCartesian(cx, cy, radius+15, angle+90)
		anchor := "middle"
		if labelPoint.X < cx-10 {
			anchor = "end"
		} else if labelPoint.X > cx+10 {
			anchor = "start"
		}

		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="%s" font-size="10" fill="%s">%s</text>`,
			labelPoint.X, labelPoint.Y+3, anchor, Colors.Gray, escapeText(label)))
	}

	dataPoints := make([]Point, numPoints)
	for i, val := range values {
		angle := float64(i)*angleStep - 90
		normalizedVal := val / maxVal
		r := radius * normalizedVal
		dataPoints[i] = polarToCartesian(cx, cy, r, angle+90)
	}

	polygonPath := buildPolygonPath(dataPoints)
	sb.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" fill-opacity="0.2" stroke="%s" stroke-width="2"/>`,
		polygonPath, opts.FillColor, opts.Color))

	for _, p := range dataPoints {
		sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="4" fill="white" stroke="%s" stroke-width="2"/>`,
			p.X, p.Y, opts.Color))
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

func generateRadarGrid(cx, cy, radius float64, numPoints, levels int) string {
	var sb strings.Builder

	angleStep := 360.0 / float64(numPoints)

	for level := 1; level <= levels; level++ {
		r := radius * float64(level) / float64(levels)
		points := make([]Point, numPoints)

		for i := 0; i < numPoints; i++ {
			angle := float64(i)*angleStep - 90
			points[i] = polarToCartesian(cx, cy, r, angle+90)
		}

		polygonPath := buildPolygonPath(points)
		sb.WriteString(fmt.Sprintf(`<path d="%s" fill="none" stroke="#e5e7eb" stroke-width="1"/>`, polygonPath))
	}

	return sb.String()
}

func buildPolygonPath(points []Point) string {
	if len(points) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("M %.1f %.1f", points[0].X, points[0].Y))

	for i := 1; i < len(points); i++ {
		sb.WriteString(fmt.Sprintf(" L %.1f %.1f", points[i].X, points[i].Y))
	}

	sb.WriteString(" Z")
	return sb.String()
}
