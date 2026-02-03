import SwiftUI

struct PlayerProfileView: View {
    let playerId: String
    @StateObject private var viewModel = PlayerProfileViewModel()

    var body: some View {
        ZStack {
            AnimatedBackgroundView().ignoresSafeArea()
            content
        }
        .navigationTitle(viewModel.profile?.name ?? "Игрок")
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(.hidden, for: .navigationBar)
        .task { await viewModel.load(playerId: playerId) }
    }

    @ViewBuilder
    private var content: some View {
        if viewModel.isLoading {
            ProgressView().tint(.srCyan)
        } else if let error = viewModel.errorMessage {
            errorView(error)
        } else if let profile = viewModel.profile {
            profileContent(profile)
        }
    }

    private func profileContent(_ p: PlayerProfileDTO) -> some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                heroCard(p)
                if let stats = p.stats {
                    statsGrid(stats)
                }
                if viewModel.aggregatedBySeason.count >= 2 {
                    PlayerChartsSection(seasonData: viewModel.aggregatedBySeason)
                }
                if !viewModel.statsHistory.isEmpty {
                    PlayerStatsHistoryView(groupedHistory: viewModel.groupedHistory)
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.sm)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
    }

    // MARK: - Hero Card

    private func heroCard(_ p: PlayerProfileDTO) -> some View {
        VStack(spacing: AppSpacing.md) {
            playerAvatar(p)
            Text(p.name)
                .font(.srHeading3)
                .foregroundColor(.srTextPrimary)
            Text(p.team)
                .font(.srBody)
                .foregroundColor(.srTextSecondary)

            HStack(spacing: AppSpacing.lg) {
                infoPill(viewModel.positionLocalized)
                if p.birthYear > 0 {
                    infoPill("\(p.birthYear) г.р.")
                }
                if let h = p.height, h > 0 {
                    infoPill("\(h) см")
                }
                if let w = p.weight, w > 0 {
                    infoPill("\(w) кг")
                }
            }

            if !viewModel.handednessLocalized.isEmpty {
                Text("Хват: \(viewModel.handednessLocalized)")
                    .font(.srCaption)
                    .foregroundColor(.srTextMuted)
            }
        }
        .frame(maxWidth: .infinity)
        .glassCard()
    }

    private func playerAvatar(_ p: PlayerProfileDTO) -> some View {
        CachedAsyncImage(
            url: p.photoUrl.flatMap { URL(string: $0) }
        ) {
            avatarPlaceholder(p)
        }
        .frame(width: 88, height: 88)
        .clipShape(Circle())
        .overlay(
            Circle().stroke(Color.srCyan.opacity(0.3), lineWidth: 2)
        )
    }

    private func avatarPlaceholder(_ p: PlayerProfileDTO) -> some View {
        ZStack {
            Circle().fill(Color.srCyan.opacity(0.15))
            Text(String(p.name.prefix(1)))
                .font(.system(size: 32, weight: .bold))
                .foregroundColor(.srCyan)
        }
    }

    private func infoPill(_ text: String) -> some View {
        Text(text)
            .font(.system(size: 11, weight: .medium))
            .foregroundColor(.srTextSecondary)
            .padding(.horizontal, 8)
            .padding(.vertical, 4)
            .background(Capsule().fill(Color.white.opacity(0.05)))
    }

    // MARK: - Stats Grid

    private func statsGrid(_ stats: PlayerStatsDTO) -> some View {
        LazyVGrid(columns: Array(repeating: GridItem(.flexible(), spacing: AppSpacing.sm), count: 3), spacing: AppSpacing.sm) {
            statCard("Игры", value: "\(stats.games)", color: .srTextPrimary)
            statCard("Голы", value: "\(stats.goals)", color: .srCyan)
            statCard("Передачи", value: "\(stats.assists)", color: .srPurple)
            statCard("Очки", value: "\(stats.points)", color: .srAmber)
            statCard("+/-", value: "\(stats.plusMinus)", color: stats.plusMinus >= 0 ? .srSuccess : .srError)
            statCard("Штраф", value: "\(stats.penaltyMinutes)", color: .srTextSecondary)
        }
    }

    private func statCard(_ title: String, value: String, color: Color) -> some View {
        VStack(spacing: 4) {
            Text(value)
                .font(.system(size: 22, weight: .bold))
                .foregroundColor(color)
            Text(title)
                .font(.system(size: 10, weight: .medium))
                .foregroundColor(.srTextMuted)
        }
        .frame(maxWidth: .infinity)
        .glassCard(padding: AppSpacing.sm)
    }

    // MARK: - Error

    private func errorView(_ message: String) -> some View {
        VStack(spacing: AppSpacing.md) {
            Image(systemName: "person.slash")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text(message)
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
            Button("Повторить") {
                Task { await viewModel.load(playerId: playerId) }
            }
            .font(.srBodyMedium)
            .foregroundColor(.srCyan)
        }
        .padding(AppSpacing.xl)
    }
}
