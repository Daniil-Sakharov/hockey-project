import Foundation

@MainActor
final class MatchResultsViewModel: ObservableObject {
    @Published var results: [MatchDTO] = []
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
            results = try await repository.getRecentResults(
                tournament: selectedTournament,
                limit: 50
            )
        } catch {
            errorMessage = "Не удалось загрузить результаты"
        }

        isLoading = false
    }

    var groupedByDate: [(date: String, matches: [MatchDTO])] {
        let grouped = Dictionary(grouping: results, by: \.date)
        return grouped.keys.sorted().reversed().map { date in
            (date: date, matches: grouped[date] ?? [])
        }
    }

    var tournamentNames: [String] {
        Array(Set(results.compactMap(\.tournament))).sorted()
    }
}
