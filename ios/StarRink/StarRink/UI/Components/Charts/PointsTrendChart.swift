import SwiftUI
import Charts

struct PointsTrendChart: View {
    let data: [SeasonTrend]

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text("Тренд очков")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            Chart {
                ForEach(data) { trend in
                    // Goals line
                    LineMark(
                        x: .value("Месяц", trend.month),
                        y: .value("Голы", trend.goals)
                    )
                    .foregroundStyle(Color.srCyan)
                    .symbol(.circle)
                    .interpolationMethod(.catmullRom)

                    // Area under goals
                    AreaMark(
                        x: .value("Месяц", trend.month),
                        y: .value("Голы", trend.goals)
                    )
                    .foregroundStyle(
                        LinearGradient(
                            colors: [Color.srCyan.opacity(0.3), Color.srCyan.opacity(0.0)],
                            startPoint: .top,
                            endPoint: .bottom
                        )
                    )
                    .interpolationMethod(.catmullRom)

                    // Assists line
                    LineMark(
                        x: .value("Месяц", trend.month),
                        y: .value("Передачи", trend.assists)
                    )
                    .foregroundStyle(Color.srPurple)
                    .symbol(.diamond)
                    .interpolationMethod(.catmullRom)
                }
            }
            .chartXAxis {
                AxisMarks(values: .automatic) { _ in
                    AxisGridLine(stroke: StrokeStyle(lineWidth: 0.5, dash: [4]))
                        .foregroundStyle(Color.srBorder)
                    AxisValueLabel()
                        .foregroundStyle(Color.srTextSecondary)
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
                    Label("Голы", systemImage: "circle.fill")
                        .font(.srCaption)
                        .foregroundColor(.srCyan)
                    Label("Передачи", systemImage: "diamond.fill")
                        .font(.srCaption)
                        .foregroundColor(.srPurple)
                }
            }
            .frame(height: 200)
        }
        .glassCard()
    }
}

#Preview {
    PointsTrendChart(data: SeasonTrend.sampleData)
        .padding()
        .background(Color.srBackground)
}
