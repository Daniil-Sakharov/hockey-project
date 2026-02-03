import Foundation

protocol TeamProfileRepositoryProtocol {
    func getTeamProfile(id: String) async throws -> TeamProfileDTO
}

final class TeamProfileRepository: TeamProfileRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getTeamProfile(id: String) async throws -> TeamProfileDTO {
        try await apiClient.request(endpoint: .teamProfile(id: id))
    }
}
