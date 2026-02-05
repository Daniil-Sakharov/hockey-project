import Foundation

enum RankingSortOption: String, CaseIterable {
    case points = "Очки"
    case goals = "Голы"
    case assists = "Передачи"
    case plusMinus = "+/-"
    case penaltyMinutes = "Штраф"

    var apiValue: String {
        switch self {
        case .points: return "points"
        case .goals: return "goals"
        case .assists: return "assists"
        case .plusMinus: return "plusMinus"
        case .penaltyMinutes: return "penaltyMinutes"
        }
    }
}

@MainActor
final class RankingsViewModel: ObservableObject {
    @Published var players: [RankedPlayerDTO] = []
    @Published var season: String = ""
    @Published var sortBy: RankingSortOption = .points
    @Published var isLoading = false
    @Published var errorMessage: String?

    // Filters
    @Published var filters: RankingsFiltersResponse?
    @Published var selectedBirthYear: Int?
    @Published var selectedDomain: String?
    @Published var selectedTournamentId: String?
    @Published var selectedGroup: String?
    @Published var isFiltersLoading = false

    private let repository: RankingsRepositoryProtocol

    init(repository: RankingsRepositoryProtocol = RankingsRepository()) {
        self.repository = repository
    }

    func load() async {
        isLoading = true
        errorMessage = nil

        do {
            let response = try await repository.getRankings(
                sort: sortBy.apiValue,
                limit: 50,
                birthYear: selectedBirthYear,
                domain: selectedDomain,
                tournamentId: selectedTournamentId,
                groupName: selectedGroup
            )
            players = response.players
            season = response.season
        } catch {
            errorMessage = "Не удалось загрузить рейтинг"
        }

        isLoading = false
    }

    func loadFilters() async {
        isFiltersLoading = true
        do {
            filters = try await repository.getFilters()
        } catch {
            // Filters are optional, silently fail
        }
        isFiltersLoading = false
    }

    func changeSortAndReload(_ option: RankingSortOption) {
        guard sortBy != option else { return }
        sortBy = option
        Task { await load() }
    }

    func applyFilters() {
        Task { await load() }
    }

    func resetFilters() {
        selectedBirthYear = nil
        selectedDomain = nil
        selectedTournamentId = nil
        selectedGroup = nil
        Task { await load() }
    }

    var availableTournaments: [TournamentOption] {
        guard let filters else { return [] }
        return filters.tournaments.filter { t in
            if let domain = selectedDomain, !domain.isEmpty {
                return t.domain == domain
            }
            return true
        }
    }

    var availableGroups: [GroupOption] {
        guard let filters else { return [] }
        return filters.groups.filter { g in
            if let tid = selectedTournamentId, !tid.isEmpty {
                return g.tournamentId == tid
            }
            return true
        }
    }
}
