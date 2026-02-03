import SwiftUI
import Charts

struct StatsView: View {
    @StateObject private var viewModel = StatsViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                // Header with time range picker
                headerSection

                // Summary cards
                summaryCardsSection

                // Points trend line chart
                if !viewModel.seasonTrends.isEmpty {
                    PointsTrendChart(data: viewModel.seasonTrends)
                }

                // Game stats bar chart
                if !viewModel.gameStats.isEmpty {
                    GameStatsBarChart(data: viewModel.gameStats)
                }

                // Detailed stats list
                if !viewModel.gameStats.isEmpty {
                    detailedStatsList
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.vertical, AppSpacing.md)
        }
        .background(Color.clear)
        .refreshable {
            await viewModel.refresh()
        }
        .task {
            await viewModel.loadStats()
        }
    }

    // MARK: - Sections

    private var headerSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            Text("Статистика")
                .font(.srHeading2)
                .foregroundColor(.srTextPrimary)

            // Time range picker
            Picker("Период", selection: $viewModel.selectedTimeRange) {
                ForEach(StatsViewModel.TimeRange.allCases, id: \.self) { range in
                    Text(range.rawValue).tag(range)
                }
            }
            .pickerStyle(.segmented)
        }
    }

    private var summaryCardsSection: some View {
        HStack(spacing: AppSpacing.sm) {
            SummaryStatCard(
                value: "\(viewModel.totalGoals)",
                label: "Голы",
                icon: "sportscourt.fill",
                color: .srCyan
            )
            SummaryStatCard(
                value: "\(viewModel.totalAssists)",
                label: "Передачи",
                icon: "arrow.triangle.branch",
                color: .srPurple
            )
            SummaryStatCard(
                value: String(format: "%.1f", viewModel.averagePointsPerGame),
                label: "Очков/игра",
                icon: "chart.line.uptrend.xyaxis",
                color: .srAmber
            )
        }
    }

    private var detailedStatsList: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            HStack {
                Text("Последние игры")
                    .font(.srHeading4)
                    .foregroundColor(.srTextPrimary)
                Spacer()
                Text("\(viewModel.gamesPlayed) игр")
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
            }

            ForEach(viewModel.gameStats.prefix(5)) { game in
                GameStatRow(game: game)
            }

            // Plus/Minus summary
            HStack {
                Text("+/-")
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextSecondary)
                Spacer()
                Text(viewModel.plusMinus >= 0 ? "+\(viewModel.plusMinus)" : "\(viewModel.plusMinus)")
                    .font(.srHeading4)
                    .foregroundColor(viewModel.plusMinus >= 0 ? .srSuccess : .srError)
            }
            .padding(.top, AppSpacing.sm)
        }
        .glassCard()
    }
}

// MARK: - Supporting Views

struct SummaryStatCard: View {
    let value: String
    let label: String
    let icon: String
    let color: Color

    var body: some View {
        VStack(spacing: AppSpacing.xs) {
            Image(systemName: icon)
                .font(.title3)
                .foregroundColor(color)
            Text(value)
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
            Text(label)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity)
        .glassCard(padding: AppSpacing.sm)
    }
}

struct GameStatRow: View {
    let game: GameStat

    var body: some View {
        HStack {
            VStack(alignment: .leading, spacing: 2) {
                Text("Игра \(game.gameNumber)")
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                Text(game.formattedDate)
                    .font(.srCaption)
                    .foregroundColor(.srTextMuted)
            }

            Spacer()

            HStack(spacing: AppSpacing.md) {
                StatPill(value: game.goals, label: "Г", color: .srCyan)
                StatPill(value: game.assists, label: "П", color: .srPurple)
                StatPill(value: game.points, label: "О", color: .srAmber)
            }
        }
        .padding(.vertical, AppSpacing.xs)
    }
}

struct StatPill: View {
    let value: Int
    let label: String
    let color: Color

    var body: some View {
        HStack(spacing: 4) {
            Text("\(value)")
                .font(.srBodyMedium)
                .foregroundColor(.srTextPrimary)
            Text(label)
                .font(.srCaption)
                .foregroundColor(color)
        }
    }
}

#Preview {
    StatsView()
        .environmentObject(AuthViewModel())
}
