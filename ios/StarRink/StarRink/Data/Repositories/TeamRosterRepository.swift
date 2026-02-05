import Foundation

protocol TeamRosterRepositoryProtocol {
    func getTournamentTeams(tournamentId: String) async throws -> TeamsResponse
    func getTeamRoster(teamId: String, tournamentId: String) async throws -> TeamRosterResponse
}

final class TeamRosterRepository: TeamRosterRepositoryProtocol {
    func getTournamentTeams(tournamentId: String) async throws -> TeamsResponse {
        try await APIClient.shared.request(endpoint: .tournamentTeams(id: tournamentId))
    }

    func getTeamRoster(teamId: String, tournamentId: String) async throws -> TeamRosterResponse {
        try await APIClient.shared.request(endpoint: .teamRoster(teamId: teamId, tournamentId: tournamentId))
    }
}
