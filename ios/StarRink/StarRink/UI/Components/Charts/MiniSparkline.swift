import SwiftUI
import Charts

struct MiniSparkline: View {
    let data: [Int]
    let color: Color
    let label: String
    let currentValue: String

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.xs) {
            HStack {
                Text(currentValue)
                    .font(.srHeading3)
                    .foregroundColor(color)
                Spacer()
                // Trend indicator
                if let trend = calculateTrend() {
                    HStack(spacing: 2) {
                        Image(systemName: trend >= 0 ? "arrow.up.right" : "arrow.down.right")
                            .font(.caption2)
                        Text("\(abs(trend))%")
                            .font(.srCaption)
                    }
                    .foregroundColor(trend >= 0 ? .srSuccess : .srError)
                }
            }

            Chart {
                ForEach(Array(data.enumerated()), id: \.offset) { index, value in
                    LineMark(
                        x: .value("Index", index),
                        y: .value("Value", value)
                    )
                    .foregroundStyle(color)
                    .interpolationMethod(.catmullRom)

                    AreaMark(
                        x: .value("Index", index),
                        y: .value("Value", value)
                    )
                    .foregroundStyle(
                        LinearGradient(
                            colors: [color.opacity(0.4), color.opacity(0.0)],
                            startPoint: .top,
                            endPoint: .bottom
                        )
                    )
                    .interpolationMethod(.catmullRom)
                }
            }
            .chartXAxis(.hidden)
            .chartYAxis(.hidden)
            .frame(height: 40)

            Text(label)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity)
        .padding(AppSpacing.sm)
        .background(
            RoundedRectangle(cornerRadius: AppSpacing.radiusSmall)
                .fill(.ultraThinMaterial)
                .overlay(
                    RoundedRectangle(cornerRadius: AppSpacing.radiusSmall)
                        .stroke(Color.srBorder.opacity(0.3), lineWidth: 0.5)
                )
        )
    }

    private func calculateTrend() -> Int? {
        guard data.count >= 2 else { return nil }
        let recentCount = min(3, data.count / 2)
        let recent = data.suffix(recentCount).reduce(0, +)
        let previous = data.prefix(recentCount).reduce(0, +)
        guard previous > 0 else { return nil }
        return Int(((Double(recent) - Double(previous)) / Double(previous)) * 100)
    }
}

// Mini bar chart for points - enhanced version
struct MiniPointsChart: View {
    let data: [Int]
    let label: String

    private var maxValue: Int { max(data.max() ?? 1, 1) }
    private var total: Int { data.reduce(0, +) }
    private var average: Double { Double(total) / Double(max(data.count, 1)) }

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            // Header with stats
            HStack(alignment: .top) {
                VStack(alignment: .leading, spacing: 2) {
                    Text(label)
                        .font(.srBodyMedium)
                        .foregroundColor(.srTextPrimary)
                    Text("Последние \(data.count) игр")
                        .font(.caption2)
                        .foregroundColor(.srTextMuted)
                }
                Spacer()
                VStack(alignment: .trailing, spacing: 2) {
                    Text("\(total)")
                        .font(.srHeading3)
                        .foregroundStyle(
                            LinearGradient(
                                colors: [.srCyan, .srPurple],
                                startPoint: .leading,
                                endPoint: .trailing
                            )
                        )
                    Text("Ср: \(String(format: "%.1f", average))")
                        .font(.caption2)
                        .foregroundColor(.srTextMuted)
                }
            }

            // Custom bar chart with labels
            HStack(alignment: .bottom, spacing: 6) {
                ForEach(Array(data.enumerated()), id: \.offset) { index, value in
                    VStack(spacing: 4) {
                        // Value label
                        Text("\(value)")
                            .font(.system(size: 10, weight: .semibold))
                            .foregroundColor(value == maxValue ? .srCyan : .srTextSecondary)

                        // Bar
                        RoundedRectangle(cornerRadius: 4)
                            .fill(
                                LinearGradient(
                                    colors: barColors(for: value),
                                    startPoint: .bottom,
                                    endPoint: .top
                                )
                            )
                            .frame(height: barHeight(for: value))
                            .overlay(
                                RoundedRectangle(cornerRadius: 4)
                                    .stroke(
                                        value == maxValue ? Color.srCyan.opacity(0.5) : Color.clear,
                                        lineWidth: 1
                                    )
                            )

                        // Game number
                        Text("И\(index + 1)")
                            .font(.system(size: 8))
                            .foregroundColor(.srTextMuted)
                    }
                    .frame(maxWidth: .infinity)
                }
            }
            .frame(height: 80)
        }
        .glassCard(padding: AppSpacing.md)
    }

    private func barHeight(for value: Int) -> CGFloat {
        let minHeight: CGFloat = 4
        let maxHeight: CGFloat = 50
        guard maxValue > 0 else { return minHeight }
        return minHeight + (maxHeight - minHeight) * CGFloat(value) / CGFloat(maxValue)
    }

    private func barColors(for value: Int) -> [Color] {
        if value == maxValue {
            return [Color.srCyan, Color.srPurple]
        } else if value == 0 {
            return [Color.srBorder.opacity(0.3), Color.srBorder.opacity(0.5)]
        } else {
            return [Color.srCyan.opacity(0.5), Color.srPurple.opacity(0.5)]
        }
    }
}

#Preview {
    VStack(spacing: 16) {
        MiniSparkline(
            data: [1, 0, 2, 1, 0, 3],
            color: .srCyan,
            label: "Голы",
            currentValue: "7"
        )
        MiniPointsChart(
            data: [3, 1, 2, 4, 0, 4],
            label: "Очки за сезон"
        )
    }
    .padding()
    .background(Color.srBackground)
}
