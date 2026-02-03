import SwiftUI

struct RankingsView: View {
    @StateObject private var viewModel = RankingsViewModel()
    @State private var selectedSort: String = RankingSortOption.points.rawValue

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                headerSection
                sortPicker
                playersList
            }
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .overlay {
            if viewModel.isLoading && viewModel.players.isEmpty {
                ProgressView().tint(.srCyan)
            }
        }
        .task { await viewModel.load() }
        .onChange(of: selectedSort) { _, newValue in
            if let option = RankingSortOption.allCases.first(where: { $0.rawValue == newValue }) {
                viewModel.changeSortAndReload(option)
            }
        }
    }

    private var headerSection: some View {
        HStack {
            VStack(alignment: .leading, spacing: AppSpacing.xxs) {
                Text("Рейтинг игроков")
                    .font(.srHeading3)
                    .foregroundColor(.srTextPrimary)
                if !viewModel.season.isEmpty {
                    Text("Сезон \(viewModel.season)")
                        .font(.srCaption)
                        .foregroundColor(.srTextSecondary)
                }
            }
            Spacer()
            Image(systemName: "medal.fill")
                .font(.system(size: 28))
                .foregroundColor(.srAmber)
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
    }

    private var sortPicker: some View {
        HorizontalChipPicker(
            items: RankingSortOption.allCases.map(\.rawValue),
            selectedItem: $selectedSort
        )
    }

    private var playersList: some View {
        LazyVStack(spacing: 0) {
            ForEach(viewModel.players) { player in
                RankingPlayerRow(player: player)
                Divider().background(Color.srBorder.opacity(0.3))
            }
        }
        .glassCard(padding: 0)
        .padding(.horizontal, AppSpacing.screenHorizontal)
    }
}

// MARK: - Row

private struct RankingPlayerRow: View {
    let player: RankedPlayerDTO

    private var rankColor: Color {
        switch player.rank {
        case 1: return .srAmber
        case 2: return .srTextSecondary
        case 3: return Color(red: 0.72, green: 0.45, blue: 0.2)
        default: return .srCyan
        }
    }

    var body: some View {
        HStack(spacing: AppSpacing.sm) {
            Text("\(player.rank)")
                .font(.srHeading4)
                .foregroundColor(rankColor)
                .frame(width: 30)

            CachedAsyncImage(url: URL(string: player.photoUrl ?? "")) {
                Image(systemName: "person.circle.fill")
                    .resizable()
                    .foregroundColor(.srTextMuted)
            }
            .frame(width: 36, height: 36)
            .clipShape(Circle())

            VStack(alignment: .leading, spacing: 2) {
                Text(player.name)
                    .font(.srBodyMedium)
                    .foregroundColor(.srTextPrimary)
                    .lineLimit(1)
                Text(player.team)
                    .font(.srCaption)
                    .foregroundColor(.srTextSecondary)
                    .lineLimit(1)
            }

            Spacer()

            HStack(spacing: AppSpacing.sm) {
                StatCell(value: player.games, label: "И")
                StatCell(value: player.goals, label: "Г")
                StatCell(value: player.assists, label: "П")
                StatCell(value: player.points, label: "О")
            }
        }
        .padding(.horizontal, AppSpacing.md)
        .padding(.vertical, AppSpacing.sm)
    }
}

private struct StatCell: View {
    let value: Int
    let label: String

    var body: some View {
        VStack(spacing: 1) {
            Text("\(value)")
                .font(.system(size: 13, weight: .semibold, design: .monospaced))
                .foregroundColor(.srTextPrimary)
            Text(label)
                .font(.system(size: 9))
                .foregroundColor(.srTextMuted)
        }
        .frame(width: 28)
    }
}
