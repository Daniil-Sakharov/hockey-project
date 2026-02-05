import Foundation

@MainActor
final class TeamRosterViewModel: ObservableObject {
    @Published var team: TeamInfoDTO?
    @Published var players: [RosterPlayerDTO] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: TeamRosterRepositoryProtocol

    init(repository: TeamRosterRepositoryProtocol = TeamRosterRepository()) {
        self.repository = repository
    }

    func load(teamId: String, tournamentId: String) async {
        isLoading = true
        errorMessage = nil

        do {
            let response = try await repository.getTeamRoster(
                teamId: teamId,
                tournamentId: tournamentId
            )
            team = response.team
            players = response.players
        } catch {
            errorMessage = "Не удалось загрузить состав"
        }

        isLoading = false
    }
}
