import Foundation

protocol ExploreOverviewRepositoryProtocol {
    func getOverview() async throws -> ExploreOverviewDTO
}

final class ExploreOverviewRepository: ExploreOverviewRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getOverview() async throws -> ExploreOverviewDTO {
        try await apiClient.request(endpoint: .exploreOverview)
    }
}
