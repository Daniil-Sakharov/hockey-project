import Foundation

@MainActor
final class TournamentsViewModel: ObservableObject {
    @Published var allTournaments: [TournamentItemDTO] = []
    @Published var selectedSeason: String?
    @Published var selectedBirthYear: Int?
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: ExploreRepositoryProtocol

    init(repository: ExploreRepositoryProtocol = ExploreRepository()) {
        self.repository = repository
    }

    // MARK: - Computed

    var seasons: [String] {
        let set = Set(allTournaments.map(\.season))
        return Array(set).sorted().reversed()
    }

    var activeSeason: String {
        selectedSeason ?? seasons.first ?? ""
    }

    var seasonTournaments: [TournamentItemDTO] {
        guard !activeSeason.isEmpty else { return [] }
        return allTournaments.filter { $0.season == activeSeason }
    }

    var birthYears: [Int] {
        var years = Set<Int>()
        for tournament in seasonTournaments {
            guard let groups = tournament.birthYearGroups else { continue }
            for key in groups.keys {
                if let year = Int(key) {
                    years.insert(year)
                }
            }
        }
        return Array(years).sorted()
    }

    var activeBirthYear: Int? {
        selectedBirthYear ?? birthYears.first
    }

    var filteredTournaments: [TournamentItemDTO] {
        guard let year = activeBirthYear else { return seasonTournaments }
        return seasonTournaments.filter { tournament in
            guard let groups = tournament.birthYearGroups else { return true }
            return groups.keys.contains(String(year))
        }
    }

    // MARK: - Actions

    func loadTournaments(source: String) async {
        isLoading = true
        errorMessage = nil

        do {
            allTournaments = try await repository.getTournaments(source: source)
        } catch let error as APIError {
            errorMessage = error.errorDescription
        } catch {
            errorMessage = "Не удалось загрузить турниры"
        }

        isLoading = false
    }

    func selectSeason(_ season: String) {
        selectedSeason = season
        selectedBirthYear = nil
    }

    func selectBirthYear(_ year: Int) {
        selectedBirthYear = year
    }
}
