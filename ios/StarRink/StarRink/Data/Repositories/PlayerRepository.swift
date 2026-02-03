import Foundation

protocol PlayerRepositoryProtocol {
    func searchPlayers(query: String, position: String, season: String, birthYear: Int, limit: Int, offset: Int) async throws -> PlayersSearchResponse
    func getProfile(id: String, season: String?) async throws -> PlayerProfileDTO
    func getStats(id: String) async throws -> [PlayerStatEntryDTO]
    func getSeasons() async throws -> [String]
}

final class PlayerRepository: PlayerRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func searchPlayers(query: String, position: String, season: String, birthYear: Int, limit: Int, offset: Int) async throws -> PlayersSearchResponse {
        var items: [URLQueryItem] = []
        if !query.isEmpty { items.append(URLQueryItem(name: "q", value: query)) }
        if !position.isEmpty, position != "all" { items.append(URLQueryItem(name: "position", value: position)) }
        if !season.isEmpty { items.append(URLQueryItem(name: "season", value: season)) }
        if birthYear > 0 { items.append(URLQueryItem(name: "birthYear", value: String(birthYear))) }
        items.append(URLQueryItem(name: "limit", value: String(limit)))
        items.append(URLQueryItem(name: "offset", value: String(offset)))

        return try await apiClient.request(endpoint: .explorePlayers, queryItems: items)
    }

    func getProfile(id: String, season: String?) async throws -> PlayerProfileDTO {
        var items: [URLQueryItem]?
        if let season, !season.isEmpty {
            items = [URLQueryItem(name: "season", value: season)]
        }
        return try await apiClient.request(endpoint: .playerProfile(id: id), queryItems: items)
    }

    func getStats(id: String) async throws -> [PlayerStatEntryDTO] {
        let response: PlayerStatsHistoryResponse = try await apiClient.request(
            endpoint: .playerStatsHistory(id: id)
        )
        return response.stats
    }

    func getSeasons() async throws -> [String] {
        let response: SeasonsResponse = try await apiClient.request(endpoint: .seasons)
        return response.seasons
    }
}
