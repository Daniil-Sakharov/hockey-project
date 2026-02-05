import Foundation

enum TournamentTab: String, CaseIterable {
    case standings = "Таблица"
    case matches = "Матчи"
    case scorers = "Бомбардиры"
    case teams = "Команды"
}

@MainActor
final class TournamentDetailViewModel: ObservableObject {
    @Published var selectedTab: TournamentTab = .standings
    @Published var standings: [StandingDTO] = []
    @Published var matches: [MatchDTO] = []
    @Published var scorers: [ScorerDTO] = []
    @Published var teams: [TeamItemDTO] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: TournamentDetailRepositoryProtocol
    private let teamsRepository: TeamRosterRepositoryProtocol

    init(
        repository: TournamentDetailRepositoryProtocol = TournamentDetailRepository(),
        teamsRepository: TeamRosterRepositoryProtocol = TeamRosterRepository()
    ) {
        self.repository = repository
        self.teamsRepository = teamsRepository
    }

    func loadAll(tournamentId: String, birthYear: Int?, group: String?) async {
        isLoading = true
        errorMessage = nil

        do {
            async let s = repository.getStandings(tournamentId: tournamentId, birthYear: birthYear, group: group)
            async let m = repository.getMatches(tournamentId: tournamentId, birthYear: birthYear, group: group, limit: 50)
            async let sc = repository.getScorers(tournamentId: tournamentId, birthYear: birthYear, group: group, limit: 50)
            async let t = teamsRepository.getTournamentTeams(tournamentId: tournamentId)

            let (standingsResult, matchesResult, scorersResult, teamsResult) = try await (s, m, sc, t)
            standings = standingsResult
            matches = matchesResult
            scorers = scorersResult
            teams = teamsResult.teams
        } catch {
            errorMessage = "Не удалось загрузить данные"
        }

        isLoading = false
    }
}
