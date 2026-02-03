import SwiftUI

struct PlayerChartsSection: View {
    let seasonData: [SeasonAggregated]
    @State private var isExpanded = false

    var body: some View {
        VStack(spacing: AppSpacing.md) {
            toggleButton
            if isExpanded {
                chartsGrid
                    .transition(.opacity.combined(with: .move(edge: .top)))
            }
        }
    }

    private var toggleButton: some View {
        Button {
            withAnimation(.spring(response: 0.3, dampingFraction: 0.8)) {
                isExpanded.toggle()
            }
        } label: {
            HStack {
                Image(systemName: "chart.line.uptrend.xyaxis")
                    .foregroundColor(.srCyan)
                Text("Аналитика")
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                Spacer()
                Image(systemName: "chevron.down")
                    .foregroundColor(.srTextMuted)
                    .rotationEffect(.degrees(isExpanded ? 180 : 0))
                    .animation(.spring(response: 0.3), value: isExpanded)
            }
            .padding(.horizontal, AppSpacing.md)
            .padding(.vertical, AppSpacing.sm)
        }
        .glassCard(padding: 0)
    }

    private var chartsGrid: some View {
        LazyVGrid(
            columns: [GridItem(.flexible()), GridItem(.flexible())],
            spacing: AppSpacing.sm
        ) {
            SeasonProgressChart(data: seasonData)
            PointsBreakdownChart(data: seasonData)
            AvgPerGameChart(data: seasonData)
            PenaltyChart(data: seasonData)
            PlusMinusChart(data: seasonData)
            RadarChartView(data: seasonData)
        }
    }
}
