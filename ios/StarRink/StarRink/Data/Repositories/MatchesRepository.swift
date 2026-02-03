import Foundation

protocol MatchesRepositoryProtocol {
    func getRecentResults(tournament: String?, limit: Int) async throws -> [MatchDTO]
    func getUpcomingMatches(tournament: String?, limit: Int) async throws -> [MatchDTO]
}

final class MatchesRepository: MatchesRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getRecentResults(tournament: String?, limit: Int) async throws -> [MatchDTO] {
        let response: MatchListResponse = try await apiClient.request(
            endpoint: .exploreResults,
            queryItems: buildQuery(tournament: tournament, limit: limit)
        )
        return response.matches
    }

    func getUpcomingMatches(tournament: String?, limit: Int) async throws -> [MatchDTO] {
        let response: MatchListResponse = try await apiClient.request(
            endpoint: .exploreCalendar,
            queryItems: buildQuery(tournament: tournament, limit: limit)
        )
        return response.matches
    }

    private func buildQuery(tournament: String?, limit: Int) -> [URLQueryItem] {
        var items = [URLQueryItem(name: "limit", value: String(limit))]
        if let tournament, !tournament.isEmpty {
            items.append(URLQueryItem(name: "tournament", value: tournament))
        }
        return items
    }
}
