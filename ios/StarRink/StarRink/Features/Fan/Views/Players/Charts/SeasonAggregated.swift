import Foundation

struct SeasonAggregated: Identifiable {
    let id: String
    let season: String
    let games: Int
    let goals: Int
    let assists: Int
    let points: Int
    let plusMinus: Int
    let penaltyMinutes: Int

    var avgGoals: Double { games > 0 ? Double(goals) / Double(games) : 0 }
    var avgAssists: Double { games > 0 ? Double(assists) / Double(games) : 0 }
    var avgPoints: Double { games > 0 ? Double(points) / Double(games) : 0 }
}
