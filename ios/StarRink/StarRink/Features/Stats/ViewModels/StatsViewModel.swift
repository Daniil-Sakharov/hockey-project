import Foundation

@MainActor
final class StatsViewModel: ObservableObject {
    @Published var isLoading = false
    @Published var gameStats: [GameStat] = []
    @Published var seasonTrends: [SeasonTrend] = []
    @Published var selectedTimeRange: TimeRange = .season

    enum TimeRange: String, CaseIterable {
        case month = "Месяц"
        case season = "Сезон"
        case allTime = "Все время"
    }

    // MARK: - Computed Properties

    var totalGoals: Int { gameStats.reduce(0) { $0 + $1.goals } }
    var totalAssists: Int { gameStats.reduce(0) { $0 + $1.assists } }
    var totalPoints: Int { gameStats.reduce(0) { $0 + $1.points } }
    var gamesPlayed: Int { gameStats.count }

    var averagePointsPerGame: Double {
        guard gamesPlayed > 0 else { return 0 }
        return Double(totalPoints) / Double(gamesPlayed)
    }

    var plusMinus: Int { gameStats.reduce(0) { $0 + $1.plusMinus } }

    // MARK: - Init

    init() {
        loadSampleData()
    }

    // MARK: - Data Loading

    func loadSampleData() {
        gameStats = GameStat.sampleData
        seasonTrends = SeasonTrend.sampleData
    }

    func loadStats() async {
        isLoading = true
        defer { isLoading = false }

        // TODO: Fetch from API
        // For now, use sample data
        try? await Task.sleep(nanoseconds: 500_000_000)
        loadSampleData()
    }

    func refresh() async {
        await loadStats()
    }
}
