import SwiftUI

struct RegionTournamentsView: View {
    let region: Region
    @StateObject private var viewModel = TournamentsViewModel()

    var body: some View {
        ZStack {
            AnimatedBackgroundView()
                .ignoresSafeArea()

            content
        }
        .navigationTitle(region.name)
        .navigationBarTitleDisplayMode(.inline)
        .toolbarBackground(.hidden, for: .navigationBar)
        .task {
            await viewModel.loadTournaments(source: region.source)
        }
    }

    @ViewBuilder
    private var content: some View {
        if viewModel.isLoading {
            loadingView
        } else if let error = viewModel.errorMessage {
            errorView(error)
        } else {
            tournamentsContent
        }
    }

    private var tournamentsContent: some View {
        ScrollView {
            VStack(spacing: AppSpacing.md) {
                pickersSection
                tournamentsSection
            }
            .padding(.top, AppSpacing.sm)
            .padding(.bottom, 100)
        }
        .scrollContentBackground(.hidden)
    }

    // MARK: - Pickers

    private var pickersSection: some View {
        VStack(spacing: AppSpacing.xs) {
            if !viewModel.seasons.isEmpty {
                pickerLabel("Сезон")
                seasonPicker
            }

            if !viewModel.birthYears.isEmpty {
                pickerLabel("Год рождения")
                birthYearPicker
            }
        }
    }

    private var seasonPicker: some View {
        HorizontalChipPicker(
            items: viewModel.seasons,
            selectedItem: Binding(
                get: { viewModel.activeSeason },
                set: { viewModel.selectSeason($0) }
            )
        )
    }

    private var birthYearPicker: some View {
        HorizontalChipPicker(
            items: viewModel.birthYears.map(String.init),
            selectedItem: Binding(
                get: { String(viewModel.activeBirthYear ?? 0) },
                set: { if let y = Int($0) { viewModel.selectBirthYear(y) } }
            )
        )
    }

    private func pickerLabel(_ title: String) -> some View {
        HStack {
            Text(title)
                .font(.system(size: 11, weight: .semibold))
                .foregroundColor(.srTextMuted)
                .textCase(.uppercase)
                .tracking(1)
            Spacer()
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
        .padding(.top, AppSpacing.xs)
    }

    // MARK: - Tournaments List

    private var tournamentsSection: some View {
        VStack(spacing: AppSpacing.sm) {
            if viewModel.filteredTournaments.isEmpty {
                emptyState
            } else {
                ForEach(expandedCards, id: \.route) { card in
                    NavigationLink(value: card.route) {
                        TournamentRow(
                            name: card.displayName,
                            groupName: card.groupName,
                            teamsCount: card.teamsCount,
                            matchesCount: card.matchesCount,
                            isEnded: card.isEnded
                        )
                    }
                }
            }
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
    }

    private var expandedCards: [(route: TournamentRoute, displayName: String, groupName: String?, teamsCount: Int, matchesCount: Int, isEnded: Bool)] {
        let year = viewModel.activeBirthYear
        return viewModel.filteredTournaments.flatMap { t in
            let cleanName = TournamentNameHelper.cleanName(t.name)
            guard let year, let groups = t.birthYearGroups?[String(year)] else {
                return [(
                    route: TournamentRoute(tournamentId: t.id, name: cleanName, birthYear: year, groupName: nil),
                    displayName: cleanName,
                    groupName: nil as String?,
                    teamsCount: t.teamsCount,
                    matchesCount: t.matchesCount,
                    isEnded: t.isEnded
                )]
            }
            if groups.count <= 1 {
                let g = groups.first
                return [(
                    route: TournamentRoute(tournamentId: t.id, name: cleanName, birthYear: year, groupName: g?.name),
                    displayName: cleanName,
                    groupName: g?.name,
                    teamsCount: g?.teamsCount ?? t.teamsCount,
                    matchesCount: g?.matchesCount ?? t.matchesCount,
                    isEnded: t.isEnded
                )]
            }
            return groups.map { g in (
                route: TournamentRoute(tournamentId: t.id, name: cleanName, birthYear: year, groupName: g.name),
                displayName: cleanName,
                groupName: g.name as String?,
                teamsCount: g.teamsCount,
                matchesCount: g.matchesCount,
                isEnded: t.isEnded
            )}
        }
    }

    // MARK: - States

    private var loadingView: some View {
        VStack(spacing: AppSpacing.md) {
            ProgressView()
                .tint(.srCyan)
            Text("Загрузка турниров...")
                .font(.srCaption)
                .foregroundColor(.srTextSecondary)
        }
    }

    private func errorView(_ message: String) -> some View {
        VStack(spacing: AppSpacing.md) {
            Image(systemName: "wifi.slash")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text(message)
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
                .multilineTextAlignment(.center)
            Button("Повторить") {
                Task { await viewModel.loadTournaments(source: region.source) }
            }
            .font(.srBodyMedium)
            .foregroundColor(.srCyan)
        }
        .padding(AppSpacing.xl)
    }

    private var emptyState: some View {
        VStack(spacing: AppSpacing.sm) {
            Image(systemName: "trophy")
                .font(.system(size: 36))
                .foregroundColor(.srTextMuted)
            Text("Нет турниров")
                .font(.srBody)
                .foregroundColor(.srTextSecondary)
            Text("Попробуйте выбрать другой сезон или год рождения")
                .font(.srCaption)
                .foregroundColor(.srTextMuted)
                .multilineTextAlignment(.center)
        }
        .frame(maxWidth: .infinity)
        .padding(.vertical, AppSpacing.xxl)
    }
}
