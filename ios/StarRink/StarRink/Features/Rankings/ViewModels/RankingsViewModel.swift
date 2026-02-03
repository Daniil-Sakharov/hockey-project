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
                limit: 50
            )
            players = response.players
            season = response.season
        } catch {
            errorMessage = "Не удалось загрузить рейтинг"
        }

        isLoading = false
    }

    func changeSortAndReload(_ option: RankingSortOption) {
        guard sortBy != option else { return }
        sortBy = option
        Task { await load() }
    }
}
