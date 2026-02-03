import Foundation
import Combine

enum PositionFilter: String, CaseIterable {
    case all = "Все"
    case forward = "Нападающий"
    case defender = "Защитник"
    case goalie = "Вратарь"

    var apiValue: String {
        switch self {
        case .all: return "all"
        case .forward: return "forward"
        case .defender: return "defender"
        case .goalie: return "goalie"
        }
    }
}

@MainActor
final class PlayersSearchViewModel: ObservableObject {
    @Published var searchText = ""
    @Published var selectedPosition: PositionFilter = .all
    @Published var selectedBirthYear = 0
    @Published var selectedSeason = ""
    @Published var players: [PlayerItemDTO] = []
    @Published var total = 0
    @Published var seasons: [String] = []
    @Published var isLoading = false
    @Published var errorMessage: String?

    private let repository: PlayerRepositoryProtocol
    private var searchTask: Task<Void, Never>?

    init(repository: PlayerRepositoryProtocol = PlayerRepository()) {
        self.repository = repository
    }

    func loadSeasons() async {
        do {
            seasons = try await repository.getSeasons()
        } catch {
            seasons = []
        }
    }

    func search() {
        searchTask?.cancel()
        searchTask = Task {
            isLoading = true
            errorMessage = nil

            do {
                let response = try await repository.searchPlayers(
                    query: searchText,
                    position: selectedPosition.apiValue,
                    season: selectedSeason,
                    birthYear: selectedBirthYear,
                    limit: 30,
                    offset: 0
                )
                guard !Task.isCancelled else { return }
                players = response.players
                total = response.total
            } catch {
                guard !Task.isCancelled else { return }
                errorMessage = "Не удалось найти игроков"
            }

            isLoading = false
        }
    }

    func onSearchTextChanged() {
        searchTask?.cancel()
        searchTask = Task {
            try? await Task.sleep(nanoseconds: 400_000_000)
            guard !Task.isCancelled else { return }
            search()
        }
    }

    func resetFilters() {
        selectedPosition = .all
        selectedBirthYear = 0
        selectedSeason = ""
        search()
    }
}
