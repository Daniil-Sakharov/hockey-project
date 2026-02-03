import SwiftUI
import Charts

struct AvgPerGameChart: View {
    let data: [SeasonAggregated]

    private struct BarEntry: Identifiable {
        let id = UUID()
        let season: String
        let metric: String
        let value: Double
        let color: Color
    }

    private var entries: [BarEntry] {
        data.flatMap { item in
            let s = shortSeason(item.season)
            return [
                BarEntry(season: s, metric: "Г", value: item.avgGoals, color: .srCyan),
                BarEntry(season: s, metric: "П", value: item.avgAssists, color: .srPurple),
                BarEntry(season: s, metric: "О", value: item.avgPoints, color: .srAmber)
            ]
        }
    }

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text("Средние за игру")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)

            Chart(entries) { entry in
                BarMark(
                    x: .value("Сезон", entry.season),
                    y: .value("Среднее", entry.value)
                )
                .foregroundStyle(entry.color)
                .position(by: .value("Метрика", entry.metric))
                .cornerRadius(3)
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

            HStack(spacing: AppSpacing.sm) {
                LegendDot(color: .srCyan, label: "Г")
                LegendDot(color: .srPurple, label: "П")
                LegendDot(color: .srAmber, label: "О")
            }
        }
        .glassCard(padding: AppSpacing.sm)
    }

    private func shortSeason(_ s: String) -> String {
        if s.count > 5 { return String(s.suffix(5)) }
        return s
    }
}
