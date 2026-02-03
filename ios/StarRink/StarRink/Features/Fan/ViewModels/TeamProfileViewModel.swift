import Foundation

@MainActor
final class TeamProfileViewModel: ObservableObject {
    @Published var team: TeamProfileDTO?
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: TeamProfileRepositoryProtocol

    init(repository: TeamProfileRepositoryProtocol = TeamProfileRepository()) {
        self.repository = repository
    }

    func load(teamId: String) async {
        isLoading = true
        errorMessage = nil

        do {
            team = try await repository.getTeamProfile(id: teamId)
        } catch {
            errorMessage = "Не удалось загрузить профиль команды"
        }

        isLoading = false
    }
}
