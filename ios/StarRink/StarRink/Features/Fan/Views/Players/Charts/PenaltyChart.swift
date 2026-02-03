import SwiftUI
import Charts

struct PenaltyChart: View {
    let data: [SeasonAggregated]

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text("Штрафные минуты")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)

            Chart(data) { item in
                BarMark(
                    x: .value("Сезон", shortSeason(item.season)),
                    y: .value("Мин", item.penaltyMinutes)
                )
                .foregroundStyle(Color.srError.opacity(0.7))
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
        }
        .glassCard(padding: AppSpacing.sm)
    }

    private func shortSeason(_ s: String) -> String {
        if s.count > 5 { return String(s.suffix(5)) }
        return s
    }
}
