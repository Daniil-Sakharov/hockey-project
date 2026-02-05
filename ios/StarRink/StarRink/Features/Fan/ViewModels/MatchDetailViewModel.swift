import Foundation

@MainActor
final class MatchDetailViewModel: ObservableObject {
    @Published var match: MatchDetailResponse?
    @Published var eventsByPeriod: [(period: String, events: [MatchEventDTO])] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: MatchDetailRepositoryProtocol

    init(repository: MatchDetailRepositoryProtocol = MatchDetailRepository()) {
        self.repository = repository
    }

    func load(matchId: String) async {
        isLoading = true
        errorMessage = nil

        do {
            let detail = try await repository.getMatchDetail(id: matchId)
            match = detail
            eventsByPeriod = groupEventsByPeriod(detail.events)
        } catch {
            errorMessage = "Не удалось загрузить матч"
        }

        isLoading = false
    }

    private func groupEventsByPeriod(_ events: [MatchEventDTO]) -> [(period: String, events: [MatchEventDTO])] {
        let grouped = Dictionary(grouping: events) { $0.period ?? 0 }
        return grouped
            .sorted { $0.key < $1.key }
            .map { (period: periodLabel($0.key), events: $0.value.sorted { $0.time < $1.time }) }
    }

    private func periodLabel(_ raw: Int) -> String {
        switch raw {
        case 1: return "1-й период"
        case 2: return "2-й период"
        case 3: return "3-й период"
        case 4: return "Овертайм"
        case 5: return "Буллиты"
        default: return "\(raw)-й период"
        }
    }
}
