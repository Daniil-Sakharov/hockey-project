import Foundation

@MainActor
final class PlayerProfileViewModel: ObservableObject {
    @Published var profile: PlayerProfileDTO?
    @Published var statsHistory: [PlayerStatEntryDTO] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: PlayerRepositoryProtocol

    init(repository: PlayerRepositoryProtocol = PlayerRepository()) {
        self.repository = repository
    }

    func load(playerId: String) async {
        isLoading = true
        errorMessage = nil

        do {
            async let p = repository.getProfile(id: playerId, season: nil)
            async let s = repository.getStats(id: playerId)

            let (profileResult, statsResult) = try await (p, s)
            profile = profileResult
            statsHistory = statsResult
        } catch {
            errorMessage = "Не удалось загрузить профиль"
        }

        isLoading = false
    }

    var groupedHistory: [(season: String, entries: [PlayerStatEntryDTO])] {
        let grouped = Dictionary(grouping: statsHistory, by: \.season)
        return grouped.keys.sorted().reversed().map { season in
            (season: season, entries: grouped[season] ?? [])
        }
    }

    var aggregatedBySeason: [SeasonAggregated] {
        let grouped = Dictionary(grouping: statsHistory, by: \.season)
        return grouped.keys.sorted().map { season in
            let entries = grouped[season] ?? []
            return SeasonAggregated(
                id: season,
                season: season,
                games: entries.reduce(0) { $0 + $1.games },
                goals: entries.reduce(0) { $0 + $1.goals },
                assists: entries.reduce(0) { $0 + $1.assists },
                points: entries.reduce(0) { $0 + $1.points },
                plusMinus: entries.reduce(0) { $0 + $1.plusMinus },
                penaltyMinutes: entries.reduce(0) { $0 + $1.penaltyMinutes }
            )
        }
    }

    var positionLocalized: String {
        guard let pos = profile?.position else { return "" }
        switch pos {
        case "forward": return "Нападающий"
        case "defender": return "Защитник"
        case "goalie": return "Вратарь"
        default: return pos
        }
    }

    var handednessLocalized: String {
        guard let h = profile?.handedness else { return "" }
        switch h {
        case "left": return "Левый"
        case "right": return "Правый"
        default: return h
        }
    }
}
