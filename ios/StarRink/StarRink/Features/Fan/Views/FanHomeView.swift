import SwiftUI

struct FanHomeView: View {
    @EnvironmentObject var authViewModel: AuthViewModel
    @StateObject private var viewModel = FanHomeViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.lg) {
                welcomeSection
                platformStatsSection
                topScorersSection
                recentMatchesSection
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .overlay {
            if viewModel.isLoading && viewModel.overview == nil {
                ProgressView().tint(.srCyan)
            }
        }
        .task { await viewModel.load() }
    }

    private var welcomeSection: some View {
        HStack {
            VStack(alignment: .leading, spacing: AppSpacing.xxs) {
                Text("Привет, \(authViewModel.currentUser?.displayName ?? "Гость")!")
                    .font(.srHeading3)
                    .foregroundColor(.srTextPrimary)
                Text("Следи за хоккеем")
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
            }
            Spacer()
            Image(systemName: "figure.hockey")
                .font(.system(size: 32))
                .foregroundStyle(Color.srGradientPrimary)
        }
        .glassCard()
    }

    private var platformStatsSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            Text("Платформа")
                .font(.srHeading4)
                .foregroundColor(.srTextPrimary)

            HStack(spacing: AppSpacing.sm) {
                let o = viewModel.overview
                FanStatCard(value: "\(o?.tournaments ?? 0)", label: "Турниров", icon: "trophy.fill", color: .srCyan)
                FanStatCard(value: "\(o?.teams ?? 0)", label: "Команд", icon: "person.3.fill", color: .srPurple)
                FanStatCard(value: "\(o?.players ?? 0)", label: "Игроков", icon: "figure.hockey", color: .srAmber)
            }
        }
    }

    private var topScorersSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            HStack {
                Text("Топ бомбардиры")
                    .font(.srHeading4)
                    .foregroundColor(.srTextPrimary)
                Spacer()
            }

            if viewModel.topScorers.isEmpty && !viewModel.isLoading {
                Text("Нет данных")
                    .font(.srCaption)
                    .foregroundColor(.srTextMuted)
                    .glassCard()
            } else {
                VStack(spacing: 0) {
                    ForEach(Array(viewModel.topScorers.enumerated()), id: \.element.id) { index, scorer in
                        NavigationLink(value: PlayerRoute(playerId: scorer.id)) {
                            FanScorerRow(scorer: scorer)
                        }
                        if index < viewModel.topScorers.count - 1 {
                            Divider().background(Color.srBorder.opacity(0.3))
                        }
                    }
                }
                .glassCard(padding: 0)
            }
        }
    }

    private var recentMatchesSection: some View {
        VStack(alignment: .leading, spacing: AppSpacing.md) {
            HStack {
                Text("Последние матчи")
                    .font(.srHeading4)
                    .foregroundColor(.srTextPrimary)
                Spacer()
            }

            if viewModel.recentMatches.isEmpty && !viewModel.isLoading {
                Text("Нет данных")
                    .font(.srCaption)
                    .foregroundColor(.srTextMuted)
                    .glassCard()
            } else {
                VStack(spacing: AppSpacing.sm) {
                    ForEach(viewModel.recentMatches) { match in
                        NavigationLink(value: MatchRoute(matchId: match.id)) {
                            FanMatchCard(match: match)
                        }
                    }
                }
            }
        }
    }
}
