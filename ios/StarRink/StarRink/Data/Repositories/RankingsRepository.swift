import Foundation

protocol RankingsRepositoryProtocol {
    func getRankings(sort: String, limit: Int, birthYear: Int?, domain: String?, tournamentId: String?, groupName: String?) async throws -> RankingsResponse
    func getFilters() async throws -> RankingsFiltersResponse
}

final class RankingsRepository: RankingsRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getRankings(sort: String, limit: Int, birthYear: Int?, domain: String?, tournamentId: String?, groupName: String?) async throws -> RankingsResponse {
        var items = [
            URLQueryItem(name: "sort", value: sort),
            URLQueryItem(name: "limit", value: String(limit))
        ]
        if let birthYear, birthYear > 0 {
            items.append(URLQueryItem(name: "birthYear", value: String(birthYear)))
        }
        if let domain, !domain.isEmpty {
            items.append(URLQueryItem(name: "domain", value: domain))
        }
        if let tournamentId, !tournamentId.isEmpty {
            items.append(URLQueryItem(name: "tournamentId", value: tournamentId))
        }
        if let groupName, !groupName.isEmpty {
            items.append(URLQueryItem(name: "groupName", value: groupName))
        }
        return try await apiClient.request(endpoint: .exploreRankings, queryItems: items)
    }

    func getFilters() async throws -> RankingsFiltersResponse {
        try await apiClient.request(endpoint: .rankingsFilters)
    }
}
