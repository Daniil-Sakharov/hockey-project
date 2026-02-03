import Foundation

protocol TournamentDetailRepositoryProtocol {
    func getStandings(tournamentId: String, birthYear: Int?, group: String?) async throws -> [StandingDTO]
    func getMatches(tournamentId: String, birthYear: Int?, group: String?, limit: Int) async throws -> [MatchDTO]
    func getScorers(tournamentId: String, birthYear: Int?, group: String?, limit: Int) async throws -> [ScorerDTO]
}

final class TournamentDetailRepository: TournamentDetailRepositoryProtocol {
    private let apiClient: APIClient

    init(apiClient: APIClient = .shared) {
        self.apiClient = apiClient
    }

    func getStandings(tournamentId: String, birthYear: Int?, group: String?) async throws -> [StandingDTO] {
        let response: StandingsResponse = try await apiClient.request(
            endpoint: .tournamentStandings(id: tournamentId),
            queryItems: buildQuery(birthYear: birthYear, group: group)
        )
        return response.standings
    }

    func getMatches(tournamentId: String, birthYear: Int?, group: String?, limit: Int) async throws -> [MatchDTO] {
        var items = buildQuery(birthYear: birthYear, group: group)
        items.append(URLQueryItem(name: "limit", value: String(limit)))
        let response: MatchListResponse = try await apiClient.request(
            endpoint: .tournamentMatches(id: tournamentId),
            queryItems: items
        )
        return response.matches
    }

    func getScorers(tournamentId: String, birthYear: Int?, group: String?, limit: Int) async throws -> [ScorerDTO] {
        var items = buildQuery(birthYear: birthYear, group: group)
        items.append(URLQueryItem(name: "limit", value: String(limit)))
        let response: ScorersResponse = try await apiClient.request(
            endpoint: .tournamentScorers(id: tournamentId),
            queryItems: items
        )
        return response.scorers
    }

    private func buildQuery(birthYear: Int?, group: String?) -> [URLQueryItem] {
        var items: [URLQueryItem] = []
        if let birthYear {
            items.append(URLQueryItem(name: "birthYear", value: String(birthYear)))
        }
        if let group, !group.isEmpty {
            items.append(URLQueryItem(name: "group", value: group))
        }
        return items
    }
}
