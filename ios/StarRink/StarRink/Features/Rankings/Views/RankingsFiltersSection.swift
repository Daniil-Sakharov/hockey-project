import SwiftUI

struct RankingsFiltersSection: View {
    @ObservedObject var viewModel: RankingsViewModel

    var body: some View {
        VStack(alignment: .leading, spacing: AppSpacing.sm) {
            if let filters = viewModel.filters {
                birthYearFilter(filters.birthYears)
                domainFilter(filters.domains)
                tournamentFilter
                groupFilter
                resetButton
            }
        }
        .padding(.horizontal, AppSpacing.screenHorizontal)
        .transition(.opacity.combined(with: .move(edge: .top)))
    }

    private func birthYearFilter(_ years: [Int]) -> some View {
        Group {
            if !years.isEmpty {
                filterLabel("Год рождения")
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: AppSpacing.xs) {
                        ForEach(years, id: \.self) { year in
                            chipButton(
                                String(year),
                                isSelected: viewModel.selectedBirthYear == year
                            ) {
                                viewModel.selectedBirthYear = viewModel.selectedBirthYear == year ? nil : year
                                viewModel.applyFilters()
                            }
                        }
                    }
                }
            }
        }
    }

    private func domainFilter(_ domains: [DomainOption]) -> some View {
        Group {
            if !domains.isEmpty {
                filterLabel("Регион")
                ScrollView(.horizontal, showsIndicators: false) {
                    HStack(spacing: AppSpacing.xs) {
                        ForEach(domains) { d in
                            chipButton(
                                d.label,
                                isSelected: viewModel.selectedDomain == d.domain
                            ) {
                                viewModel.selectedDomain = viewModel.selectedDomain == d.domain ? nil : d.domain
                                viewModel.selectedTournamentId = nil
                                viewModel.selectedGroup = nil
                                viewModel.applyFilters()
                            }
                        }
                    }
                }
            }
        }
    }

    private var tournamentFilter: some View {
        Group {
            if !viewModel.availableTournaments.isEmpty {
                filterLabel("Турнир")
                Menu {
                    Button("Все турниры") {
                        viewModel.selectedTournamentId = nil
                        viewModel.selectedGroup = nil
                        viewModel.applyFilters()
                    }
                    ForEach(viewModel.availableTournaments) { t in
                        Button(t.name) {
                            viewModel.selectedTournamentId = t.id
                            viewModel.selectedGroup = nil
                            viewModel.applyFilters()
                        }
                    }
                } label: {
                    menuLabel(
                        viewModel.availableTournaments.first(where: { $0.id == viewModel.selectedTournamentId })?.name ?? "Все турниры"
                    )
                }
            }
        }
    }

    private var groupFilter: some View {
        Group {
            if !viewModel.availableGroups.isEmpty {
                filterLabel("Группа")
                Menu {
                    Button("Все группы") {
                        viewModel.selectedGroup = nil
                        viewModel.applyFilters()
                    }
                    ForEach(viewModel.availableGroups) { g in
                        Button(g.name) {
                            viewModel.selectedGroup = g.name
                            viewModel.applyFilters()
                        }
                    }
                } label: {
                    menuLabel(viewModel.selectedGroup ?? "Все группы")
                }
            }
        }
    }

    private var resetButton: some View {
        Group {
            if viewModel.selectedBirthYear != nil ||
               viewModel.selectedDomain != nil ||
               viewModel.selectedTournamentId != nil ||
               viewModel.selectedGroup != nil {
                Button("Сбросить фильтры") { viewModel.resetFilters() }
                    .font(.system(size: 12, weight: .medium))
                    .foregroundColor(.srAmber)
            }
        }
    }

    // MARK: - Helpers

    private func filterLabel(_ text: String) -> some View {
        Text(text)
            .font(.system(size: 10, weight: .semibold))
            .foregroundColor(.srTextMuted)
            .textCase(.uppercase)
            .tracking(0.5)
    }

    private func chipButton(_ title: String, isSelected: Bool, action: @escaping () -> Void) -> some View {
        Button(action: action) {
            Text(title)
                .font(.system(size: 12, weight: .medium))
                .foregroundColor(isSelected ? .srBackground : .srTextPrimary)
                .padding(.horizontal, 10)
                .padding(.vertical, 6)
                .background(isSelected ? Color.srCyan : Color.srCard)
                .clipShape(Capsule())
        }
    }

    private func menuLabel(_ text: String) -> some View {
        HStack {
            Text(text)
                .font(.system(size: 12, weight: .medium))
                .foregroundColor(.srTextPrimary)
            Image(systemName: "chevron.down")
                .font(.system(size: 10))
                .foregroundColor(.srTextMuted)
        }
        .padding(.horizontal, 12)
        .padding(.vertical, 8)
        .background(Color.srCard)
        .clipShape(Capsule())
    }
}
