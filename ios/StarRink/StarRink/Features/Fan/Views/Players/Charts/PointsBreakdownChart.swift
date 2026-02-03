import SwiftUI
import Charts

struct PointsBreakdownChart: View {
    let data: [SeasonAggregated]

    private var totalGoals: Int { data.reduce(0) { $0 + $1.goals } }
    private var totalAssists: Int { data.reduce(0) { $0 + $1.assists } }

    private var slices: [(label: String, value: Int, color: Color)] {
        [
            ("Голы", totalGoals, .srCyan),
            ("Передачи", totalAssists, .srPurple)
        ]
    }

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            Text("Голы / Передачи")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)

            Chart(slices, id: \.label) { slice in
                SectorMark(
                    angle: .value(slice.label, slice.value),
                    innerRadius: .ratio(0.5),
                    angularInset: 2
                )
                .foregroundStyle(slice.color)
                .cornerRadius(4)
            }
            .frame(height: 120)

            HStack(spacing: AppSpacing.sm) {
                LegendDot(color: .srCyan, label: "\(totalGoals)Г")
                LegendDot(color: .srPurple, label: "\(totalAssists)П")
            }
        }
        .glassCard(padding: AppSpacing.sm)
    }
}
