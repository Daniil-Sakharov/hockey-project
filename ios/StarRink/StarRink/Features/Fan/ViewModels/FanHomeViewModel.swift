import Foundation

@MainActor
final class FanHomeViewModel: ObservableObject {
    @Published var overview: ExploreOverviewDTO?
    @Published var recentMatches: [MatchDTO] = []
    @Published var topScorers: [RankedPlayerDTO] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let overviewRepo: ExploreOverviewRepositoryProtocol
    private let matchesRepo: MatchesRepositoryProtocol
    private let rankingsRepo: RankingsRepositoryProtocol

    init(
        overviewRepo: ExploreOverviewRepositoryProtocol = ExploreOverviewRepository(),
        matchesRepo: MatchesRepositoryProtocol = MatchesRepository(),
        rankingsRepo: RankingsRepositoryProtocol = RankingsRepository()
    ) {
        self.overviewRepo = overviewRepo
        self.matchesRepo = matchesRepo
        self.rankingsRepo = rankingsRepo
    }

    func load() async {
        isLoading = true
        errorMessage = nil

        do {
            async let o = overviewRepo.getOverview()
            async let m = matchesRepo.getRecentResults(tournament: nil, limit: 5)
            async let r = rankingsRepo.getRankings(sort: "points", limit: 3, birthYear: nil, domain: nil, tournamentId: nil, groupName: nil)

            let (overviewResult, matchesResult, rankingsResult) = try await (o, m, r)
            overview = overviewResult
            recentMatches = matchesResult
            topScorers = rankingsResult.players
        } catch {
            errorMessage = "Не удалось загрузить данные"
        }

        isLoading = false
    }
}
