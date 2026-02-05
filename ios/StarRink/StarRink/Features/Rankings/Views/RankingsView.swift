import SwiftUI

struct RankingsView: View {
    @StateObject private var viewModel = RankingsViewModel()
    @State private var selectedSort: String = RankingSortOption.points.rawValue
    @State private var showFilters = false

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                headerSection
                sortPicker
                if showFilters {
                    RankingsFiltersSection(viewModel: viewModel)
                }
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
        .task {
            async let _ = viewModel.load()
            async let _ = viewModel.loadFilters()
        }
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
            filterToggleButton
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
    }

    private var filterToggleButton: some View {
        Button {
            withAnimation(.easeInOut(duration: 0.2)) { showFilters.toggle() }
        } label: {
            Image(systemName: showFilters ? "line.3.horizontal.decrease.circle.fill" : "line.3.horizontal.decrease.circle")
                .font(.system(size: 24))
                .foregroundColor(hasActiveFilters ? .srCyan : .srTextMuted)
        }
    }

    private var hasActiveFilters: Bool {
        viewModel.selectedBirthYear != nil ||
        viewModel.selectedDomain != nil ||
        viewModel.selectedTournamentId != nil ||
        viewModel.selectedGroup != nil
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
