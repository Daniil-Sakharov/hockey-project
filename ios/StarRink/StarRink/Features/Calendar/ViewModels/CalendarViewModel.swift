import Foundation

@MainActor
final class CalendarViewModel: ObservableObject {
    @Published var upcomingMatches: [MatchDTO] = []
    @Published var selectedTournament: String?
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: MatchesRepositoryProtocol

    init(repository: MatchesRepositoryProtocol = MatchesRepository()) {
        self.repository = repository
    }

    func load() async {
        isLoading = true
        errorMessage = nil

        do {
            upcomingMatches = try await repository.getUpcomingMatches(
                tournament: selectedTournament,
                limit: 50
            )
        } catch {
            errorMessage = "Не удалось загрузить календарь"
        }

        isLoading = false
    }

    var groupedByDate: [(date: String, matches: [MatchDTO])] {
        let grouped = Dictionary(grouping: upcomingMatches, by: \.date)
        return grouped.keys.sorted().map { date in
            (date: date, matches: grouped[date] ?? [])
        }
    }

    var tournamentNames: [String] {
        Array(Set(upcomingMatches.compactMap(\.tournament))).sorted()
    }
}
