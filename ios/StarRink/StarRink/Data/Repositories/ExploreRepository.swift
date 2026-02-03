import Foundation

protocol ExploreRepositoryProtocol {
    func getTournaments(source: String?) async throws -> [TournamentItemDTO]
}

final class ExploreRepository: ExploreRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getTournaments(source: String?) async throws -> [TournamentItemDTO] {
        var queryItems: [URLQueryItem]?
        if let source, !source.isEmpty {
            queryItems = [URLQueryItem(name: "source", value: source)]
        }

        let response: TournamentsResponse = try await apiClient.request(
            endpoint: .exploreTournaments,
            queryItems: queryItems
        )
        return response.tournaments
    }
}
