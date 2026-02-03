import SwiftUI
import Charts

struct GameStatsBarChart: View {
    let data: [GameStat]

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text("Статистика по играм")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            Chart(data) { game in
                BarMark(
                    x: .value("Игра", "И\(game.gameNumber)"),
                    y: .value("Очки", game.points)
                )
                .foregroundStyle(
                    LinearGradient(
                        colors: [Color.srCyan, Color.srPurple],
                        startPoint: .bottom,
                        endPoint: .top
                    )
                )
                .cornerRadius(4)

                // Goals points on top
                PointMark(
                    x: .value("Игра", "И\(game.gameNumber)"),
                    y: .value("Голы", game.goals)
                )
                .foregroundStyle(Color.srAmber)
                .symbolSize(50)
            }
            .chartXAxis {
                AxisMarks(values: .automatic) { _ in
                    AxisValueLabel()
                        .foregroundStyle(Color.srTextSecondary)
                        .font(.system(size: 10))
                }
            }
            .chartYAxis {
                AxisMarks(position: .leading) { _ in
                    AxisGridLine(stroke: StrokeStyle(lineWidth: 0.5, dash: [4]))
                        .foregroundStyle(Color.srBorder)
                    AxisValueLabel()
                        .foregroundStyle(Color.srTextSecondary)
                }
            }
            .chartLegend(position: .bottom, alignment: .center) {
                HStack(spacing: AppSpacing.lg) {
                    Label("Очки", systemImage: "chart.bar.fill")
                        .font(.srCaption)
                        .foregroundColor(.srCyan)
                    Label("Голы", systemImage: "circle.fill")
                        .font(.srCaption)
                        .foregroundColor(.srAmber)
                }
            }
            .frame(height: 180)
        }
        .glassCard()
    }
}

#Preview {
    GameStatsBarChart(data: GameStat.sampleData)
        .padding()
        .background(Color.srBackground)
}
