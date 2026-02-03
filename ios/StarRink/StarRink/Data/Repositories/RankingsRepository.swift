import Foundation

protocol RankingsRepositoryProtocol {
    func getRankings(sort: String, limit: Int) async throws -> RankingsResponse
}

final class RankingsRepository: RankingsRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getRankings(sort: String, limit: Int) async throws -> RankingsResponse {
        let items = [
            URLQueryItem(name: "sort", value: sort),
            URLQueryItem(name: "limit", value: String(limit))
        ]
        return try await apiClient.request(endpoint: .exploreRankings, queryItems: items)
    }
}
