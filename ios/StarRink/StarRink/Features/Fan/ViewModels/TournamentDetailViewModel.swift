import Foundation

enum TournamentTab: String, CaseIterable {
    case standings = "Таблица"
    case matches = "Матчи"
    case scorers = "Бомбардиры"
}

@MainActor
final class TournamentDetailViewModel: ObservableObject {
    @Published var selectedTab: TournamentTab = .standings
    @Published var standings: [StandingDTO] = []
    @Published var matches: [MatchDTO] = []
    @Published var scorers: [ScorerDTO] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: TournamentDetailRepositoryProtocol

    init(repository: TournamentDetailRepositoryProtocol = TournamentDetailRepository()) {
        self.repository = repository
    }

    func loadAll(tournamentId: String, birthYear: Int?, group: String?) async {
        isLoading = true
        errorMessage = nil

        do {
            async let s = repository.getStandings(tournamentId: tournamentId, birthYear: birthYear, group: group)
            async let m = repository.getMatches(tournamentId: tournamentId, birthYear: birthYear, group: group, limit: 50)
            async let sc = repository.getScorers(tournamentId: tournamentId, birthYear: birthYear, group: group, limit: 50)

            let (standingsResult, matchesResult, scorersResult) = try await (s, m, sc)
            standings = standingsResult
            matches = matchesResult
            scorers = scorersResult
        } catch {
            errorMessage = "Не удалось загрузить данные"
        }

        isLoading = false
    }
}
