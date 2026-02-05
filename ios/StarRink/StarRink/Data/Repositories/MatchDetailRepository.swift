import Foundation

protocol MatchDetailRepositoryProtocol {
    func getMatchDetail(id: String) async throws -> MatchDetailResponse
}

final class MatchDetailRepository: MatchDetailRepositoryProtocol {
    func getMatchDetail(id: String) async throws -> MatchDetailResponse {
        try await APIClient.shared.request(endpoint: .matchDetail(id: id))
    }
}
