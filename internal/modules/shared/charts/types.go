package charts

import (
	"fmt"
	"math"
)

// Colors - стандартная цветовая палитра
var Colors = struct {
	Primary     string
	PrimaryDark string
	Accent      string
	AccentLight string
	Success     string
	Warning     string
	Gray        string
	White       string
}{
	Primary:     "#4a90d9",
	PrimaryDark: "#1a3a5c",
	Accent:      "#7bb8e8",
	AccentLight: "#a8d4f0",
	Success:     "#10b981",
	Warning:     "#f59e0b",
	Gray:        "#6b7280",
	White:       "#ffffff",
}

// ChartPalette - палитра для графиков
var ChartPalette = []string{
	"#4a90d9", // Primary blue
	"#7bb8e8", // Light blue
	"#2c5aa0", // Dark blue
	"#10b981", // Green
	"#f59e0b", // Orange
	"#ef4444", // Red
}

// Size содержит размеры графика
type Size struct {
	Width  int
	Height int
}

// Padding содержит отступы
type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// DefaultSize возвращает размер по умолчанию
func DefaultSize() Size {
	return Size{Width: 300, Height: 180}
}

// DefaultPadding возвращает отступы по умолчанию
func DefaultPadding() Padding {
	return Padding{Top: 20, Right: 20, Bottom: 30, Left: 40}
}

// Point представляет точку на графике
type Point struct {
	X, Y float64
}

// degToRad конвертирует градусы в радианы
func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

// polarToCartesian конвертирует полярные координаты в декартовы
func polarToCartesian(cx, cy, r, angleDeg float64) Point {
	rad := degToRad(angleDeg - 90) // -90 чтобы начинать сверху
	return Point{
		X: cx + r*math.Cos(rad),
		Y: cy + r*math.Sin(rad),
	}
}

// describeArc создаёт SVG path для дуги
func describeArc(cx, cy, r, startAngle, endAngle float64) string {
	start := polarToCartesian(cx, cy, r, endAngle)
	end := polarToCartesian(cx, cy, r, startAngle)

	largeArcFlag := 0
	if endAngle-startAngle > 180 {
		largeArcFlag = 1
	}

	return fmt.Sprintf("M %.2f %.2f A %.2f %.2f 0 %d 0 %.2f %.2f L %.2f %.2f Z",
		start.X, start.Y,
		r, r,
		largeArcFlag,
		end.X, end.Y,
		cx, cy,
	)
}

// maxInt находит максимум в слайсе int
func maxInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// maxFloat находит максимум в слайсе float64
func maxFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// sumInt суммирует слайс int
func sumInt(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}

// escapeText экранирует текст для SVG
func escapeText(s string) string {
	result := ""
	for _, r := range s {
		switch r {
		case '<':
			result += "&lt;"
		case '>':
			result += "&gt;"
		case '&':
			result += "&amp;"
		case '"':
			result += "&quot;"
		default:
			result += string(r)
		}
	}
	return result
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
