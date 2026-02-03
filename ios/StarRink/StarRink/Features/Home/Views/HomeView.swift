import SwiftUI
import Charts

struct HomeView: View {
    @EnvironmentObject var authViewModel: AuthViewModel

    // Sample data for sparklines
    private let goalsData = [1, 0, 2, 1, 0, 3, 1]
    private let assistsData = [2, 1, 0, 3, 0, 1, 2]
    private let pointsData = [3, 1, 2, 4, 0, 4, 3]

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                welcomeSection
                quickStatsWithChartsSection
                recentActivitySection
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100) // Extra space for TabBar
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
    }

    // MARK: - Sections

    private var welcomeSection: some View {
        HStack {
            VStack(alignment: .leading, spacing: AppSpacing.xxs) {
                Text("Привет, \(authViewModel.currentUser?.displayName ?? "Игрок")!")
                    .font(.srHeading3)
                    .foregroundColor(.srTextPrimary)
                Text("Сезон 2024/25")
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
            }
            Spacer()
            ZStack {
                Circle()
                    .fill(
                        LinearGradient(
                            colors: [Color.srCyan, Color.srPurple],
                            startPoint: .topLeading,
                            endPoint: .bottomTrailing
                        )
                    )
                    .frame(width: 50, height: 50)
                Text(authViewModel.currentUser?.displayName.prefix(1).uppercased() ?? "?")
                    .font(.srHeading4)
                    .foregroundColor(.white)
            }
        }
        .glassCard()
    }

    private var quickStatsWithChartsSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            HStack {
                Text("Статистика")
                    .font(.srHeading4)
                    .foregroundColor(.srTextPrimary)
                Spacer()
                NavigationLink(destination: StatsView()) {
                    HStack(spacing: 4) {
                        Text("Подробнее")
                            .font(.srCaption)
                        Image(systemName: "chevron.right")
                            .font(.caption2)
                    }
                    .foregroundColor(.srCyan)
                }
            }

            // Mini sparkline cards
            HStack(spacing: AppSpacing.sm) {
                MiniSparkline(
                    data: goalsData,
                    color: .srCyan,
                    label: "Голы",
                    currentValue: "\(goalsData.reduce(0, +))"
                )
                MiniSparkline(
                    data: assistsData,
                    color: .srPurple,
                    label: "Передачи",
                    currentValue: "\(assistsData.reduce(0, +))"
                )
            }

            // Points trend mini chart
            MiniPointsChart(data: pointsData, label: "Очки за сезон")
        }
    }

    private var recentActivitySection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            Text("Последняя активность")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            VStack(spacing: 0) {
                ActivityRow(
                    icon: "sportscourt.fill",
                    title: "Игра vs Метеор",
                    subtitle: "2 гола, 1 передача",
                    time: "Вчера",
                    color: .srCyan
                )
                Divider().background(Color.srBorder.opacity(0.3))
                ActivityRow(
                    icon: "figure.hockey",
                    title: "Тренировка",
                    subtitle: "90 минут",
                    time: "2 дня назад",
                    color: .srPurple
                )
                Divider().background(Color.srBorder.opacity(0.3))
                ActivityRow(
                    icon: "trophy.fill",
                    title: "Достижение",
                    subtitle: "10 голов за сезон",
                    time: "3 дня назад",
                    color: .srAmber
                )
            }
            .glassCard(padding: 0)
        }
    }
}

// MARK: - Supporting Views

struct StatCard: View {
    let value: String
    let label: String
    let color: Color

    var body: some View {
        VStack(spacing: AppSpacing.xs) {
            Text(value)
                .font(.srHeading2)
                .foregroundColor(color)
            Text(label)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
        .frame(maxWidth: .infinity)
        .glassCard(padding: AppSpacing.md)
    }
}

struct ActivityRow: View {
    let icon: String
    let title: String
    let subtitle: String
    let time: String
    var color: Color = .srCyan

    var body: some View {
        HStack(spacing: AppSpacing.md) {
            ZStack {
                RoundedRectangle(cornerRadius: 10)
                    .fill(color.opacity(0.15))
                    .frame(width: 40, height: 40)
                Image(systemName: icon)
                    .font(.body)
                    .foregroundColor(color)
            }

            VStack(alignment: .leading, spacing: 2) {
                Text(title)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                Text(subtitle)
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
            }

            Spacer()

            Text(time)
                .font(.srCaption)
                .foregroundColor(.srTextMuted)
        }
        .padding(AppSpacing.md)
    }
}

#Preview {
    HomeView()
        .environmentObject(AuthViewModel())
}
