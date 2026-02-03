import SwiftUI

struct PlayersSearchView: View {
    @StateObject private var viewModel = PlayersSearchViewModel()

    var body: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                searchBar
                filtersSection
                resultsSection
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
            .padding(.top, AppSpacing.md)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
        .background(Color.clear)
        .task {
            await viewModel.loadSeasons()
            viewModel.search()
        }
    }

    // MARK: - Search Bar

    private var searchBar: some View {
        HStack(spacing: AppSpacing.sm) {
            Image(systemName: "magnifyingglass")
                .foregroundColor(.srTextMuted)
            TextField("Поиск игроков...", text: $viewModel.searchText)
                .font(.srBody)
                .foregroundColor(.srTextPrimary)
                .autocorrectionDisabled()
                .onChange(of: viewModel.searchText) { _, _ in
                    viewModel.onSearchTextChanged()
                }
            if !viewModel.searchText.isEmpty {
                Button { viewModel.searchText = ""; viewModel.search() } label: {
                    Image(systemName: "xmark.circle.fill")
                        .foregroundColor(.srTextMuted)
                }
            }
        }
        .padding(AppSpacing.md)
        .background(
            RoundedRectangle(cornerRadius: AppSpacing.radiusMedium)
                .fill(.ultraThinMaterial.opacity(0.5))
                .overlay(
                    RoundedRectangle(cornerRadius: AppSpacing.radiusMedium)
                        .stroke(Color.srBorder.opacity(0.3), lineWidth: 0.5)
                )
        )
    }

    // MARK: - Filters

    private var filtersSection: some View {
        VStack(spacing: AppSpacing.xs) {
            positionFilter
            if !viewModel.seasons.isEmpty {
                seasonFilter
            }
            birthYearFilter
        }
    }

    private var positionFilter: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: AppSpacing.xs) {
                ForEach(PositionFilter.allCases, id: \.self) { pos in
                    chipButton(pos.rawValue, isSelected: viewModel.selectedPosition == pos) {
                        viewModel.selectedPosition = pos
                        viewModel.search()
                    }
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
        }
        .padding(.horizontal, -AppSpacing.screenHorizontal)
    }

    private var seasonFilter: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: AppSpacing.xs) {
                chipButton("Все сезоны", isSelected: viewModel.selectedSeason.isEmpty) {
                    viewModel.selectedSeason = ""
                    viewModel.search()
                }
                ForEach(viewModel.seasons, id: \.self) { season in
                    chipButton(season, isSelected: viewModel.selectedSeason == season) {
                        viewModel.selectedSeason = season
                        viewModel.search()
                    }
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
        }
        .padding(.horizontal, -AppSpacing.screenHorizontal)
    }

    private var birthYearFilter: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack(spacing: AppSpacing.xs) {
                chipButton("Все годы", isSelected: viewModel.selectedBirthYear == 0) {
                    viewModel.selectedBirthYear = 0
                    viewModel.search()
                }
                ForEach((2008...2016).reversed(), id: \.self) { year in
                    chipButton(String(year), isSelected: viewModel.selectedBirthYear == year) {
                        viewModel.selectedBirthYear = year
                        viewModel.search()
                    }
                }
            }
            .padding(.horizontal, AppSpacing.screenHorizontal)
        }
        .padding(.horizontal, -AppSpacing.screenHorizontal)
    }

    private func chipButton(_ title: String, isSelected: Bool, action: @escaping () -> Void) -> some View {
        Button(action: action) {
            Text(title)
                .font(.system(size: 12, weight: isSelected ? .semibold : .regular))
                .foregroundColor(isSelected ? .white : .srTextSecondary)
                .padding(.horizontal, 12)
                .padding(.vertical, 6)
                .background(
                    Capsule().fill(isSelected ? Color.srCyan.opacity(0.3) : Color.white.opacity(0.05))
                )
                .overlay(
                    Capsule().stroke(isSelected ? Color.srCyan.opacity(0.5) : Color.clear, lineWidth: 1)
                )
        }
    }

    // MARK: - Results

    @ViewBuilder
    private var resultsSection: some View {
        if viewModel.isLoading {
            ProgressView().tint(.srCyan).padding(.top, AppSpacing.xl)
        } else if let error = viewModel.errorMessage {
            Text(error)
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
                .padding(.top, AppSpacing.xl)
        } else if viewModel.players.isEmpty {
            VStack(spacing: AppSpacing.sm) {
                Image(systemName: "person.slash")
                    .font(.system(size: 36))
                    .foregroundColor(.srTextMuted)
                Text("Игроки не найдены")
                    .font(.srBody)
                    .foregroundColor(.srTextSecondary)
            }
            .padding(.top, AppSpacing.xxl)
        } else {
            VStack(alignment: .leading, spacing: AppSpacing.xs) {
                Text("Найдено: \(viewModel.total)")
                    .font(.srCaption)
                    .foregroundColor(.srTextMuted)
                ForEach(viewModel.players) { player in
                    NavigationLink(value: PlayerRoute(playerId: player.id)) {
                        PlayerCardView(player: player)
                    }
                }
            }
        }
    }
}
