import SwiftUI
import Charts

struct PlusMinusChart: View {
    let data: [SeasonAggregated]

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text("+/- по сезонам")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)

            Chart(data) { item in
                BarMark(
                    x: .value("Сезон", shortSeason(item.season)),
                    y: .value("+/-", item.plusMinus)
                )
                .foregroundStyle(item.plusMinus >= 0 ? Color.srSuccess : Color.srError)
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
