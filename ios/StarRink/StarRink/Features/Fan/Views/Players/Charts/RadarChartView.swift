import SwiftUI

struct RadarChartView: View {
    let data: [SeasonAggregated]

    private let axes = ["Голы", "Передачи", "+/-", "Игры", "Штраф"]
    private let gridLevels = 3

    private var current: SeasonAggregated? { data.last }
    private var previous: SeasonAggregated? { data.count >= 2 ? data[data.count - 2] : nil }

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text("Радар")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)

            GeometryReader { geo in
                let size = min(geo.size.width, geo.size.height)
                let center = CGPoint(x: geo.size.width / 2, y: size / 2)
                let radius = size * 0.38

                ZStack {
                    gridShape(center: center, radius: radius)
                    axisLines(center: center, radius: radius)

                    if let cur = current {
                        dataPolygon(values: normalize(cur), center: center, radius: radius, color: .srCyan)
                    }
                    if let prev = previous {
                        dataPolygon(values: normalize(prev), center: center, radius: radius, color: .srPurple)
                    }

                    axisLabels(center: center, radius: radius)
                }
            }
            .frame(height: 120)

            HStack(spacing: AppSpacing.sm) {
                LegendDot(color: .srCyan, label: "Текущий")
                if previous != nil {
                    LegendDot(color: .srPurple, label: "Прошлый")
                }
            }
        }
        .glassCard(padding: AppSpacing.sm)
    }

    // MARK: - Grid

    private func gridShape(center: CGPoint, radius: CGFloat) -> some View {
        ForEach(1...gridLevels, id: \.self) { level in
            let r = radius * CGFloat(level) / CGFloat(gridLevels)
            PolygonPath(sides: axes.count, radius: r, center: center)
                .stroke(Color.srBorder.opacity(0.3), lineWidth: 0.5)
        }
    }

    private func axisLines(center: CGPoint, radius: CGFloat) -> some View {
        ForEach(0..<axes.count, id: \.self) { i in
            let angle = angleFor(index: i)
            Path { path in
                path.move(to: center)
                path.addLine(to: point(center: center, radius: radius, angle: angle))
            }
            .stroke(Color.srBorder.opacity(0.2), lineWidth: 0.5)
        }
    }

    private func axisLabels(center: CGPoint, radius: CGFloat) -> some View {
        ForEach(0..<axes.count, id: \.self) { i in
            let angle = angleFor(index: i)
            let p = point(center: center, radius: radius + 14, angle: angle)
            Text(axes[i])
                .font(.system(size: 7))
                .foregroundColor(.srTextMuted)
                .position(p)
        }
    }

    // MARK: - Data Polygon

    private func dataPolygon(values: [Double], center: CGPoint, radius: CGFloat, color: Color) -> some View {
        let path = Path { path in
            for (i, v) in values.enumerated() {
                let angle = angleFor(index: i)
                let r = radius * CGFloat(v)
                let p = point(center: center, radius: r, angle: angle)
                if i == 0 { path.move(to: p) } else { path.addLine(to: p) }
            }
            path.closeSubpath()
        }

        return ZStack {
            path.fill(color.opacity(0.15))
            path.stroke(color, lineWidth: 1.5)
        }
    }

    // MARK: - Helpers

    private func normalize(_ s: SeasonAggregated) -> [Double] {
        let maxVals = maxValues()
        return [
            safe(Double(s.goals), max: maxVals.0),
            safe(Double(s.assists), max: maxVals.1),
            safe(Double(s.plusMinus), max: maxVals.2),
            safe(Double(s.games), max: maxVals.3),
            safe(Double(s.penaltyMinutes), max: maxVals.4)
        ]
    }

    private func maxValues() -> (Double, Double, Double, Double, Double) {
        let g = data.map(\.goals).max().map(Double.init) ?? 1
        let a = data.map(\.assists).max().map(Double.init) ?? 1
        let pm = data.map { abs($0.plusMinus) }.max().map(Double.init) ?? 1
        let gm = data.map(\.games).max().map(Double.init) ?? 1
        let pen = data.map(\.penaltyMinutes).max().map(Double.init) ?? 1
        return (max(g, 1), max(a, 1), max(pm, 1), max(gm, 1), max(pen, 1))
    }

    private func safe(_ val: Double, max: Double) -> Double {
        guard max > 0 else { return 0 }
        return min(abs(val) / max, 1.0)
    }

    private func angleFor(index: Int) -> Double {
        let slice = 2 * .pi / Double(axes.count)
        return slice * Double(index) - .pi / 2
    }

    private func point(center: CGPoint, radius: CGFloat, angle: Double) -> CGPoint {
        CGPoint(
            x: center.x + radius * CGFloat(cos(angle)),
            y: center.y + radius * CGFloat(sin(angle))
        )
    }
}

// MARK: - Polygon Path

private struct PolygonPath: Shape {
    let sides: Int
    let radius: CGFloat
    let center: CGPoint

    func path(in rect: CGRect) -> Path {
        Path { path in
            for i in 0..<sides {
                let angle = 2 * .pi / Double(sides) * Double(i) - .pi / 2
                let p = CGPoint(
                    x: center.x + radius * CGFloat(cos(angle)),
                    y: center.y + radius * CGFloat(sin(angle))
                )
                if i == 0 { path.move(to: p) } else { path.addLine(to: p) }
            }
            path.closeSubpath()
        }
    }
}
