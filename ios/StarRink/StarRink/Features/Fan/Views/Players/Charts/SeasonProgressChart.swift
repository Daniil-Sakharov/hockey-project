import SwiftUI
import Charts

struct SeasonProgressChart: View {
    let data: [SeasonAggregated]

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text("Прогресс")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)

            Chart {
                ForEach(data) { item in
                    AreaMark(
                        x: .value("Сезон", shortSeason(item.season)),
                        y: .value("Очки", item.points)
                    )
                    .foregroundStyle(Color.srAmber.opacity(0.2))

                    LineMark(
                        x: .value("Сезон", shortSeason(item.season)),
                        y: .value("Очки", item.points)
                    )
                    .foregroundStyle(Color.srAmber)
                    .symbol(Circle())

                    LineMark(
                        x: .value("Сезон", shortSeason(item.season)),
                        y: .value("Голы", item.goals)
                    )
                    .foregroundStyle(Color.srCyan)
                    .symbol(Circle())
                }
            }
            .chartXAxis {
                AxisMarks { _ in
                    AxisValueLabel()
                        .font(.system(size: 8))
                        .foregroundStyle(Color.srTextMuted)
                }
            }
            .chartYAxis {
                AxisMarks { _ in
                    AxisValueLabel()
                        .font(.system(size: 8))
                        .foregroundStyle(Color.srTextMuted)
                }
            }
            .frame(height: 120)

            legendRow
        }
        .glassCard(padding: AppSpacing.sm)
    }

    private var legendRow: some View {
        HStack(spacing: AppSpacing.sm) {
            LegendDot(color: .srAmber, label: "О")
            LegendDot(color: .srCyan, label: "Г")
        }
    }

    private func shortSeason(_ s: String) -> String {
        if s.count > 5 { return String(s.suffix(5)) }
        return s
    }
}

struct LegendDot: View {
    let color: Color
    let label: String

    var body: some View {
        HStack(spacing: 3) {
            Circle().fill(color).frame(width: 6, height: 6)
            Text(label)
                .font(.system(size: 9))
                .foregroundColor(.srTextMuted)
        }
    }
}
